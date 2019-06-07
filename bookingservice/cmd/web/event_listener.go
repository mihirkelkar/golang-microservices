package main

import (
	"fmt"
	"log"

	"github.com/mihirkelkar/microservices/lib/msgqueue"
)

func (app *Application) ProcessEvents() {
	log.Println("Listening to Events")

	recieved, errors, err := app.EventListener.Listen("event.created")
	if err != nil {
		app.ErrorLog.Panic(err)
	}
	for {
		select {
		case evt := <-recieved:
			fmt.Printf("got event %T", evt)
			app.handleEvent(evt)
		case err := <-errors:
			fmt.Printf("got error while receiving event: %s\n", err)
		}
	}
}

func (app *Application) handleEvent(event msgqueue.Event) {
	fmt.Print(event)
}
