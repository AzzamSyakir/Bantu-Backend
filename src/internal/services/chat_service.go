package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ChatService struct {
	ChatRepository *repository.ChatRepository
	RabbitMq       *configs.RabbitMqConfig
	Producer       *producer.ServicesProducer
}

func NewChatService(chatRepo *repository.ChatRepository, producer *producer.ServicesProducer, rabbitMq *configs.RabbitMqConfig) *ChatService {
	return &ChatService{
		ChatRepository: chatRepo,
		RabbitMq:       rabbitMq,
		Producer:       producer,
	}
}

func (chatService *ChatService) GetChatsService(reader *http.Request) error {
	queryParams := reader.URL.Query()
	var senderID, receiverID string
	s, sok := queryParams["sender_id"]
	r, rok := queryParams["receiver_id"]
	if !sok || len(s) == 0 || !rok || len(r) == 0 {
		return errors.New("no chats available")
	}
	senderID = s[0]
	receiverID = r[0]
	getChat, err := chatService.ChatRepository.GetChatsRepository(senderID, receiverID)
	if err != nil {
		return chatService.Producer.CreateMessageError(chatService.RabbitMq.Channel, "get chat is failed", http.StatusBadRequest)
	}
	return chatService.Producer.CreateMessageChat(chatService.RabbitMq.Channel, getChat)
}

func (chatService *ChatService) GetOrCreateRoomService(senderID, receiverID string) (string, error) {

	chat, err := chatService.ChatRepository.GetOrCreateRoomRepository(senderID, receiverID)
	fmt.Println(chat)

	if err != nil {
		return "", err
	}

	return chat, nil
}

func (chatService *ChatService) CreateChatService(roomID, senderID, receiverID, message string) error {
	chat := entity.ChatEntity{
		ID:         uuid.NewString(),
		RoomID:     roomID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
	}
	getJobsByRedis, err := chatService.ChatRepository.CreateChatRepository(&chat)
	if err != nil {
		return chatService.Producer.CreateMessageError(chatService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}
	return chatService.Producer.CreateMessageChat(chatService.RabbitMq.Channel, getJobsByRedis)
}
