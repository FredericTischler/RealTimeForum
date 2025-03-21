package services

import (
	"forum/models"
	"forum/repositories"
)

type MessageService struct {
	MessageRepo *repositories.MessageRepository
}

func (ms *MessageService) GetMessage(userId, withUserId string, offset int) ([]*models.Message, error) {
	return ms.MessageRepo.LoadMessages(userId, withUserId, offset)
}
