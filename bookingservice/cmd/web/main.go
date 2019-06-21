package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/mihirkelkar/microservices/lib/configuration"
	"github.com/mihirkelkar/microservices/lib/msgqueue"
	msgqueue_amqp "github.com/mihirkelkar/microservices/lib/msgqueue/amqp"
	msgqueue_kfka "github.com/mihirkelkar/microservices/lib/msgqueue/kafka"
	"github.com/streadway/amqp"
)

func main() {
	infLog := log.New(os.Stdout, "INFO\n", log.Ldate|log.Ltime)
	//define the error log.
	errLog := log.New(os.Stderr, "ERROR\n", log.Ldate|log.Ltime)

	//get the configurations
	conf, err := configuration.ReadConfig("booking_config.json")
	if err != nil {
		infLog.Print("Using Default Configuration")
	}

	var eventListener msgqueue.EventListener

	switch conf.EventBroker {
	case "rabbitmq":
		conn, err := amqp.Dial(conf.EventBrokerURL)
		if err != nil {
			errLog.Panic(err)
		}
		defer conn.Close()

		//create a new event Listener
		eventListener, err = msgqueue_amqp.NewEventListener(conn, conf.Exchange, conf.Queue)
	case "kafka":
		config := sarama.NewConfig()
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
		fmt.Println(conf.EventBrokerURL)
		client, err := sarama.NewClient([]string{conf.EventBrokerURL}, config)
		if err != nil {
			errLog.Panic(err)
		}

		eventListener, err = msgqueue_kfka.NewKafkaEventListener(client, []int32{})
		if err != nil {
			errLog.Panic(err)
		}

	}

	//create a new event service

	//instantiate a routes app
	app := Application{
		ErrorLog:      errLog,
		InfoLog:       infLog,
		EventListener: eventListener,
	}

	//instatntiate a server that accepts the routes as handlers.
	server := http.Server{
		Addr:     conf.Port,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infLog.Printf("Starting to process events")
	go app.ProcessEvents()
	infLog.Printf("Starting the server on port :%s", conf.Port)
	err = server.ListenAndServe()
	if err != nil {
		server.ErrorLog.Fatal(err.Error())
	}
}
