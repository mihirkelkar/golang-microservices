package main

import "github.com/streadway/amqp"

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic("Error: Could not establish connection to amqp")
	}

	channel, err := connection.Channel()
	if err != nil {
		panic("Error: Could not create a channel")
	}

	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	message := amqp.Publishing{Body: []byte("Hello World !")}

	err = channel.Publish("events", "routing-key", false, false, message)
	if err != nil {
		panic("Error publishing message")
	}
}
