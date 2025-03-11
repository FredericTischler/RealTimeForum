package config

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartServer(handler http.Handler) {
	server := &http.Server{
		Addr:              ":8443",
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}
	fmt.Println("Server is running at http://localhost:8443")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
