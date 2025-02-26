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
	status  int
	message string
	userId  uuid.UUID
}

type LoginResponse struct {
	status  int
	message string
	userId  uuid.UUID
	token   string
}

type LogoutResponse struct {
	status  int
	message string
}
