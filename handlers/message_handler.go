package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"forum/services"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request, messageService *services.MessageService, sessionService *services.SessionService) {
	fmt.Println("handler")
	// Récupérer l'ID de l'utilisateur connecté
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	userID := session.UserId

	url := strings.Split(r.URL.Path, "/")
	receiverID := url[len(url)-1]

	// Récupérer les messages entre les deux utilisateurs
	messages, err := messageService.GetMessagesBetweenUsers(userID.String(), receiverID)
	if err != nil {
		fmt.Println(err.Error())
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Renvoyer les messages en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
