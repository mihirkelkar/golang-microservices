package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/mihirkelkar/microservices/contracts"
	"github.com/mihirkelkar/microservices/lib/msgqueue"
	"github.com/streadway/amqp"
)

const eventNameHeader = "x-event-name"

type ampqEventListener struct {
	connection *amqp.Connection
	exchange   string
	queue      string
}

func (al *ampqEventListener) setup() error {
	channel, err := al.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	//declare an exchange
	err = channel.ExchangeDeclare(al.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}
	//declare a queue
	_, err = channel.QueueDeclare(al.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not declare queue %s: %s", al.queue, err)
	}

	return nil
}

func (al *ampqEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := al.connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	//create a binding between the queue and the exhange for every eventName that we want to listen to.
	for _, event := range eventNames {
		if err = channel.QueueBind(al.queue, event, al.exchange, false, nil); err != nil {
			return nil, nil, fmt.Errorf("Could not bind event %s to queue %s: %s", event, al.queue, err)
		}
	}

	//if the binding is succesful go ahead and consume the events.
	msgs, err := channel.Consume(al.queue, "", false, false, false, false, nil)

	if err != nil {
		return nil, nil, fmt.Errorf("Could not consume events in queue %s", al.queue)
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers[eventNameHeader]
			if !ok {
				errors <- fmt.Errorf("message did not contain %s header", eventNameHeader)
				msg.Nack(false, false)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("header did not contain string")
				msg.Nack(false, false)
				continue
			}
			var event msgqueue.Event
			if eventName == "event.created" {
				event = &contracts.EventCreatedEvent{}
				err := json.Unmarshal(msg.Body, event)
				fmt.Println(event)
				if err != nil {
					errors <- fmt.Errorf("could not unmarshal event into json")
					msg.Nack(false, false)
					continue
				}
			} else {
				errors <- fmt.Errorf("Not an event created event. We can only process those right now")
				msg.Nack(false, false)
				continue
			}
			fmt.Println(event)
			events <- event
			msg.Ack(false)
		}
	}()
	return events, errors, nil
}

//Creates New Event Listener
func NewEventListener(conn *amqp.Connection, exchange string, queue string) (msgqueue.EventListener, error) {
	listener := ampqEventListener{connection: conn, exchange: exchange, queue: queue}

	if err := listener.setup(); err != nil {
		return nil, err
	}
	return &listener, nil
}
