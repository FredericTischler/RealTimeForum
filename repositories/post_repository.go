package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"

	"github.com/gofrs/uuid"
)

type PostsRepository struct {
	DB *sql.DB
}

// InsertPost insère un nouveau post dans la base de données.
func (pr *PostsRepository) InsertPost(userId uuid.UUID, title, content, category string, postUUID uuid.UUID) (int64, error) {
	createdAt := time.Now()

	query := `INSERT INTO posts (uuid, title, content, category, created_at, user_id) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := pr.DB.Exec(query, postUUID.String(), title, content, category, createdAt, userId.String())

	if err != nil {
		return 0, fmt.Errorf("failed to insert post: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}
	return id, nil
}

// GetPostById récupère un post à partir de son UUID.
func (pr *PostsRepository) GetPostById(postUUID uuid.UUID) (*models.Post, error) {
	var post models.Post
	query := `SELECT uuid, user_id, title, content, category, created_at FROM posts WHERE uuid = ?`
	row := pr.DB.QueryRow(query, postUUID.String())
	err := row.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetAllPosts récupère l'ensemble des posts.
func (pr *PostsRepository) GetPosts(limit, offset int, category, keyword string) ([]*models.Post, error) {
	baseQuery := `SELECT uuid, user_id, title, content, category, created_at FROM posts WHERE 1=1`
	args := []interface{}{}

	// Filtrage par catégorie
	if category != "" {
		baseQuery += " AND category = ?"
		args = append(args, category)
	}

	// Filtrage par mot-clé dans le titre ou le contenu
	if keyword != "" {
		baseQuery += " AND (title LIKE ? OR content LIKE ?)"
		keywordParam := "%" + keyword + "%"
		args = append(args, keywordParam, keywordParam)
	}

	// Tri et pagination
	baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := pr.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Category, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	fmt.Println(posts)
	return posts, nil
}
