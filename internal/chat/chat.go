package chat

import (
	"go-chat/pkg/models"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatRoom struct {
    Name    string
    Members map[*Client]bool
    Mutex   sync.Mutex
}

type Client struct {
    Username string
    Conn     *websocket.Conn
    Room     *ChatRoom
}
// create chat room in table
var ChatRooms = make(map[string]*ChatRoom)
var Clients = make(map[string]*Client)
var Mutex = sync.Mutex{}


func CreateChatRoom(name string) {
    // Create chatroom in DB
    chatRoom := &models.ChatRoom{Name: name}
    chatRoom.CreateChatRoom()
}

// JoinChatRoom adds a client to a chat room
// func JoinChatRoom(roomName string, client *Client) {
//     room, exists := ChatRooms[roomName]
//     if !exists {
//         room = CreateChatRoom(roomName)
//     }
//     room.Mutex.Lock()
//     room.Members[client] = true
//     room.Mutex.Unlock()
//     client.Room = room
// }

// LeaveChatRoom removes a client from a chat room
func LeaveChatRoom(roomName string, client *Client) {
    room, exists := ChatRooms[roomName]
    if exists {
        room.Mutex.Lock()
        delete(room.Members, client)
        room.Mutex.Unlock()
        client.Room = nil
    }
}

// BroadcastMessage sends a message to all members of a chat room
func BroadcastMessage(roomName, message string) {
    room, exists := ChatRooms[roomName]
    if exists {
        room.Mutex.Lock()
        for member := range room.Members {
            member.Conn.WriteMessage(websocket.TextMessage, []byte(message))
        }
        room.Mutex.Unlock()
    }
}

// SendDirectMessage sends a direct message to a specific user
func SendDirectMessage(username string, message string) {
    Mutex.Lock()
    client, exists := Clients[username]
    Mutex.Unlock()
    if exists {
        client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
    }
}
