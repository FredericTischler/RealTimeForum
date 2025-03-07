package repositories

import (
	"database/sql"
	"forum/models"
	"time"

	"github.com/gofrs/uuid"
)

type CommentsRepositories struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentsRepositories {
	return &CommentsRepositories{DB: db}
}

func (cr *CommentsRepositories) InsertComment(userID uuid.UUID, postID, content string, commentUUID uuid.UUID) (int64, error) {
	query := `INSERT INTO comments (uuid, content, created_at, user_id, post_id) VALUES (?, ?, ?, ?, ?)`
	result, err := cr.DB.Exec(query, commentUUID.String(), content, time.Now(), userID, postID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (cr *CommentsRepositories) GetCommentsByPost(postID string) ([]models.Comment, error) {
	query := `SELECT uuid, content, created_at, user_id FROM comments WHERE post_id = ? ORDER BY created_at ASC`
	rows, err := cr.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var userID uuid.UUID
		if err := rows.Scan(&comment.CommentId, &comment.Content, &comment.CreatedAt, &userID); err != nil {
			return nil, err
		}
		comment.UserId = userID
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
