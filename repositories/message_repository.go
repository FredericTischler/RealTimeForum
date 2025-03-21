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

func (mr *MessageRepository) LoadMessages(userId, withUserId string, offset int) (*[]models.Message, error) {
	query := `SELECT uuid, sender_id, receiver_id, content, created_at 
			  FROM messages 
			  WHERE (sender_id = ? AND receiver_id = ?)
			  OR (sender_id = ? AND receiver_id = ?)
			  LIMIT 10
			  OFFSET ?
			  `
	rows, err := mr.DB.Query(query, userId, withUserId, withUserId, userId, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err = rows.Scan(&msg.MessageId, &msg.SenderId, &msg.ReceiverId, &msg.Content, &msg.SentAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return &messages, nil
}

func (mr *MessageRepository) InsertMessage(message *models.Message) error {
	newUuid, _ := uuid.NewV4()
	message.MessageId = newUuid.String()
	message.SentAt = time.Now()

	query := `INSERT INTO messages (uuid, sender_id, receiver_id, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := mr.DB.Exec(query, message.MessageId, message.SenderId, message.ReceiverId, message.Content, message.SentAt)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}
	return nil
}
