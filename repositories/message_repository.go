package repositories

import (
	"database/sql"
	"forum/models"
)

type MessageRepository struct {
	DB *sql.DB
}

func (mr *MessageRepository) LoadMessages(userId, withUserId string, offset int) ([]*models.Message, error) {
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

	var messages []*models.Message
	for rows.Next() {
		var msg *models.Message
		err = rows.Scan(&msg.MessageId, &msg.SenderId, &msg.ReceiverId, &msg.Content, &msg.SentAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
