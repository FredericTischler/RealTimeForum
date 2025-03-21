package handlers

import (
	"encoding/json"
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
