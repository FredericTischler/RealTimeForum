package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	MessageID  uuid.UUID `json:"messageuuid"`
	SenderID   string    `json:"senderuuid"`
	ReceiverID string    `json:"receiveruuid"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"timestamp"`
	IsRead     bool      `json:"is_read"`
}

type MessageResponse struct {
	response *ResponseBody
	message  *Message
}

type GetMessagesResponse struct {
	response *ResponseBody
	messages []*Message
}
