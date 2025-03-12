package services

import (
	"forum/models"
	"forum/repositories"
	"time"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	SessionRepo *repositories.SessionRepository
	UserService *UserService
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

func (ss *SessionService) GetActiveSessions() ([]models.Session, error) {
	query := `SELECT * FROM sessions WHERE expires_at > ?`
	rows, err := ss.SessionRepo.DB.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(&session.SessionId, &session.UserId, &session.Token, &session.CreatedAt, &session.ExpireAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (ss *SessionService) GetUserBySession(session models.Session) (*models.User, error) {
	// Utilisez UserService pour récupérer l'utilisateur par son UUID
	user, err := ss.UserService.GetUserByUUID(session.UserId.String())
	if err != nil {
		return nil, err
	}
	return user, nil
}
