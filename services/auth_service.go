package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService    *UserService
	SessionService *SessionService
}

// Login vérifie l'identifiant (username ou email) et le mot de passe,
// crée une session dans la base de données et retourne l'UUID de l'utilisateur
// ainsi que le token de session à stocker côté client (via cookie).
func (as *AuthService) Login(identifier, password string) (uuid.UUID, string, error) {
	// Recherche de l'utilisateur par username, puis par email si nécessaire.
	user, err := as.UserService.GetUserByUsername(identifier)
	if err != nil {
		user, err = as.UserService.GetUserByEmail(identifier)
		if err != nil {
			return uuid.Nil, "", errors.New(err.Error())
		}
	}

	// Vérification du mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return uuid.Nil, "", errors.New("invalid password")
	}

	// Génération d'un UUID pour la session (clé primaire de la table sessions)
	sessionUUID, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, "", err
	}

	// Génération d'un token de session à retourner au client (ce token sera stocké dans un cookie)
	tokenUUID, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, "", err
	}

	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour) // Exemple : session valide 24 heures

	// Insertion de la session dans la base de données via le session_service
	_, err = as.SessionService.InsertSession(sessionUUID, user.UserId, tokenUUID.String(), createdAt, &expiresAt)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to create session: %v", err)
	}

	// Retourne l'UUID de l'utilisateur et le token de session
	return user.UserId, tokenUUID.String(), nil
}
