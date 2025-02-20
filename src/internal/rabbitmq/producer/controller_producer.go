package producer

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type ControllerProducer struct{}

func CreateNewControllerProducer() *ControllerProducer {
	controllerProducer := &ControllerProducer{}
	return controllerProducer
}

func (*ControllerProducer) ProducerJobsQueue(channelRabbitMQ *amqp.Channel, data map[string]interface{}) error {
	queueName := "JobQueue"
	payload := map[string]interface{}{
		"message": "Start Scraping",
		"job":     data,
		"channel": channelRabbitMQ,
	}
	messageBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message body: %w", err)
	}
	message := amqp.Publishing{
		ContentType: "application/text",
		Body:        []byte(messageBody),
	}
	if err := channelRabbitMQ.Publish(
		"",
		queueName,
		false,
		false,
		message,
	); err != nil {
		return fmt.Errorf("failed to publish message to queue: %w", err)
	}
	return nil
}
