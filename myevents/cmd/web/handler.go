package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/microservices/contracts"
	"github.com/mihirkelkar/microservices/myevents/pkg/models"
)

/*
This is the module where we define all the handler functions for this
micro-service.
*/

func (app *Application) GetEvents(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Print("The GetEvents Function Handler was Invoked")
	events, err := app.EventService.FindAllEvents()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&JsonError{Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(events)
	return
}

func (app *Application) GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	app.InfoLog.Printf("The GetEventsByID Function Handler was Invoked for ID: %s", id)
	event, err := app.EventService.FindEventByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&JsonError{Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(event)
	return
}

func (app *Application) GetEventByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	app.InfoLog.Printf("The GetEventsByName Function Handler was invoked for Name:%s", name)
	event, err := app.EventService.FindEventByName(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&JsonError{Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(event)
	return
}

func (app *Application) CreateEvent(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Print("The function to create an event was called")
	event := models.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&JsonError{Message: err.Error()})
		return
	}
	id, err := app.EventService.AddEvent(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&JsonError{Message: err.Error()})
		return
	}

	msg := contracts.EventCreatedEvent{
		ID:        hex.EncodeToString([]byte(id)),
		Name:      event.Name,
		StartDate: time.Unix(event.StartDate, 0),
		EndDate:   time.Unix(event.EndDate, 0),
		Location:  string(event.Location.ID),
	}

	err = app.EventEmitter.Emit(&msg)
	if err != nil {
		fmt.Println(err)
	}
	return
}
