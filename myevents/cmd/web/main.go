package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mihirkelkar/microservices/lib/configuration"
	"github.com/mihirkelkar/microservices/lib/msgqueue"
	msgqueue_amqp "github.com/mihirkelkar/microservices/lib/msgqueue/amqp"
	msgqueue_kfka "github.com/mihirkelkar/microservices/lib/msgqueue/kafka"
	"github.com/mihirkelkar/microservices/myevents/pkg/models"
	"github.com/streadway/amqp"
)

func main() {

	var eventEmitter msgqueue.EventEmitter

	infLog := log.New(os.Stdout, "INFO\n", log.Ldate|log.Ltime)
	//define the error log.
	errLog := log.New(os.Stderr, "ERROR\n", log.Ldate|log.Ltime)

	//get the configurations
	conf, err := configuration.ReadConfig("event_config.json")
	if err != nil {
		infLog.Print("Using Default Configuration")
	}

	//open a connection to the mongodb database
	session, err := models.NewMongoConnection(conf.DatabaseURL)
	if err != nil {
		errLog.Panic(err)
	}
	defer session.Close()

	switch conf.EventBroker {
	case "rabbitmq":
		conn, err := amqp.Dial(conf.EventBrokerURL)
		if err != nil {
			errLog.Panic(err)
		}
		eventEmitter, err = msgqueue_amqp.NewEventEmitter(conn, "events")
		if err != nil {
			errLog.Panic(err)
		}

	case "kafka":
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Net.DialTimeout = 3 * time.Second
		config.Net.ReadTimeout = 3 * time.Second
		config.Net.WriteTimeout = 3 * time.Second
		fmt.Println(conf.EventBrokerURL)
		client, err := sarama.NewClient([]string{conf.EventBrokerURL}, config)
		if err != nil {
			errLog.Panic(err)
		}

		eventEmitter, err = msgqueue_kfka.NewKafkaEventEmitter(client)
		if err != nil {
			errLog.Panic(err)
		}

	default:
		errLog.Panic("No event broker specified")
	}

	//create a new event service
	eventService := models.NewEventService(session, "myevents", "events")

	//instantiate a routes app
	app := Application{ErrorLog: errLog,
		InfoLog:      infLog,
		EventService: eventService,
		EventEmitter: eventEmitter}

	//instatntiate a server that accepts the routes as handlers.
	server := http.Server{
		Addr:     conf.Port,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infLog.Printf("Starting the server on port :%s", conf.Port)
	err = server.ListenAndServe()
	if err != nil {
		server.ErrorLog.Fatal(err.Error())
	}

}
