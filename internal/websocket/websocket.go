package websocket

import (
	"go-chat/internal/chat"
	"go-chat/pkg/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Message struct {
    Type    string `json:"type"`
    Room    string `json:"room"`
    Sender  string `json:"sender"`
    Target  string `json:"target,omitempty"`
    Content string `json:"content"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch msg.Type {
		// case "join":
		// 	handleJoin(conn, msg)
		case "leave":
			handleLeave(msg)
		case "broadcast":
			handleBroadcast(msg)
		case "dm":
			handleDirectMessage(msg)
		}
	}
}

// func handleJoin(conn *websocket.Conn, msg Message) {
//     chat.Clients[msg.Sender] = &chat.Client{
// 		Username : msg.Sender,
// 		Conn:     conn,
// 	}
// 	chat.JoinChatRoom(msg.Room, chat.Clients[msg.Sender])
// 	chat.BroadcastMessage(msg.Room, msg.Sender+" joined the room")
// 	storeMessage(msg.Sender, msg.Room, msg.Sender+" joined the room")
// }

func handleLeave(msg Message) {
    chat.LeaveChatRoom(msg.Room, chat.Clients[msg.Sender])
	chat.BroadcastMessage(msg.Room, msg.Sender+" left the room")
	storeMessage(msg.Sender, msg.Room, msg.Sender+" left the room")
}

func handleBroadcast(msg Message) {
    chat.BroadcastMessage(msg.Room, msg.Sender+": "+msg.Content)
	storeMessage(msg.Sender, msg.Room, msg.Content)
}

func handleDirectMessage(msg Message) {
    chat.Clients[msg.Target].Conn.WriteJSON(msg)
	storeMessage(msg.Sender, msg.Room, msg.Content)
}

func storeMessage(sender string, room string, message string) {
    err := models.DB.Create(&models.Message{
		Sender:  sender,
		Room:    room,
		Message: message,
	}).Error

	if err != nil {
		log.Printf("Error storing message: %v", err)
	}
}
