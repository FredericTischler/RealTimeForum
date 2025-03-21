package handlers

import (
	"fmt"
	"forum/models"
	"forum/services"
	"forum/utils"
	"net/http"

	"github.com/gofrs/uuid"
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
		// Upgrade de la connexion HTTP
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		// Récupérer le nom d'utilisateur depuis la session
		user := getUsernameFromRequest(r, sessionService, userService)
		if user == nil {
			conn.Close()
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Ajout de la connexion au hub et diffusion de la liste mise à jour
		hub.AddClient(conn, user)
		hub.BroadcastUsers()

		// Boucle de lecture pour maintenir la connexion ouverte
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}

		// À la déconnexion, retirer la connexion du hub et diffuser la nouvelle liste
		hub.RemoveClient(conn)
		hub.BroadcastUsers()
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

var privateConnections = make(map[uuid.UUID]*websocket.Conn)

func PrivateMessageHandler(w http.ResponseWriter, r *http.Request, sessionService *services.SessionService) {
	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, err.Error())
	}
	// Récupérer la session associée au token
	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
	}

	privateConnections[session.UserId] = conn
	defer conn.Close()

	for {
		var msg struct {
			To      uuid.UUID `json:"to"`
			Message string    `json:"message"`
		}
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		if receiverConn, ok := privateConnections[msg.To]; ok {
			err = receiverConn.WriteJSON(map[string]interface{}{
				"from":    session.UserId,
				"message": msg.Message,
			})

			if err != nil {
				fmt.Println("Failed to send message:", err)
			}
		}
	}
}
