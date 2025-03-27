package models

import (
	"time"
)

type Message struct {
	MessageId  string
	SenderId   string
	ReceiverId string
	Content    string
	SentAt     time.Time
	Is_read    bool
}

type MessageResponse struct {
	response *ResponseBody
	message  *Message
}

type GetMessagesResponse struct {
	response *ResponseBody
	messages []*Message
}
