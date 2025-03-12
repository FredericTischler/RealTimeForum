package config

import (
	"forum/handlers"
	"forum/services"
	"net/http"
)

func InitializeRoutes(mux *http.ServeMux, userService *services.UserService, authService *services.AuthService, postsService *services.PostsService, commentService *services.CommentsService, sessionService *services.SessionService) {
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))
	mux.HandleFunc("/", handlers.HomeHandler)

	// Utiliser des closures pour capturer les services et les passer aux gestionnaires de routes
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, userService)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, authService)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutHandler(w, r, sessionService)
	})

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.PostsHandler(w, r, postsService, sessionService)
		} else if r.Method == http.MethodGet {
			handlers.GetPostsHandler(w, r, postsService, userService)
		}
	})
	mux.HandleFunc("/auth/status", func(w http.ResponseWriter, r *http.Request) {
		handlers.AuthStatusHandler(w, r, sessionService)
	})
	mux.HandleFunc("/posts/{id}", handlers.GetPostsByIdHandler)
	mux.HandleFunc("/posts/comment/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.PostCommentHandler(w, r, commentService, sessionService)
		} else if r.Method == http.MethodGet {
			handlers.GetCommentHandler(w, r, commentService, userService)
		}
	})
	mux.HandleFunc("/users/connected", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetConnectedUsersHandler(w, r, sessionService)
	})
	mux.HandleFunc("/posts/comments/{postid}", handlers.GetCommentsHandler)
	mux.HandleFunc("/message", handlers.MessageHandler)
	mux.HandleFunc("/message/{id}", handlers.GetMessageHandler)
	mux.HandleFunc("/users", handlers.GetUsersHandler)
	mux.HandleFunc("/users/{id}", handlers.GetUsersHandler)
}
