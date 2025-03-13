package utils

import (
	"encoding/json"
	"forum/models"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub central pour gérer les connexions
type Hub struct {
	clients map[*websocket.Conn]*models.UserList // Vous pouvez associer chaque connexion à un identifiant utilisateur
	mu      sync.Mutex
}

// NewHub initialise un hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]*models.UserList),
	}
}

// AddClient ajoute une nouvelle connexion et son identifiant utilisateur
func (h *Hub) AddClient(conn *websocket.Conn, user *models.UserList) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = user
}

// RemoveClient supprime une connexion
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
