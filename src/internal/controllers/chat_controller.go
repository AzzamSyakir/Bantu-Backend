package controllers

import (
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
)

type ChatController struct {
	ChatService     *services.ChatService
	ResponseChannel chan response.Response[any]
}

func NewChatController(authService *services.ChatService) *ChatController {
	return &ChatController{
		ChatService:     authService,
		ResponseChannel: make(chan response.Response[any], 1),
	}
}
