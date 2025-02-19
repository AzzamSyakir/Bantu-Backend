package services

import (
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
)

type ChatService struct {
	ChatRepository *repository.ChatRepository
	Producer       *producer.ServicesProducer
}

func NewChatService(chatRepo *repository.ChatRepository, producer *producer.ServicesProducer) *ChatService {
	return &ChatService{
		ChatRepository: chatRepo,
	}
}
