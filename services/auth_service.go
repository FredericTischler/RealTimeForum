package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserMethod struct {
	DB *sql.DB
}

func (u *UserMethod) InsertInUser(userName, email, password, firstName, lastName, gender string, age int, userUUID uuid.UUID) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	createdAt := time.Now()

	query := `INSERT INTO users (uuid, user_name, email, password, first_name, last_name, gender, age, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := u.DB.Exec(query, userUUID.String(), userName, email, passwordHash, firstName, lastName, gender, age, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return userID, nil
}
