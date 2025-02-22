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
	Message    string `json:"message"`
	Data       any    `json:"data"`
	StatusCode int    `json:"status_code"`
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
	go consumerEntrypoint.ControllerConsumer.ConsumeError(consumerEntrypoint.RabbitMq)
	go consumerEntrypoint.ControllerConsumer.ConsumeSuccess(consumerEntrypoint.RabbitMq)
}
func ServiceConsumerStart(consumerEntrypoint *ConsumerEntrypoint) {
}
func (consumerEntrypoint *ConsumerEntrypoint) ConsumerEntrypointStart() {
	ControllerConsumerStart(consumerEntrypoint)
	ServiceConsumerStart(consumerEntrypoint)
}
