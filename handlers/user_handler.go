package handlers

import (
	"encoding/json"
	"forum/models"
	"forum/services"
	"net/http"
)

func GetConnectedUsersHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService) {
	// Récupérer la liste des sessions actives
	activeSessions, err := sessionService.GetActiveSessions()
	if err != nil {
		http.Error(w, "Failed to retrieve active sessions", http.StatusInternalServerError)
		return
	}

	// Convertir les sessions en une liste d'utilisateurs connectés
	var connectedUsers []*models.User
	for _, session := range activeSessions {
		user, err := sessionService.GetUserBySession(session)
		if err != nil {
			continue // Ignorer les erreurs et continuer avec les autres sessions
		}
		connectedUsers = append(connectedUsers, user)
	}

	// Retourner la liste des utilisateurs connectés en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connectedUsers)
}
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

}

func GetUsersByIdHandler(w http.ResponseWriter, r *http.Request) {

}
