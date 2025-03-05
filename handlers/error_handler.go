package handlers

import (
	"encoding/json"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": status,
		"error":  message,
	})
}
