package controllers

import (
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
	"errors"
	"net/http"
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

func (chatController *ChatController) GetChats(writer http.ResponseWriter, reader *http.Request) {
	// chatController.ChatService.GetChatsService(reader)
	if chatController.ChatService == nil {
		http.Error(writer, "ChatService is not initialized", http.StatusInternalServerError)
		return
	}
	chatController.ChatService.GetChatsService(reader)
	select {
	case responseError := <-chatController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-chatController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (chatController *ChatController) GetOrCreateRoom(senderID, receiverID string) (string, error) {
	room, _ := chatController.ChatService.GetOrCreateRoomService(senderID, receiverID)
	return room, nil
}

func (chatController *ChatController) CreateChat(roomID, senderID, receiverID, message string) error {
	chatController.ChatService.CreateChatService(roomID, senderID, receiverID, message)
	select {
	case <-chatController.ResponseChannel.ResponseError:
		return errors.New("create chat is failed")
	case <-chatController.ResponseChannel.ResponseSuccess:
		return nil
	}
}
