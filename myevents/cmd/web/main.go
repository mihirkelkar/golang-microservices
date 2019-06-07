package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mihirkelkar/microservices/lib/configuration"
	msgqueue_amqp "github.com/mihirkelkar/microservices/lib/msgqueue/amqp"
	"github.com/mihirkelkar/microservices/myevents/pkg/models"
	"github.com/streadway/amqp"
)

func main() {

	infLog := log.New(os.Stdout, "INFO\n", log.Ldate|log.Ltime)
	//define the error log.
	errLog := log.New(os.Stderr, "ERROR\n", log.Ldate|log.Ltime)

	//get the configurations
	conf, err := configuration.ReadConfig("")
	if err != nil {
		infLog.Print("Using Default Configuration")
	}

	//open a connection to the mongodb database
	session, err := models.NewMongoConnection(conf.DatabaseURL)
	if err != nil {
		errLog.Panic(err)
	}
	defer session.Close()

	conn, err := amqp.Dial(conf.EventBrokerURL)
	if err != nil {
		errLog.Panic(err)
	}
	defer conn.Close()

	//create a new event Emitter.
	eventEmitter, err := msgqueue_amqp.NewEventEmitter(conn, conf.Exchange)

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
