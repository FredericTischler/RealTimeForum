package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	MessageId  uuid.UUID
	SenderId   uuid.UUID
	ReceiverId uuid.UUID
	Content    string
	SentAt     time.Time
}

type MessageResponse struct {
	response *ResponseBody
	message  *Message
}

type GetMessagesResponse struct {
	response *ResponseBody
	messages []*Message
}
