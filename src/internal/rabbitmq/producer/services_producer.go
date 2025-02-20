package producer

import (
	"bantu-backend/src/configs"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type ServicesProducer struct {
	Env *configs.RabbitMqEnv
}

func CreateNewServicesProducer(env *configs.RabbitMqEnv) *ServicesProducer {
	return &ServicesProducer{
		Env: env,
	}
}
func (servicesProducer *ServicesProducer) CreateMessageError(channelRabbitMQ *amqp.Channel, data interface{}) error {
	queueName := servicesProducer.Env.Queues[6]
	payload := map[string]interface{}{
		"data": data,
	}
	messageBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message body: %w", err)
	}
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(messageBody),
	}
	if err := channelRabbitMQ.Publish(
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,     // immediate
		msg,       // message to publish
	); err != nil {
		return fmt.Errorf("failed to publish message to queue: %w", err)
	}
	return nil
}

func (servicesProducer *ServicesProducer) CreateMessageAuth(channelRabbitMQ *amqp.Channel, data any) error {
	queueName := servicesProducer.Env.Queues[0]
	payload := map[string]any{
		"data":    data,
		"channel": channelRabbitMQ,
	}
	messageBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message body: %w", err)
	}
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(messageBody),
	}
	if err := channelRabbitMQ.Publish(
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,     // immediate
		msg,       // message to publish
	); err != nil {
		return fmt.Errorf("failed to publish message to queue: %w", err)
	}
	return nil
}

func (servicesProducer *ServicesProducer) CreateMessageJob(channelRabbitMQ *amqp.Channel, messageType string, data interface{}) error {
	queueName := servicesProducer.Env.Queues[2]
	payload := map[string]interface{}{
		"message": messageType,
		"data":    data,
		"channel": channelRabbitMQ,
	}
	messageBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message body: %w", err)
	}
	message := amqp.Publishing{
		ContentType: "application/json",
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
