package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Post struct {
	PostId    uuid.UUID
	UserId    uuid.UUID
	Title     string
	Content   string
	Category  string
	CreatedAt time.Time
}

type PostPayload struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type PostResponse struct {
	response *ResponseBody
	post     *Post
}

type GetPostsResponse struct {
	response *ResponseBody
	post     []*Post
}

type GetPostResponse struct {
	response *ResponseBody
	post     *Post
}

type PostAuthor struct {
	Post     Post
	Username string
}
