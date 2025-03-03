package services

import (
	"errors"
	"forum/utils"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService *UserService
}

func (as *AuthService) Login(identifier, password string) (uuid.UUID, string, error) {
	user, err := as.UserService.GetUserByUsername(identifier)
	if err != nil {
		user, err = as.UserService.GetUserByEmail(identifier)
		if err != nil {
			return uuid.Nil, "", errors.New(err.Error())
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return uuid.Nil, "", errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(user.UserId)
	if err != nil {
		return uuid.Nil, "", err
	}

	return user.UserId, token, nil
}
