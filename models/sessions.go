package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Session struct {
	SessionId uuid.UUID
	UserId    uuid.UUID
	Token     string
	CreatedAt time.Time
	ExpireAt  time.Time
}
