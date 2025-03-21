package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"

	"github.com/gofrs/uuid"
)

type MessageRepository struct {
	DB *sql.DB
}

func (mr *MessageRepository) InsertMessage(message *models.Message) error {
	query := `INSERT INTO messages (uuid, sender_id, receiver_id, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := mr.DB.Exec(query, message.MessageId.String(), message.SenderId, message.ReceiverId, message.Content, message.SentAt)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}
	return nil
}

func (mr *MessageRepository) GetMessagesByUserID(userID uuid.UUID) ([]*models.Message, error) {
	query := `SELECT uuid, sender_id, receiver_id, content, created_at FROM messages WHERE sender_id = ? OR receiver_id = ? ORDER BY created_at ASC`
	rows, err := mr.DB.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %v", err)
	}
	defer rows.Close()
	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.MessageId, &message.SenderId, &message.ReceiverId, &message.Content, &message.SentAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (mr *MessageRepository) GetMessagesBetweenUsers(senderID, receiverID string) ([]models.MessageResponse, error) {
	var messages []models.MessageResponse
	fmt.Println(senderID)
	fmt.Println(receiverID)
	query := `
        SELECT uuid, sender_id, receiver_id, content, created_at
       	FROM messages
        WHERE (sender_id = ? AND receiver_id = ?)
        OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at ASC
    `
	rows, err := mr.DB.Query(query, senderID, receiverID, receiverID, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.MessageResponse
		err := rows.Scan(
			&msg.MessageID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.SentAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	fmt.Println(messages)
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
