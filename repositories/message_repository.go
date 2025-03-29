package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"

	"github.com/gofrs/uuid"
)

type MessageRepository struct {
	DB *sql.DB
}

func (mr *MessageRepository) LoadMessages(userId, withUserId string, offset int, onlyUnread bool) (*[]models.Message, error) {
	var query string
	var args []interface{}

	// Construction dynamique de la requête
	baseQuery := `SELECT uuid, sender_id, receiver_id, content, created_at, is_read
                 FROM messages 
                 WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))`

	args = append(args, userId, withUserId, withUserId, userId)

	if onlyUnread {
		baseQuery += " AND is_read = false AND receiver_id = ?"
		args = append(args, userId)
	}

	query = baseQuery + ` ORDER BY created_at DESC LIMIT 10 OFFSET ?`
	args = append(args, offset)

	rows, err := mr.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err = rows.Scan(&msg.MessageId, &msg.SenderId, &msg.ReceiverId, &msg.Content, &msg.SentAt, &msg.Is_read)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}
		messages = append(messages, msg)
	}

	return &messages, nil
}

func (mr *MessageRepository) InsertMessage(message *models.Message) error {
	newUuid, _ := uuid.NewV4()
	message.MessageId = newUuid.String()
	message.SentAt = time.Now()

	// Marquer comme lu si c'est un message que j'envoie (sender)
	isRead := false
	if message.SenderId == message.ReceiverId {
		isRead = true
	}

	query := `INSERT INTO messages (uuid, sender_id, receiver_id, content, created_at, is_read) 
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := mr.DB.Exec(query,
		message.MessageId,
		message.SenderId,
		message.ReceiverId,
		message.Content,
		message.SentAt,
		isRead) // Ajout du statut de lecture
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}
	return nil
}

func (mr *MessageRepository) GetUnreadMessagesCount(userId string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND is_read = 0`
	err := mr.DB.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %v", err)
	}
	return count, nil
}

func (mr *MessageRepository) MarkMessagesAsRead(senderId, receiverId string) error {
	query := `UPDATE messages SET is_read = true 
              WHERE sender_id = ? AND receiver_id = ? AND is_read = false`
	_, err := mr.DB.Exec(query, senderId, receiverId)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %v", err)
	}
	return nil
}
