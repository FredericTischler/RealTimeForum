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
