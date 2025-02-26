package config

import (
	"forum/handlers"
	"net/http"
)

func InitializeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("POST /register", handlers.RegisterHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)
	mux.HandleFunc("POST /posts", handlers.PostsHandler)
	mux.HandleFunc("GET /posts", handlers.GetPostsHandler)
	mux.HandleFunc("GET /posts/{id}", handlers.GetPostsByIdHandler)
	mux.HandleFunc("POST /posts/comment/{id}", handlers.PostCommentHandler)
	mux.HandleFunc("GET /posts/comments/{postid}", handlers.GetCommentsHandler)
	mux.HandleFunc("POST /message", handlers.MessageHandler)
	mux.HandleFunc("GET /message/{id}", handlers.GetMessageHandler)
	mux.HandleFunc("GET /users", handlers.GetUsersHandler)
	mux.HandleFunc("GET /users/{id}", handlers.GetUsersHandler)
}
