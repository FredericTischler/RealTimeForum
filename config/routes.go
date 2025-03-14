package config

import (
	"encoding/json"
	"fmt"
	"forum/handlers"
	"forum/services"
	"forum/utils"
	"log"
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

	// Création du hub WebSocket
	hub := utils.NewHub()

	mux.HandleFunc("/ws", handlers.WebsocketHandler(hub, sessionService, userService))

	mux.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		sender := r.URL.Query().Get("sender")
		fmt.Println(sender)
		recipient := r.URL.Query().Get("recipient")

		log.Printf("Requête reçue - sender: %s, recipient: %s", sender, recipient)

		if sender == "" || recipient == "" {
			log.Println("Paramètres manquants")
			http.Error(w, "Paramètres manquants", http.StatusBadRequest)
			return
		}

		messages, err := utils.GetMessages(sender, recipient)
		if err != nil {
			log.Println("Erreur lors de la récupération des messages :", err)
			http.Error(w, "Erreur lors de la récupération des messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	})

	mux.HandleFunc("/posts/comments/{postid}", handlers.GetCommentsHandler)
	mux.HandleFunc("/message", handlers.MessageHandler)
	mux.HandleFunc("/message/{id}", handlers.GetMessageHandler)
	mux.HandleFunc("/users", handlers.GetUsersHandler)
	mux.HandleFunc("/users/{id}", handlers.GetUsersHandler)
}
