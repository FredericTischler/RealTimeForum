package handlers

import (
	"encoding/json"
	"fmt"
	"forum/services"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request, service *services.UserService) {
	users, err := service.GetUsers()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve users: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
	}
}

func GetUsersByIdHandler(w http.ResponseWriter, r *http.Request) {

}
