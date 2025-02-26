package config

import (
	"forum/handlers"
	"net/http"
)

func InitializeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.HomeHandler)
}
