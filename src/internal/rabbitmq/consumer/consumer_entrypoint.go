package consumer

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/controllers"
)

type ConsumerEntrypoint struct {
	ControllerConsumer *ControllerConsumer
	ServicesConsumer   *ControllerConsumer
}

func NewConsumerEntrypointInit(
	rabbitMQConfig *configs.RabbitMqConfig,
	authController        *controllers.AuthController,
	chatController        *controllers.ChatController,
	jobController         *controllers.JobController,
	proposalController    *controllers.ProposalController,
	transactionController *controllers.TransactionController,
	userController        *controllers.UserController,
) *ConsumerEntrypoint {
	return &ConsumerEntrypoint{
		ControllerConsumer: &ControllerConsumer{
			AuthController: authController,
			ChatController: chatController,
			JobController: jobController,
			ProposalController: proposalController,
			TransactionController: transactionController,
			UserController: userController,
		},
		ServicesConsumer: ,
	}
}

func (consumerEntrypoint *ConsumerEntrypoint) ConsumerEntrypointStart() {
	go consumerEntrypoint.ScrapingConsumer.ConsumeMessageAllSellerProduct(consumerEntrypoint.RabbitMQ)
	go consumerEntrypoint.ScrapingConsumer.ConsumeMessageSoldSellerProduct(consumerEntrypoint.RabbitMQ)
	go consumerEntrypoint.MainConsumer.ConsumeSellerProductResponse(consumerEntrypoint.RabbitMQ)
}
