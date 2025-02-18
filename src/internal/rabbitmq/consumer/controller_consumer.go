package consumer

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/controllers"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/helper/response"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type ControllerConsumer struct {
	AuthController        *controllers.AuthController
	ChatController        *controllers.ChatController
	JobController         *controllers.JobController
	ProposalController    *controllers.ProposalController
	TransactionController *controllers.TransactionController
	UserController        *controllers.UserController
}

func NewControllerConsumer(
	authController *controllers.AuthController,
	chatController *controllers.ChatController,
	jobController *controllers.JobController,
	proposalController *controllers.ProposalController,
	transactionController *controllers.TransactionController,
	userController *controllers.UserController,
) *ControllerConsumer {
	return &ControllerConsumer{
		AuthController:        authController,
		ChatController:        chatController,
		JobController:         jobController,
		ProposalController:    proposalController,
		TransactionController: transactionController,
		UserController:        userController,
	}
}

func (controller ControllerConsumer) ConsumeAuthQueue(rabbitMQConfig *configs.RabbitMqConfig) {
	expectedQueueName := "AuthQueue"
	var queueName string
	for _, name := range rabbitMQConfig.Queue {
		if expectedQueueName == name.Name {
			queueName = name.Name
			break
		}
	}
	msgs, err := rabbitMQConfig.Channel.Consume(
		queueName,
		"AuthConsumer",
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

		// Handle error response
		if strings.HasPrefix(payload.Message, "responseError") {
			errorMessage := strings.TrimPrefix(payload.Message, "responseError")
			errorMessage = strings.TrimSpace(errorMessage)

			if errorMessage == "" {
				controller.AuthController.ResponseChannel <- response.Response[interface{}]{
					Code:    500,
					Message: "Error message is empty after 'responseError'",
					Data:    payload.Data,
				}
				continue
			}

			controller.AuthController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: fmt.Sprintf("Error occurred: %s", errorMessage),
				Data:    payload.Data,
			}
			continue
		}

		// Handle success response
		if payload.Message == "responseSuccess" {
			dataBytes, err := json.Marshal(payload.Data)
			if err != nil {
				fmt.Printf("Failed to marshal response data: %v\n", err)
				continue
			}

			var responseData entity.UserEntity
			err = json.Unmarshal(dataBytes, &responseData)
			if err != nil {
				fmt.Printf("Failed to unmarshal category products: %v\n", err)
				continue
			}
			controller.AuthController.ResponseChannel <- response.Response[interface{}]{
				Code:    200,
				Message: "Success",
				Data:    responseData,
			}
		} else {
			controller.AuthController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: "Unknown message type",
				Data:    nil,
			}
		}
	}
}
func (controller ControllerConsumer) ConsumeChatQueue(rabbitMQConfig *configs.RabbitMqConfig) {
	expectedQueueName := "ChatQueue"
	var queueName string
	for _, name := range rabbitMQConfig.Queue {
		if expectedQueueName == name.Name {
			queueName = name.Name
			break
		}
	}
	msgs, err := rabbitMQConfig.Channel.Consume(
		queueName,
		"ChatConsumer",
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

		// Handle error response
		if strings.HasPrefix(payload.Message, "responseError") {
			errorMessage := strings.TrimPrefix(payload.Message, "responseError")
			errorMessage = strings.TrimSpace(errorMessage)

			if errorMessage == "" {
				controller.ChatController.ResponseChannel <- response.Response[interface{}]{
					Code:    500,
					Message: "Error message is empty after 'responseError'",
					Data:    payload.Data,
				}
				continue
			}

			controller.ChatController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: fmt.Sprintf("Error occurred: %s", errorMessage),
				Data:    payload.Data,
			}
			continue
		}

		// Handle success response
		if payload.Message == "responseSuccess" {
			dataBytes, err := json.Marshal(payload.Data)
			if err != nil {
				fmt.Printf("Failed to marshal response data: %v\n", err)
				continue
			}

			var responseData *entity.ChatEntity
			err = json.Unmarshal(dataBytes, &responseData)
			if err != nil {
				fmt.Printf("Failed to unmarshal category products: %v\n", err)
				continue
			}
			controller.ChatController.ResponseChannel <- response.Response[interface{}]{
				Code:    200,
				Message: "Success",
				Data:    responseData,
			}
		} else {
			controller.ChatController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: "Unknown message type",
				Data:    nil,
			}
		}
	}
}
func (controller ControllerConsumer) ConsumeJobQueue(rabbitMQConfig *configs.RabbitMqConfig) {
	expectedQueueName := "JobQueue"
	var queueName string
	for _, name := range rabbitMQConfig.Queue {
		if expectedQueueName == name.Name {
			queueName = name.Name
			break
		}
	}
	msgs, err := rabbitMQConfig.Channel.Consume(
		queueName,
		"JobConsumer",
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

		// Handle error response
		if strings.HasPrefix(payload.Message, "responseError") {
			errorMessage := strings.TrimPrefix(payload.Message, "responseError")
			errorMessage = strings.TrimSpace(errorMessage)

			if errorMessage == "" {
				controller.JobController.ResponseChannel <- response.Response[interface{}]{
					Code:    500,
					Message: "Error message is empty after 'responseError'",
					Data:    payload.Data,
				}
				continue
			}

			controller.JobController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: fmt.Sprintf("Error occurred: %s", errorMessage),
				Data:    payload.Data,
			}
			continue
		}

		// Handle success response
		if payload.Message == "responseSuccess" {
			dataBytes, err := json.Marshal(payload.Data)
			if err != nil {
				fmt.Printf("Failed to marshal response data: %v\n", err)
				continue
			}

			var responseData *entity.JobEntity
			err = json.Unmarshal(dataBytes, &responseData)
			if err != nil {
				fmt.Printf("Failed to unmarshal category products: %v\n", err)
				continue
			}
			controller.JobController.ResponseChannel <- response.Response[interface{}]{
				Code:    200,
				Message: "Success",
				Data:    responseData,
			}
		} else {
			controller.JobController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: "Unknown message type",
				Data:    nil,
			}
		}
	}
}
func (controller ControllerConsumer) ConsumeProposalQueue(rabbitMQConfig *configs.RabbitMqConfig) {
	expectedQueueName := "ProposalQueue"
	var queueName string
	for _, name := range rabbitMQConfig.Queue {
		if expectedQueueName == name.Name {
			queueName = name.Name
			break
		}
	}
	msgs, err := rabbitMQConfig.Channel.Consume(
		queueName,
		"ProposalConsumer",
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

		// Handle error response
		if strings.HasPrefix(payload.Message, "responseError") {
			errorMessage := strings.TrimPrefix(payload.Message, "responseError")
			errorMessage = strings.TrimSpace(errorMessage)

			if errorMessage == "" {
				controller.ProposalController.ResponseChannel <- response.Response[interface{}]{
					Code:    500,
					Message: "Error message is empty after 'responseError'",
					Data:    payload.Data,
				}
				continue
			}

			controller.ProposalController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: fmt.Sprintf("Error occurred: %s", errorMessage),
				Data:    payload.Data,
			}
			continue
		}

		// Handle success response
		if payload.Message == "responseSuccess" {
			dataBytes, err := json.Marshal(payload.Data)
			if err != nil {
				fmt.Printf("Failed to marshal response data: %v\n", err)
				continue
			}

			var responseData *entity.ProposalEntity
			err = json.Unmarshal(dataBytes, &responseData)
			if err != nil {
				fmt.Printf("Failed to unmarshal category products: %v\n", err)
				continue
			}
			controller.ProposalController.ResponseChannel <- response.Response[interface{}]{
				Code:    200,
				Message: "Success",
				Data:    responseData,
			}
		} else {
			controller.ProposalController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: "Unknown message type",
				Data:    nil,
			}
		}
	}
}
func (controller ControllerConsumer) ConsumeTransactionQueue(rabbitMQConfig *configs.RabbitMqConfig) {
	expectedQueueName := "TransactionQueue"
	var queueName string
	for _, name := range rabbitMQConfig.Queue {
		if expectedQueueName == name.Name {
			queueName = name.Name
			break
		}
	}
	msgs, err := rabbitMQConfig.Channel.Consume(
		queueName,
		"TransactionConsumer",
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

		// Handle error response
		if strings.HasPrefix(payload.Message, "responseError") {
			errorMessage := strings.TrimPrefix(payload.Message, "responseError")
			errorMessage = strings.TrimSpace(errorMessage)

			if errorMessage == "" {
				controller.TransactionController.ResponseChannel <- response.Response[interface{}]{
					Code:    500,
					Message: "Error message is empty after 'responseError'",
					Data:    payload.Data,
				}
				continue
			}

			controller.TransactionController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: fmt.Sprintf("Error occurred: %s", errorMessage),
				Data:    payload.Data,
			}
			continue
		}

		// Handle success response
		if payload.Message == "responseSuccess" {
			dataBytes, err := json.Marshal(payload.Data)
			if err != nil {
				fmt.Printf("Failed to marshal response data: %v\n", err)
				continue
			}

			var responseData *entity.TransactionEntity
			err = json.Unmarshal(dataBytes, &responseData)
			if err != nil {
				fmt.Printf("Failed to unmarshal category products: %v\n", err)
				continue
			}
			controller.TransactionController.ResponseChannel <- response.Response[interface{}]{
				Code:    200,
				Message: "Success",
				Data:    responseData,
			}
		} else {
			controller.TransactionController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: "Unknown message type",
				Data:    nil,
			}
		}
	}
}
func (controller ControllerConsumer) ConsumeUserQueue(rabbitMQConfig *configs.RabbitMqConfig) {
	expectedQueueName := "UserQueue"
	var queueName string
	for _, name := range rabbitMQConfig.Queue {
		if expectedQueueName == name.Name {
			queueName = name.Name
			break
		}
	}
	msgs, err := rabbitMQConfig.Channel.Consume(
		queueName,
		"UserConsumer",
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

		// Handle error response
		if strings.HasPrefix(payload.Message, "responseError") {
			errorMessage := strings.TrimPrefix(payload.Message, "responseError")
			errorMessage = strings.TrimSpace(errorMessage)

			if errorMessage == "" {
				controller.UserController.ResponseChannel <- response.Response[interface{}]{
					Code:    500,
					Message: "Error message is empty after 'responseError'",
					Data:    payload.Data,
				}
				continue
			}

			controller.UserController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: fmt.Sprintf("Error occurred: %s", errorMessage),
				Data:    payload.Data,
			}
			continue
		}

		// Handle success response
		if payload.Message == "responseSuccess" {
			dataBytes, err := json.Marshal(payload.Data)
			if err != nil {
				fmt.Printf("Failed to marshal response data: %v\n", err)
				continue
			}

			var responseData *entity.UserEntity
			err = json.Unmarshal(dataBytes, &responseData)
			if err != nil {
				fmt.Printf("Failed to unmarshal category products: %v\n", err)
				continue
			}
			controller.UserController.ResponseChannel <- response.Response[interface{}]{
				Code:    200,
				Message: "Success",
				Data:    responseData,
			}
		} else {
			controller.UserController.ResponseChannel <- response.Response[interface{}]{
				Code:    400,
				Message: "Unknown message type",
				Data:    nil,
			}
		}
	}
}
