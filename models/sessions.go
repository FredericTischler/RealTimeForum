package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Session struct {
	sessionId uuid.UUID
	userId    uuid.UUID
	token     string
	createdAt time.Time
	expireAt  time.Time
}
