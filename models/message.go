package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Message struct {
	Id              gocql.UUID `json:"id"`
	Type            string     `json:"type"`
	ConversationKey string     `json:"conversation_key"`
	SenderId        uint       `json:"sender_id"`
	ReceiverId      uint       `json:"receiver_id"`
	Content         string     `json:"content"`
	SentAt          time.Time  `json:"sent_at"`
	ReadAt          bool       `json:"read_at"`
}

type IncomingMessage struct {
	SenderId   uint   `json:"sender_id"`
	ReceiverId uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

type Error struct {
	Type    string `json:"type"`
	Content string `json:"message"`
}
