package consumer

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/controllers"
	"bantu-backend/src/internal/models/response"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

type ControllerConsumer struct {
	ResponseChannel       *response.ResponseChannel
	Env                   *configs.RabbitMqEnv
	AuthController        *controllers.AuthController
	ChatController        *controllers.ChatController
	JobController         *controllers.JobController
	ProposalController    *controllers.ProposalController
	TransactionController *controllers.TransactionController
	UserController        *controllers.UserController
}

func NewControllerConsumer(
	env *configs.RabbitMqEnv,
	authController *controllers.AuthController,
	chatController *controllers.ChatController,
	jobController *controllers.JobController,
	proposalController *controllers.ProposalController,
	transactionController *controllers.TransactionController,
	userController *controllers.UserController,
	responseChannel *response.ResponseChannel,
) *ControllerConsumer {
	return &ControllerConsumer{
		Env:                   env,
		AuthController:        authController,
		ChatController:        chatController,
		JobController:         jobController,
		ProposalController:    proposalController,
		TransactionController: transactionController,
		UserController:        userController,
		ResponseChannel:       responseChannel,
	}
}

func (consumerController ControllerConsumer) ConsumeSuccess(rabbitMQConfig *configs.RabbitMqConfig) {
	for _, q := range rabbitMQConfig.Queue {
		consumerTag := strings.Replace(q.Name, "Queue", "Consumer", 1)
		if consumerTag == "ErrorConsumer" {
			continue
		}
		msgs, err := rabbitMQConfig.Channel.Consume(
			q.Name,
			consumerTag,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Printf("Queue '%s' is not available. Retrying in 5 seconds... Error: %v\n", q.Name, err)
			continue
		}
		go func(queue string, msgs <-chan amqp.Delivery) {
			for msg := range msgs {
				var payload RabbitMQPayload
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					log.Printf("Failed to unmarshal message from queue %s: %v", queue, err)
					continue
				}
				consumerController.ResponseChannel.ResponseSuccess <- response.Response[interface{}]{
					Code:    200,
					Message: "Success",
					Data:    payload.Data,
				}
			}
		}(q.Name, msgs)
	}
}
func (consumerController ControllerConsumer) ConsumeError(rabbitMQConfig *configs.RabbitMqConfig) {
	expectedQueueName := consumerController.Env.Queues[6]
	var queueName string
	for _, name := range rabbitMQConfig.Queue {
		if expectedQueueName == name.Name {
			queueName = name.Name
			break
		}
	}
	msgs, err := rabbitMQConfig.Channel.Consume(
		queueName,
		"ErrorConsumer",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("Queue '%s' not available. Retrying in 5 seconds... Error: %v\n", queueName, err)
		return
	}
	for msg := range msgs {
		var payload RabbitMQPayload
		// Parse JSON message
		err := json.Unmarshal(msg.Body, &payload)
		if err != nil {
			log.Fatal("Failed to unmarshal message: ", err)
		}
		// Handle response
		consumerController.ResponseChannel.ResponseError <- response.Response[interface{}]{
			Code:    payload.StatusCode,
			Message: "Error",
			Data:    payload.Data,
		}
	}

}
