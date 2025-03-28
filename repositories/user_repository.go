package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func (ur *UserRepository) GetUsers() (*[]models.UserList, error) {
	var users []models.UserList
	query := `SELECT uuid, user_name, age, gender FROM users`
	result, err := ur.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	for result.Next() {
		var user models.UserList
		result.Scan(
			&user.UserId,
			&user.Username,
			&user.Age,
			&user.Gender)
		users = append(users, user)
	}
	return &users, nil
}

func (ur *UserRepository) InsertUser(userName, email, password, firstName, lastName, gender string, age int, userUUID uuid.UUID) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	createdAt := time.Now()

	query := `INSERT INTO users (uuid, user_name, email, password, first_name, last_name, gender, age, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := ur.DB.Exec(query, userUUID.String(), userName, email, passwordHash, firstName, lastName, gender, age, createdAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return userID, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT uuid, user_name, email, password, first_name, last_name, age, gender, created_at FROM users WHERE email = ?"
	err := ur.DB.QueryRow(query, email).Scan(
		&user.UserId,   // Assurez-vous que ce champ est de type uuid.UUID dans models.User
		&user.Username, // Cette valeur correspondra à la colonne user_name
		&user.Email,
		&user.Password,
		&user.Firstname,
		&user.Lastname,
		&user.Age,
		&user.Gender,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT uuid, user_name, email, password, first_name, last_name, age, gender, created_at FROM users WHERE user_name = ?`
	err := ur.DB.QueryRow(query, username).Scan(
		&user.UserId,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Firstname,
		&user.Lastname,
		&user.Age,
		&user.Gender,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByUUID(userID string) (string, error) {
	var user string
	query := `SELECT user_name FROM users WHERE uuid = ?`
	err := ur.DB.QueryRow(query, userID).Scan(&user)
	if err != nil {
		return "", err
	}
	return user, nil
}

func (ur *UserRepository) GetUsernameAndAgeAndGenderByUUID(userID string) (*models.UserList, error) {
	var user models.UserList
	query := `SELECT uuid, user_name, age, gender FROM users WHERE uuid = ?`
	err := ur.DB.QueryRow(query, userID).Scan(
		&user.UserId,
		&user.Username,
		&user.Age,
		&user.Gender,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
