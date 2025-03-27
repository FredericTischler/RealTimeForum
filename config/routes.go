package config

import (
	"forum/handlers"
	"forum/services"
	"forum/utils" // Assurez-vous que ce package contient votre struct Hub et sa méthode NewHub()
	"net/http"
)

func InitializeRoutes(mux *http.ServeMux, userService *services.UserService, authService *services.AuthService, postsService *services.PostsService, commentService *services.CommentsService, sessionService *services.SessionService, messageService *services.MessageService) {
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
		handlers.AuthStatusHandler(w, r, sessionService, userService)
	})
	mux.HandleFunc("/posts/{id}", handlers.GetPostsByIdHandler)
	mux.HandleFunc("/posts/comment/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.PostCommentHandler(w, r, commentService, sessionService)
		} else if r.Method == http.MethodGet {
			handlers.GetCommentHandler(w, r, commentService, userService)
		}
	})
	mux.HandleFunc("/posts/comments/{postid}", handlers.GetCommentsHandler)
	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPrivateMessageHandler(w, r, sessionService, messageService)
	})
	mux.HandleFunc("/message/insert", func(w http.ResponseWriter, r *http.Request) {
		handlers.InsertMessageHandler(w, r, sessionService, messageService)
	})
	mux.HandleFunc("/message/unread", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUnreadMessagesCountHandler(w, r, sessionService, messageService)
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUsersHandler(w, r, userService)
	})
	mux.HandleFunc("/users/{id}", handlers.GetUsersByIdHandler)

	// Création du hub pour les WebSockets et ajout de la route dédiée
	hub := utils.NewHub()
	mux.HandleFunc("/ws", handlers.WebsocketHandler(hub, sessionService, userService))
	mux.HandleFunc("/ws/private", func(w http.ResponseWriter, r *http.Request) {
		handlers.PrivateMessageHandler(w, r, sessionService)
	})

}
