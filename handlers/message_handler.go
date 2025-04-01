package handlers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"forum/services"
	"log"
	"net/http"
	"strconv"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {

}

func GetPrivateMessageHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, messageService *services.MessageService) {
	// Récupérer les paramètres de la requête
	queryParams := r.URL.Query()
	withUserId := queryParams.Get("with")
	offsetStr := queryParams.Get("offset")
	onlyUnread := queryParams.Get("unread") == "true" // Nouveau paramètre

	// Convertir offset en int
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Récupérer l'ID de l'utilisateur depuis la session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Récupérer les messages avec le nouveau filtre
	messages, err := messageService.GetMessages(session.UserId.String(), withUserId, offset, onlyUnread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marquer comme lus si nécessaire
	if onlyUnread {
		err = messageService.MarkMessagesAsRead(withUserId, session.UserId.String())
		if err != nil {
			log.Printf("Failed to mark messages as read: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func InsertMessageHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, messageService *services.MessageService) {
	// Vérifier la session de l'utilisateur
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Décoder le corps de la requête pour obtenir les données du message
	var message models.Message
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Valider les données du message
	if message.ReceiverId == "" || message.Content == "" {
		http.Error(w, "Receiver ID and message content are required", http.StatusBadRequest)
		return
	}

	// Définir l'expéditeur du message comme l'utilisateur actuellement connecté
	message.SenderId = session.UserId.String()

	// Enregistrer le message dans la base de données
	err = messageService.InsertMessage(&message)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert message: %v", err), http.StatusInternalServerError)
		return
	}

	// Répondre avec un statut 201 Created
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func GetUnreadMessagesCountHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, messageService *services.MessageService) {
	// Vérifier la session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := messageService.GetUnreadMessagesCount(session.UserId.String())
	if err != nil {
		http.Error(w, "Failed to get unread messages count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
