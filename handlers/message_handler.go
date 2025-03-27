package handlers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"forum/services"
	"net/http"
	"strconv"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {

}

func GetPrivateMessageHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, messageService *services.MessageService) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	withUserId := r.URL.Query().Get("with")
	offsetStr := r.URL.Query().Get("offset")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	messages, err := messageService.GetMessage(session.UserId.String(), withUserId, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func MarkMessagesAsReadHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, messageService *services.MessageService) {
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

	// Récupérer l'ID de l'expéditeur (l'utilisateur avec qui je discute)
	senderId := r.URL.Query().Get("sender_id")
	if senderId == "" {
		http.Error(w, "sender_id is required", http.StatusBadRequest)
		return
	}

	// Marquer les messages comme lus
	err = messageService.MarkMessagesAsRead(senderId, session.UserId.String())
	if err != nil {
		http.Error(w, "Failed to mark messages as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
