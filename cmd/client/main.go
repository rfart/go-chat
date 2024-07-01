package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Message struct {
    Type    string `json:"type"`
    Room    string `json:"room,omitempty"`
    Sender  string `json:"sender"`
    Target  string `json:"target,omitempty"`
    Content string `json:"content"`
}

var conn *websocket.Conn
var username string

func main() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("Enter username:")
    username, _ = reader.ReadString('\n')

    connectToWebSocket()

    for {
        fmt.Println("Enter command:")
        command, _ := reader.ReadString('\n')

        switch command[:len(command)-1] {
        case "register":
            registerUser(reader)
        case "login":
            loginUser(reader)
        case "join":
            fmt.Println("Enter room name:")
            room, _ := reader.ReadString('\n')
            joinRoom(room[:len(room)-1])
        case "send":
            fmt.Println("Enter room name:")
            room, _ := reader.ReadString('\n')
            fmt.Println("Enter message:")
            content, _ := reader.ReadString('\n')
            sendMessage(room[:len(room)-1], content[:len(content)-1])
        case "dm":
            fmt.Println("Enter recipient:")
            target, _ := reader.ReadString('\n')
            fmt.Println("Enter message:")
            content, _ := reader.ReadString('\n')
            sendDirectMessage(target[:len(target)-1], content[:len(content)-1])
        case "exit":
            os.Exit(0)
        }
    }
}

func registerUser(reader *bufio.Reader) {
    fmt.Println("Enter username:")
    username, _ := reader.ReadString('\n')
    fmt.Println("Enter password:")
    password, _ := reader.ReadString('\n')

    user := User{
        Username: username[:len(username)-1], 
        Password: password[:len(password)-1], 
    }

    userData, err := json.Marshal(user)
    if err != nil {
        log.Fatalf("Error marshalling user data: %v", err)
    }

    resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(userData))
    if err != nil {
        log.Fatalf("Error sending register request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusCreated {
        fmt.Println("Registration successful")
    } else {
        fmt.Printf("Registration failed with status: %s\n", resp.Status)
    }
}

func loginUser(reader *bufio.Reader) {
    fmt.Println("Enter username:")
    username, _ := reader.ReadString('\n')
    fmt.Println("Enter password:")
    password, _ := reader.ReadString('\n')

    user := User{
        Username: username[:len(username)-1], 
        Password: password[:len(password)-1], 
    }

    userData, err := json.Marshal(user)
    if err != nil {
        log.Fatalf("Error marshalling user data: %v", err)
    }

    resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(userData))
    if err != nil {
        log.Fatalf("Error sending login request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        fmt.Println("Login successful")
    } else {
        fmt.Printf("Login failed with status: %s\n", resp.Status)
    }
}

func connectToWebSocket() {
    var err error
    conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
    if err != nil {
        log.Fatalf("Failed to connect to WebSocket: %v", err)
    }
}

func joinRoom(room string) {
    msg := Message{
        Type:   "join",
        Room:   room,
        Sender: username,
    }
    conn.WriteJSON(msg)
}

func sendMessage(room, content string) {
    msg := Message{
        Type:    "broadcast",
        Room:    room,
        Sender:  username,
        Content: content,
    }
    conn.WriteJSON(msg)
}

func sendDirectMessage(target, content string) {
    msg := Message{
        Type:    "dm",
        Sender:  username,
        Target:  target,
        Content: content,
    }
    conn.WriteJSON(msg)
}
