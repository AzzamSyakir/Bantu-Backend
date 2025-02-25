package routes

import (
	"bantu-backend/src/internal/controllers"
	"fmt"
	"log"
	"net/url"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

type Socket struct {
	Router         *mux.Router
	ChatController *controllers.ChatController
}

func NewSocket(
	router *mux.Router,
	chatController *controllers.ChatController,
) *Socket {
	return &Socket{
		Router:         router,
		ChatController: chatController,
	}
}

func (socket *Socket) RegisterSocket() *socketio.Server {
	server := socketio.NewServer(nil)
	log.Println("Socket.IO server started")

	server.OnConnect("/", func(s socketio.Conn) error {
		rawQuery := s.URL().RawQuery
		params, err := url.ParseQuery(rawQuery)
		if err != nil {
			return fmt.Errorf("Failed to parse query params: %v", err)
		}

		senderID := params.Get("sender_id")
		receiverID := params.Get("receiver_id")

		if senderID == "" || receiverID == "" {
			return fmt.Errorf("Missing sender_id or receiver_id")
		}
		room, _ := socket.ChatController.GetOrCreateRoom(senderID, receiverID)
		log.Println("User connected:", s.ID())
		s.Join(room)
		return nil
	})

	server.OnEvent("/", "chat_message", func(s socketio.Conn, data map[string]string) {
		senderID := data["sender_id"]
		receiverID := data["receiver_id"]

		roomID, err := socket.ChatController.GetOrCreateRoom(data["sender_id"], data["receiver_id"])
		if err != nil {
			log.Println("Invalid room_id:", data["room_id"])
			return
		}

		message := data["message"]
		if message == "" {
			log.Println("Empty message received")
			return
		}

		err = socket.ChatController.CreateChat(roomID, senderID, receiverID, message)
		if err != nil {
			fmt.Println(err.Error())
			return

		}
		s.Emit("chat_response", map[string]string{
			"status":  "success",
			"message": "Message success send: " + message,
		})

		server.BroadcastToRoom("/", roomID, "chat_message", map[string]string{
			"sender_id": senderID,
			"message":   message,
		})

		log.Println("Message sent to room:", roomID, "from:", senderID)
	})

	server.OnEvent("/", "leave_room", func(s socketio.Conn, data map[string]string) {
		room, _ := socket.ChatController.GetOrCreateRoom(data["sender_id"], data["receiver_id"])
		roomStr := room
		s.Leave(roomStr)
		log.Println("User left room:", roomStr)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server failed: %v", err)
		}
	}()
	return server
}
