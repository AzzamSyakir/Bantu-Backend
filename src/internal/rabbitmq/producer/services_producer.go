package producer

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type ServicesProducer struct{}

func CreateNewServicesProducer() *ServicesProducer {
	return &ServicesProducer{}
}

func (*ServicesProducer) CreateMessageAuth(channelRabbitMQ *amqp.Channel, seller string) error {
	queueName := "AuthQueue"
	payload := map[string]interface{}{
		"message": "Auth Process",
		"seller":  seller,
		"channel": channelRabbitMQ,
	}
	messageBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message body: %w", err)
	}
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(messageBody),
	}
	if err := channelRabbitMQ.Publish(
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,     // immediate
		message,   // message to publish
	); err != nil {
		return fmt.Errorf("failed to publish message to queue: %w", err)
	}
	return nil
}
