package handlers

import (
	"forum/services"
	"forum/utils"
	"github.com/gorilla/websocket"
	"net/http"
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
		userID := getUsernameFromRequest(r, sessionService, userService)
		if userID == "" {
			conn.Close()
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Ajout de la connexion au hub et diffusion de la liste mise à jour
		hub.AddClient(conn, userID)
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
func getUsernameFromRequest(r *http.Request, sessionService *services.SessionService, userService *services.UserService) string {
	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	// Récupérer la session associée au token
	session, err := sessionService.GetSessionByToken(cookie.Value)
	if err != nil {
		return ""
	}
	// Récupérer l'utilisateur via son ID présent dans la session
	user, err := userService.GetUserByUUID(session.UserId.String())
	if err != nil {
		return ""
	}
	return user
}
