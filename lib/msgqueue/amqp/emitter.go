package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/mihirkelkar/microservices/lib/msgqueue"
	"github.com/streadway/amqp"
)

//This fits the EventEmitter interface described in msgqueue/event.go
type ampqEventEmitter struct {
	connection *amqp.Connection
	exchange   string
}

func (ae *ampqEventEmitter) setup() error {
	channel, err := ae.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	//declare an exhange.
	err = channel.ExchangeDeclare(ae.exchange, "topic", true, false, false, false, nil)
	return err
}

func (ae *ampqEventEmitter) Emit(event msgqueue.Event) error {
	channel, err := ae.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	jsonBody, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("could not JSON-serialize event: %s", err)
	}

	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event.EventName()},
		ContentType: "application/json",
		Body:        jsonBody,
	}
	//the params are exhange, queue, boolean, boolean and payload
	err = channel.Publish(ae.exchange, event.EventName(), false, false, msg)
	return err
}

//Creates a new struct that can fit the EventEmitter interface.
func NewEventEmitter(conn *amqp.Connection, exchange string) (msgqueue.EventEmitter, error) {
	emitter := ampqEventEmitter{connection: conn, exchange: exchange}

	if err := emitter.setup(); err != nil {
		return nil, err
	}
	return &emitter, nil
}
