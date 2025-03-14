package utils

import (
	"encoding/json"
	"forum/models"
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid"
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

	// Récupérer la liste des utilisateurs connectés
	var onlineUsers []*models.UserList
	for _, user := range h.clients {
		onlineUsers = append(onlineUsers, user)
	}

	// Préparer le message JSON avec un champ "type"
	response := map[string]interface{}{
		"type":  "users", // Ajoutez un champ "type" pour identifier le message
		"users": onlineUsers,
	}

	message, err := json.Marshal(response)
	if err != nil {
		log.Println("Erreur lors de la sérialisation du message :", err)
		return
	}

	// Envoyer à chaque client
	for client := range h.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Erreur lors de l'envoi du message :", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
func (h *Hub) BroadcastMessage(receiverID string, message *models.Message) {
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
func SendMessage(senderID, receiverID string, content string) (*models.Message, error) {
	message := &models.Message{
		MessageID:  uuid.Must(uuid.NewV4()), // Générer un nouvel UUID pour le message
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		SentAt:     time.Now(),
		IsRead:     false,
	}

	// Ici, vous pouvez enregistrer le message dans la base de données
	// Par exemple : db.Create(&message).Error

	return message, nil
}

func GetMessages(senderID, receiverID string) ([]*models.Message, error) {
	var messages []*models.Message

	// Ici, vous pouvez récupérer les messages depuis la base de données
	// Par exemple : db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID).Find(&messages).Error

	return messages, nil
}
