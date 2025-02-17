package services

import "bantu-backend/src/internal/repository"

type ChatService struct {
	ChatRepository *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{
		ChatRepository: chatRepo,
	}
}
