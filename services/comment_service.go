package services

import (
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

// // GetPostById retourne un post en fonction de son UUID.
// func (ps *PostsService) GetCommentsByPost(postUUID uuid.UUID) (*models.Post, error) {
// 	return ps.PostsRepo.GetCommentsByPost(postUUID)
// }

// // GetAllPosts retourne la liste de tous les posts.
// func (ps *PostsService) GetAllPosts() ([]*models.Post, error) {

// 	return ps.PostsRepo.GetAllPosts()
// }
