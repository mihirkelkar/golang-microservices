package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/microservices/lib/msgqueue"
	"github.com/mihirkelkar/microservices/myevents/pkg/models"
)

type Application struct {
	ErrorLog     *log.Logger
	InfoLog      *log.Logger
	EventService models.EventService //Event Service is the interface that lets you talk to MongoDB
	EventEmitter msgqueue.EventEmitter
}

func (app *Application) routes() *mux.Router {
	router := mux.NewRouter()
	//GET requests
	router.HandleFunc("/events", app.GetEvents).Methods("GET")
	router.HandleFunc("/events/id/{id:[a-z,0-9,_,-]+}", app.GetEventByID).Methods("GET")
	router.HandleFunc("/events/name/{name}", app.GetEventByName).Methods("GET")

	//POST requests
	router.HandleFunc("/events", app.CreateEvent).Methods("POST")
	return router
}
