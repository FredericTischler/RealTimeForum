package services

import (
	"forum/models"
	"forum/repositories"

	"github.com/gofrs/uuid"
)

type CommentsService struct {
	CommentsRepo *repositories.CommentsRepositories
}

// InsertPost appelle le repository pour insérer un post.
func (ps *CommentsService) InsertComment(userID uuid.UUID, postID, content string, commentUUID uuid.UUID) (int64, error) {
	return ps.CommentsRepo.InsertComment(userID, postID, content, commentUUID)
}

// GetCommentByPost retourne un comment en fonction de son UUID du post.
func (ps *CommentsService) GetCommentsByPost(postUUID string) ([]models.Comment, error) {
	return ps.CommentsRepo.GetCommentsByPost(postUUID)
}

// // GetAllPosts retourne la liste de tous les posts.
// func (ps *PostsService) GetAllPosts() ([]*models.Post, error) {

// 	return ps.PostsRepo.GetAllPosts()
// }
