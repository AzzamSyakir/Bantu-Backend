package controllers

import (
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
)

type ChatController struct {
	ChatService     *services.ChatService
	ResponseChannel *response.ResponseChannel
}

func NewChatController(authService *services.ChatService, responseChannel *response.ResponseChannel) *ChatController {
	return &ChatController{
		ChatService:     authService,
		ResponseChannel: responseChannel,
	}
}
