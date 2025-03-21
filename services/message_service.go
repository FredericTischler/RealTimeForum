package services

import (
	"forum/models"
	"forum/repositories"

	"github.com/gofrs/uuid"
)

type MessageService struct {
	MessageRepo *repositories.MessageRepository
}

func (ms *MessageService) InsertMessage(message *models.Message) error {
	return ms.MessageRepo.InsertMessage(message)
}

func (ms *MessageService) GetMessagesByUserID(userID uuid.UUID) ([]*models.Message, error) {
	return ms.MessageRepo.GetMessagesByUserID(userID)
}

func (ms *MessageService) GetMessagesBetweenUsers(senderID, receiverID string) ([]models.MessageResponse, error) {
	return ms.MessageRepo.GetMessagesBetweenUsers(senderID, receiverID)
}
