package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"

	"github.com/gofrs/uuid"
)

type SessionRepository struct {
	DB *sql.DB
}

// InsertSession insère une nouvelle session dans la base de données.
func (sr *SessionRepository) InsertSession(sessionUUID uuid.UUID, userID uuid.UUID, token string, createdAt time.Time, expiresAt *time.Time) (int64, error) {
	query := `INSERT INTO sessions (uuid, user_id, token, created_at, expires_at) VALUES (?, ?, ?, ?, ?)`
	result, err := sr.DB.Exec(query, sessionUUID.String(), userID.String(), token, createdAt, expiresAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert session: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return id, nil
}

// GetSessionByToken récupère une session à partir de son token.
func (sr *SessionRepository) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	query := `SELECT uuid, user_id, token, created_at, expires_at FROM sessions WHERE token = ?`
	err := sr.DB.QueryRow(query, token).Scan(&session.SessionId, &session.UserId, &session.Token, &session.CreatedAt, &session.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSessionByToken supprime une session à partir de son token.
func (sr *SessionRepository) DeleteSessionByToken(token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := sr.DB.Exec(query, token)
	return err
}
