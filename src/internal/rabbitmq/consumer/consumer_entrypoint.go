package consumer

import (
	"bantu-backend/src/configs"
)

type ConsumerEntrypoint struct {
	ControllerConsumer *ControllerConsumer
	ServicesConsumer   *ServiceConsumer
	RabbitMq           *configs.RabbitMqConfig
}
type RabbitMQPayload struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Code    int    `json:"code"`
}

func NewConsumerEntrypointInit(
	controllerConsumer *ControllerConsumer,
	// serviceConsumer *ServiceConsumer,
	rabbitMQConfig *configs.RabbitMqConfig,
) *ConsumerEntrypoint {
	return &ConsumerEntrypoint{
		ControllerConsumer: controllerConsumer,
		// ServicesConsumer:   serviceConsumer,
		RabbitMq: rabbitMQConfig,
	}
}
func ControllerConsumerStart(consumerEntrypoint *ConsumerEntrypoint) {
	go consumerEntrypoint.ControllerConsumer.ConsumeErrorQueue(consumerEntrypoint.RabbitMq)
	go consumerEntrypoint.ControllerConsumer.ConsumeAuthQueue(consumerEntrypoint.RabbitMq)
	go consumerEntrypoint.ControllerConsumer.ConsumeChatQueue(consumerEntrypoint.RabbitMq)
	go consumerEntrypoint.ControllerConsumer.ConsumeJobQueue(consumerEntrypoint.RabbitMq)
	go consumerEntrypoint.ControllerConsumer.ConsumeProposalQueue(consumerEntrypoint.RabbitMq)
	go consumerEntrypoint.ControllerConsumer.ConsumeTransactionQueue(consumerEntrypoint.RabbitMq)
	go consumerEntrypoint.ControllerConsumer.ConsumeUserQueue(consumerEntrypoint.RabbitMq)
}
func ServiceConsumerStart(consumerEntrypoint *ConsumerEntrypoint) {
}
func (consumerEntrypoint *ConsumerEntrypoint) ConsumerEntrypointStart() {
	ControllerConsumerStart(consumerEntrypoint)
	ServiceConsumerStart(consumerEntrypoint)
}
