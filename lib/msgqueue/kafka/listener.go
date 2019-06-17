package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/mihirkelkar/microservices/lib/msgqueue"
)

type kafkaEventListener struct {
	Consumer   sarama.Consumer
	Partitions []int32
}

func NewKafkaEventListener(client sarama.Client, partitions []int32) (*kafkaEventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &kafkaEventListener{
		Consumer:   consumer,
		Partitions: partitions,
	}, nil
}

func (kl *kafkaEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error

	topic := "events"
	results := make(chan msgqueue.Event)
	errors := make(chan error)
	partitions := kl.Partitions

	if len(partitions) == 0 {
		partitions, err = kl.Consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, partition := range partitions {
		log.Printf("Consuming partition %s:%d", topic, partition)

		pConsumer, err := kl.Consumer.ConsumePartition(topic, partition, 0)
		if err != nil {
			return nil, nil, err
		}

		//Consume messages from the GoConsumer
		go func() {
			//look at the message from the partition consumer
			for msg := range pConsumer.Messages() {
				body := messageEnvelope{}
				err := json.Unmarshal(msg.Value, body)
				if err != nil {
					errors <- fmt.Errorf("Error: Could not unmarshal JSON message: %v", err)
					continue
				}
				var event msgqueue.Event
				err = json.Unmarshal(body.Payload.([]byte), &event)
				if err != nil {
					errors <- fmt.Errorf("Error: Could not unmarshal JSON payload from message into event")
					continue
				}
				results <- event
			}
		}()

		//Consume errors from the GoConsumer
		go func() {
			for err := range pConsumer.Errors() {
				errors <- err
			}
		}()
	}
	return results, errors, nil
}
