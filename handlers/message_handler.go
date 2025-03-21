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
		ErrorHandler(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, err.Error())
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
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
