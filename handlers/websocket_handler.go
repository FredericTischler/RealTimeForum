package handlers

import (
	"forum/models"
	"forum/services"
	"forum/utils"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configuration de l'upgrader pour autoriser les connexions depuis n'importe quelle origine
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebsocketHandler upgrade la connexion HTTP en websocket et l'enregistre dans le hub
// Il reçoit maintenant en paramètres sessionService et userService pour récupérer le nom d'utilisateur.
func WebsocketHandler(hub *utils.Hub, sessionService *services.SessionService, userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getUsernameFromRequest(r, sessionService, userService)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// Ajout de la connexion au hub
		hub.AddClient(conn, user)
		hub.BroadcastUsers()

		// Gérer les messages WebSocket
		go func() {
			defer func() {
				// Déconnexion de l'utilisateur
				hub.RemoveClient(conn)
				hub.BroadcastUsers()
				conn.Close()
			}()

			for {
				var msg struct {
					Type       string `json:"type"`
					ReceiverID string `json:"receiveruuid"`
					Content    string `json:"content"`
				}
				err := conn.ReadJSON(&msg)
				if err != nil {
					break
				}

				if msg.Type == "message" {
					// Enregistrer le message dans la base de données
					message, err := utils.SendMessage(user.UserId, msg.ReceiverID, msg.Content)
					if err != nil {
						log.Println("Failed to save message:", err)
						continue
					}

					// Diffuser le message au destinataire
					hub.BroadcastMessage(msg.ReceiverID, message)
				}
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

	// Récupérer la session associée
	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil || session == nil {
		return nil
	}

	// Vérifie si l'UserID est bien présent
	if session.UserId.String() == "" {
		return nil
	}

	// Vérifie si l'utilisateur existe dans la base de données
	user, err := userService.GetUsernameAndAgeAndGenderByUUID(session.UserId.String())
	if err != nil || user == nil {
		return nil
	}

	return user
}
