package handlers

import (
	"forum/services"
	"net/http"
)

func PostsHandler(w http.ResponseWriter, r *http.Request, postService *services.PostsService) {
	// Implement the logic to handle POST requests for posts
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request, postService *services.PostsService) {
	// Implement the logic to handle GET requests for posts
}

func GetPostsByIdHandler(w http.ResponseWriter, r *http.Request) {

}
