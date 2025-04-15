package models

import (
	"time"
)

// Message represents a message in a chatroom
type Message struct {
	ID         string    `json:"id"`
	ChatroomID string    `json:"chatroom_id"`
	Nickname   string    `json:"nickname"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
