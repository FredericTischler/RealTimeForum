package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	UserId    uuid.UUID
	Username  string
	Email     string
	Password  string
	Firstname string
	Lastname  string
	Age       int
	Gender    string
	CreatedAt time.Time
}

type RegisterResponse struct {
	Response *ResponseBody
	UserId   uuid.UUID
}

type LoginResponse struct {
	Response *ResponseBody
	UserId   uuid.UUID
}

type GetUsersResponse struct {
	response *ResponseBody
	users    []*User
}

type GetUserResponse struct {
	response *ResponseBody
	user     *User
}
