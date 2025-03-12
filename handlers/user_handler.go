package handlers

import (
	"encoding/json"
	"forum/services"
	"log"
	"net/http"
)

func GetConnectedUsersHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, userService *services.UserService) {
	// Log pour suivre l'exécution
	log.Println("Fetching connected users...")

	// Récupérer la liste des sessions actives
	activeSessions, err := sessionService.GetActiveSessions()
	if err != nil {
		log.Printf("Error fetching active sessions: %v\n", err)
		http.Error(w, "Failed to retrieve active sessions", http.StatusInternalServerError)
		return
	}

	// Log pour afficher le nombre de sessions actives
	log.Printf("Found %d active sessions\n", len(activeSessions))

	// Convertir les sessions en une liste d'utilisateurs connectés
	var connectedUsers []string
	for _, session := range activeSessions {
		username, err := userService.GetUserByUUID(session.UserId)
		if err != nil {
			log.Printf("Error fetching username for user %s: %v\n", session.UserId, err)
			continue
		}
		connectedUsers = append(connectedUsers, username.Username)
	}

	// Log pour afficher le nombre d'utilisateurs connectés
	log.Printf("Found %d connected users\n", len(connectedUsers))

	// Retourner la liste des utilisateurs connectés en JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(connectedUsers); err != nil {
		log.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

}

func GetUsersByIdHandler(w http.ResponseWriter, r *http.Request) {

}
