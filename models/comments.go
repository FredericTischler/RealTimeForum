package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Comment struct {
	commentId uuid.UUID
	postId    uuid.UUID
	userId    uuid.UUID
	content   string
	createdAt time.Time
}
