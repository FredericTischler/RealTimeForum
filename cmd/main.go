package main

import (
	"fmt"
	"forum/config"
	"forum/middleware"
	"forum/repositories"
	"forum/services"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	config.CreateDataBase()
	// Initialiser la base de données
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	// Créer les repositories et services
	userRepo := &repositories.UserRepository{DB: db}
	userService := &services.UserService{UserRepo: userRepo}
	sessionRepo := &repositories.SessionRepository{DB: db}
	sessionService := &services.SessionService{SessionRepo: sessionRepo}
	authService := &services.AuthService{
		UserService:    userService,
		SessionService: sessionService,
	}

	postRepo := &repositories.PostsRepository{DB: db}
	postService := &services.PostsService{PostsRepo: postRepo}

	commentRepo := &repositories.CommentsRepositories{DB: db}
	commentService := &services.CommentsService{CommentsRepo: commentRepo}

	// Initialiser les routes
	config.InitializeRoutes(mux, userService, authService, postService, commentService, sessionService)

	rateLimitedHandler := middleware.SessionRateLimiter(sessionService, mux)

	fmt.Println("Starting server...")
	config.StartServer(rateLimitedHandler)
}
