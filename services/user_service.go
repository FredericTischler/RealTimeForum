package services

import (
	"forum/models"
	"forum/repositories"

	"github.com/gofrs/uuid"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func (us *UserService) GetUsers() (*[]models.UserList, error) {
	return us.UserRepo.GetUsers()
}

func (us *UserService) InsertUser(userName, email, password, firstName, lastName, gender string, age int, userUUID uuid.UUID) (int64, error) {
	return us.UserRepo.InsertUser(userName, email, password, firstName, lastName, gender, age, userUUID)
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	return us.UserRepo.GetUserByEmail(email)
}

func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	return us.UserRepo.GetUserByUsername(username)
}

func (us *UserService) GetUserByUUID(userID string) (string, error) {
	return us.UserRepo.GetUserByUUID(userID)
}

func (us *UserService) GetUsernameAndAgeAndGenderByUUID(userID string) (*models.UserList, error) {
	return us.UserRepo.GetUsernameAndAgeAndGenderByUUID(userID)
}
