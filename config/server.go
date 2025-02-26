package config

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartServer(mux *http.ServeMux) {
	server := &http.Server{
		Addr:              ":8443",
		Handler:           mux,
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
