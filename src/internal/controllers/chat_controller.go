package controllers

import "bantu-backend/src/internal/services"

type ChatController struct {
	ChatService *services.ChatService
}

func NewChatController(authService *services.ChatService) *ChatController {
	return &ChatController{
		ChatService: authService,
	}
}
