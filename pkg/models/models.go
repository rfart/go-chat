package models

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

// func OpenDB() {
//     var err error
// 	DB, err = gorm.Open(postgres.Open("postgresql://postgres:password@localhost:5432/chatapp"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// }

type User struct {
    gorm.Model
    ID           int
    Username     string
    Password     string
    CreatedAt    time.Time      // Automatically managed by GORM for creation time
    UpdatedAt    time.Time 
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Message struct {
    gorm.Model
    ID      int
    Sender  string
    Room    string
    Message string
}

type ChatRoom struct {
	gorm.Model
	Name string `json:"Name"`
    Members []User `gorm:"many2many:chat_room_users;"`
}

func (u *User) CreateUser() {
	DB.Create(&u)
}

func (c *ChatRoom) CreateChatRoom() {
    DB.Create(&c)
}