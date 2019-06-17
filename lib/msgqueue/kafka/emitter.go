package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/mihirkelkar/microservices/lib/msgqueue"
)

type kafkaEventEmitter struct {
	Producer sarama.SyncProducer
}

type messageEnvelope struct {
	EventName string      `json:"eventName"`
	Payload   interface{} `json:"payload"`
}

//Creates a new eventEmitter object to be used with kafka
func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	eventEmitter := &kafkaEventEmitter{Producer: producer}
	return eventEmitter, nil
}

//Emit: Emits events to the Kafka stream
func (ke *kafkaEventEmitter) Emit(event msgqueue.Event) error {
	msg_evlp := messageEnvelope{EventName: event.EventName(), Payload: event}
	jsonBody, err := json.Marshal(msg_evlp)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{Topic: event.EventName(), Value: sarama.ByteEncoder(jsonBody)}

	_, _, err = ke.Producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
