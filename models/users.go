package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	userId    uuid.UUID
	username  string
	email     string
	password  string
	firstname string
	lastname  string
	age       int
	gender    string
	createdAt time.Time
}

type RegisterResponse struct {
	response *ResponseBody
	userId   uuid.UUID
}

type LoginResponse struct {
	response *ResponseBody
	userId   uuid.UUID
	token    string
}

type GetUsersResponse struct {
	response *ResponseBody
	users    []*User
}

type GetUserResponse struct {
	response *ResponseBody
	user     *User
}
