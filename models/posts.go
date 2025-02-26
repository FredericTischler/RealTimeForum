package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Post struct {
	postId    uuid.UUID
	userId    uuid.UUID
	title     string
	content   string
	category  string
	createdAt time.Time
}
