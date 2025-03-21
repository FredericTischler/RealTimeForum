package handlers

import (
	"forum/services"
	"net/http"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {

}

func GetPrivateMessageHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService, messageService *services.MessageService) {

}
