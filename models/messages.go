package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	MessageId  uuid.UUID
	SenderId   string
	ReceiverId string `json:"receiver_id"`
	Content    string
	SentAt     time.Time
	isRead     bool
}

type MessageResponse struct {
	MessageID  string `json:"message_id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	SentAt     string `json:"sent_at"` // Format ISO 8601
}

// type MessageResponse struct {
// 	response *ResponseBody
// 	message  *Message
// }

type GetMessagesResponse struct {
	response *ResponseBody
	messages []*Message
}
