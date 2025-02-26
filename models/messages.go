package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Message struct {
	messageId  uuid.UUID
	senderId   uuid.UUID
	receiverId uuid.UUID
	content    string
	sentAt     time.Time
	isRead     bool
}
