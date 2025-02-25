package routes

import (
	"bantu-backend/src/internal/controllers"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Router         *mux.Router
	ChatController *controllers.ChatController
	Clients        map[string]*websocket.Conn
	Rooms          map[string]map[*websocket.Conn]bool
	Mutex          sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWebSocketServer(router *mux.Router, chatController *controllers.ChatController) *WebSocketServer {
	return &WebSocketServer{
		Router:         router,
		ChatController: chatController,
		Clients:        make(map[string]*websocket.Conn),
		Rooms:          make(map[string]map[*websocket.Conn]bool),
	}
}

func (ws *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	senderID := r.URL.Query().Get("sender_id")
	receiverID := r.URL.Query().Get("receiver_id")

	if senderID == "" || receiverID == "" {
		log.Println("Missing sender_id or receiver_id")
		return
	}

	roomID, err := ws.ChatController.GetOrCreateRoom(senderID, receiverID)
	if err != nil {
		log.Println("Error creating room:", err)
		return
	}

	ws.Mutex.Lock()
	if ws.Rooms[roomID] == nil {
		ws.Rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	ws.Rooms[roomID][conn] = true
	ws.Mutex.Unlock()

	log.Println("User connected:", senderID, "to room:", roomID)

	for {
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		senderID := msg["sender_id"]
		receiverID := msg["receiver_id"]
		message := msg["message"]

		if senderID == "" || receiverID == "" || message == "" {
			log.Println("Invalid message data")
			continue
		}

		err = ws.ChatController.CreateChat(roomID, senderID, receiverID, message)
		if err != nil {
			log.Println("Failed to save message:", err)
			continue
		}

		ws.BroadcastMessage(roomID, map[string]string{
			"sender_id": senderID,
			"message":   message,
		})
	}

	ws.Mutex.Lock()
	delete(ws.Rooms[roomID], conn)
	ws.Mutex.Unlock()
}

func (ws *WebSocketServer) BroadcastMessage(roomID string, message map[string]string) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()

	for client := range ws.Rooms[roomID] {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Close()
			delete(ws.Rooms[roomID], client)
		}
	}
}
