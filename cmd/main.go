package main

import (
	"fmt"
	"forum/config"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	config.InitializeRoutes(mux)
	config.CreateDataBase()
	fmt.Println("Starting server...")
	config.StartServer(mux)
}
