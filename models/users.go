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
