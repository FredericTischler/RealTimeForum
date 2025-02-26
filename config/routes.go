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

}
