package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	CommentId uuid.UUID
	PostId    uuid.UUID
	UserId    uuid.UUID
	Content   string
	CreatedAt time.Time
}

type CommentPayload struct {
	PostId  string `json:"post_id"`
	Content string `json:"content"`
}

type GetCommentsResponse struct {
	response *ResponseBody
	comments []*Comment
}
