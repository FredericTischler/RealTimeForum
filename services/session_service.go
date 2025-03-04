package services

import (
	"forum/models"
	"forum/repositories"
	"time"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	SessionRepo *repositories.SessionRepository
}

// InsertSession délègue l'insertion d'une session au repository.
func (ss *SessionService) InsertSession(sessionUUID uuid.UUID, userID uuid.UUID, token string, createdAt time.Time, expiresAt *time.Time) (int64, error) {
	return ss.SessionRepo.InsertSession(sessionUUID, userID, token, createdAt, expiresAt)
}

// GetSessionByToken retourne la session associée au token fourni.
func (ss *SessionService) GetSessionByToken(token string) (*models.Session, error) {
	return ss.SessionRepo.GetSessionByToken(token)
}

// DeleteSessionByToken supprime la session correspondant au token.
func (ss *SessionService) DeleteSessionByToken(token string) error {
	return ss.SessionRepo.DeleteSessionByToken(token)
}
