package services

import (
	"forum/models"
	"forum/repositories"

	"github.com/gofrs/uuid"
)

type CommentsService struct {
	CommentsRepo *repositories.CommentsRepositories
}

// InsertComment appelle le repository pour insérer un comment.
func (ps *CommentsService) InsertComment(userID uuid.UUID, postID, content string, commentUUID uuid.UUID) (int64, error) {
	return ps.CommentsRepo.InsertComment(userID, postID, content, commentUUID)
}

// GetCommentByPost retourne un comment en fonction de l'UUID du post.
func (ps *CommentsService) GetCommentsByPost(postUUID string) ([]models.Comment, error) {
	return ps.CommentsRepo.GetCommentsByPost(postUUID)
}
