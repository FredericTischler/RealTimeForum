package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Session struct {
	SessionId uuid.UUID
	UserId    uuid.UUID
	Token     string
	CreatedAt time.Time
	ExpireAt  time.Time
}

type ActiveSession struct {
	SessionId int       `json:"id"`
	UUID      string    `json:"uuid"`
	UserId    string    `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpireAt  time.Time `json:"expires_at"`
}
