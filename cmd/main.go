package main

import (
	"fmt"
	"forum/config"
	"forum/repositories"
	"forum/services"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	config.CreateDataBase()
	// Initialiser la base de données
	db, err := services.ConnectDB()
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

	// Initialiser les routes
	config.InitializeRoutes(mux, userService, authService, postService, sessionService)

	fmt.Println("Starting server...")
	config.StartServer(mux)
}
