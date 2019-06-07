package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func main() {
	//connect to amqp
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic("Error: Could not establish connection to amqp")
	}
	fmt.Println("Connected to amqp")

	channel, err := connection.Channel()
	if err != nil {
		panic("Error: Could not create a channel")
	}

	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_, err = channel.QueueDeclare("my_queue", true, false, false, false, nil)
	if err != nil {
		panic("error while connecting to the queue.")
	}

	//bind to queue.
	err = channel.QueueBind("my_queue", "#", "events", false, nil)
	if err != nil {
		panic("Could not bind queue")
	}

	msgs, err := channel.Consume("my_queue", "", false, false, false, false, nil)
	if err != nil {
		panic("error consuming this channel")
	}

	fmt.Println("Consuming")

	for msg := range msgs {
		fmt.Println("Message recieved: " + string(msg.Body))
		msg.Ack(false)
	}
}
