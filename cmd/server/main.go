package main

import (
	"go-chat/internal/auth"
	"go-chat/internal/websocket"
	"go-chat/pkg/models"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
    var err error
    models.DB, err = gorm.Open(postgres.Open("postgresql://postgres:password@localhost:5432/chatapp"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	models.DB.AutoMigrate(models.User{}, models.Message{}, models.ChatRoom{})

    http.HandleFunc("/ws", websocket.HandleConnections)
    http.HandleFunc("/register", auth.Register)
    http.HandleFunc("/login", auth.Login)

    log.Println("Server started on :8080")
    e := http.ListenAndServe(":8080", nil)
    if e != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
