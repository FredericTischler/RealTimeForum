package handlers

import (
	"fmt"
	"forum/models"
	"forum/services"
	"forum/utils"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebsocketHandler gère les connexions WebSocket
func WebsocketHandler(hub *utils.Hub, sessionService *services.SessionService, userService *services.UserService, messageService *services.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade la connexion HTTP en WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		// Récupérer l'utilisateur connecté
		user := getUsernameFromRequest(r, sessionService, userService)
		if user == nil {
			conn.Close()
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Ajouter la connexion au Hub
		hub.AddClient(conn, user)
		hub.BroadcastUsers()

		// Boucle pour lire les messages entrants
		go func() {
			defer func() {
				hub.RemoveClient(conn)
				conn.Close()
			}()

			for {
				var msg models.Message
				err := conn.ReadJSON(&msg)
				if err != nil {
					break
				}
				msg.MessageId, _ = uuid.NewV4()
				msg.SenderId = user.UserId
				msg.SentAt = time.Now()
				// Enregistrer le message en base de données
				err = messageService.InsertMessage(&msg)
				if err != nil {
					fmt.Println("Erreur lors de l'enregistrement du message:", err)
					continue
				}

				// Diffuser le message aux utilisateurs concernés
				hub.BroadcastMessage(msg.ReceiverId, &msg)
				hub.BroadcastMessage(msg.SenderId, &msg) // Envoyer aussi au sender
			}
		}()
	}
}

// getUsernameFromRequest récupère le nom d'utilisateur à partir de la session et du UserService
func getUsernameFromRequest(r *http.Request, sessionService *services.SessionService, userService *services.UserService) *models.UserList {
	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	// Récupérer la session associée au token
	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		return nil
	}
	// Récupérer l'utilisateur via son ID présent dans la session
	user, err := userService.GetUsernameAndAgeAndGenderByUUID(session.UserId.String())
	if err != nil {
		return nil
	}
	return user
}
