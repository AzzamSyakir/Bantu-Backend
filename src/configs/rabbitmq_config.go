package configs

import (
	"fmt"
	"strings"

	"github.com/streadway/amqp"
)

type RabbitMqConfig struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      []*amqp.Queue
	Env        *EnvConfig
}

func NewRabbitMqConfig(env *EnvConfig) *RabbitMqConfig {
	rabbitMq := &RabbitMqConfig{}
	connection, err := rabbitMq.RabbitMQConnection()
	if err != nil {
		panic("failed to establish RabbitMQ connection: " + err.Error())
	}
	channel, err := RabbitMqChannel(connection)
	if err != nil {
		panic("failed to establish RabbitMQ channel: " + err.Error())
	}
	queue, err := rabbitMq.RabbitMqQueue(channel)
	if err != nil {
		panic("failed to establish RabbitMQ queue: " + err.Error())
	}
	rabbitMqConfig := &RabbitMqConfig{
		Connection: connection,
		Channel:    channel,
		Queue:      queue,
		Env:        env,
	}
	return rabbitMqConfig
}

func (rabbitMqConfig *RabbitMqConfig) RabbitMQConnection() (*amqp.Connection, error) {
	rabbitMqHost := rabbitMqConfig.Env.RabbitMq.Host
	rabbitMqUser := rabbitMqConfig.Env.RabbitMq.User
	rabbitMqPass := rabbitMqConfig.Env.RabbitMq.Password
	rabbitMqPort := rabbitMqConfig.Env.RabbitMq.Port
	amqpServerURL := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitMqUser, rabbitMqPass, rabbitMqHost, rabbitMqPort)
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return nil, err
	}
	return connectRabbitMQ, nil
}

func RabbitMqChannel(connection *amqp.Connection) (*amqp.Channel, error) {

	channelRabbitMQ, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	return channelRabbitMQ, nil
}

func (rabbitMqConfig *RabbitMqConfig) RabbitMqQueue(channel *amqp.Channel) ([]*amqp.Queue, error) {
	var declaredQueues []*amqp.Queue
	queueNamesStr := rabbitMqConfig.Env.RabbitMq.Queue

	// Trim spaces and split properly
	rawQueueNames := strings.Split(queueNamesStr, ",")
	queueNames := make([]string, 0, len(rawQueueNames))

	for _, name := range rawQueueNames {
		trimmedName := strings.TrimSpace(name)
		if trimmedName != "" {
			queueNames = append(queueNames, trimmedName)
		}
	}
	for _, name := range queueNames {
		rabbitmqQueue, err := channel.QueueDeclare(
			name,
			true,  // Durable
			false, // Auto-delete
			true,  // Exclusive
			false, // No-wait
			nil,   // Argument
		)
		if err != nil {
			return nil, err
		}

		declaredQueues = append(declaredQueues, &rabbitmqQueue)
	}

	return declaredQueues, nil
}
