package services

import (
	"forum/models"
	"forum/repositories"

	"github.com/gofrs/uuid"
)

type PostsService struct {
	PostsRepo *repositories.PostsRepository
}

// InsertPost appelle le repository pour insérer un post.
func (ps *PostsService) InsertPost(userId uuid.UUID, title, content, category string, postUUID uuid.UUID) (int64, error) {
	return ps.PostsRepo.InsertPost(userId, title, content, category, postUUID)
}

// GetPostById retourne un post en fonction de son UUID.
func (ps *PostsService) GetPostById(postUUID uuid.UUID) (*models.Post, error) {
	return ps.PostsRepo.GetPostById(postUUID)
}

// GetPosts retourne la liste de tous les posts.
func (ps *PostsService) GetPosts(limit, offset int, category, keyword, author string) ([]*models.Post, error) {
	return ps.PostsRepo.GetPosts(limit, offset, category, keyword, author)
}
