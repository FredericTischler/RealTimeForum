package utils

import (
	"encoding/json"
	"forum/models"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub central pour gérer les connexions WebSocket
type Hub struct {
	clients    map[*websocket.Conn]*models.UserList // Connexions WebSocket associées à des utilisateurs
	broadcast  chan models.Message                  // Canal pour diffuser les messages
	register   chan *websocket.Conn                 // Canal pour enregistrer de nouvelles connexions
	unregister chan *websocket.Conn                 // Canal pour déconnecter les utilisateurs
	mu         sync.Mutex                           // Mutex pour la synchronisation
}

// NewHub initialise un nouveau Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]*models.UserList),
		broadcast:  make(chan models.Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// AddClient ajoute une nouvelle connexion WebSocket au Hub
func (h *Hub) AddClient(conn *websocket.Conn, user *models.UserList) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = user
}

// RemoveClient supprime une connexion WebSocket du Hub
func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

// BroadcastUsers envoie la liste des utilisateurs connectés à tous les clients
func (h *Hub) BroadcastUsers() {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Récupérer la liste des identifiants
	var onlineUsers []*models.UserList
	for _, userID := range h.clients {
		onlineUsers = append(onlineUsers, userID)
	}

	// Préparer le message JSON
	message, err := json.Marshal(onlineUsers)
	if err != nil {
		return
	}

	// Envoyer à chaque client
	for client := range h.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			client.Close()
			delete(h.clients, client)
		}
	}
}

// BroadcastMessage envoie un message à un utilisateur spécifique
func (h *Hub) BroadcastMessage(receiverID string, message *models.Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn, user := range h.clients {
		if user.UserId == receiverID {
			err := conn.WriteJSON(map[string]interface{}{
				"type":    "message",
				"message": message,
			})
			if err != nil {
				conn.Close()
				delete(h.clients, conn)
			}
		}
	}
}

// Run démarre le Hub pour gérer les connexions et les messages
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = nil // Ajouter une nouvelle connexion sans utilisateur associé
			h.mu.Unlock()

		case conn := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, conn) // Supprimer une connexion
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				err := conn.WriteJSON(message)
				if err != nil {
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}
