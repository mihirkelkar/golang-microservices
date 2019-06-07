package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/microservices/lib/msgqueue"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	EventListener msgqueue.EventListener
}

func (app *Application) routes() *mux.Router {
	router := mux.NewRouter()
	return router
}
