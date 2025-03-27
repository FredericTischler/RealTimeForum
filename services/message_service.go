package services

import (
	"forum/models"
	"forum/repositories"
)

type MessageService struct {
	MessageRepo *repositories.MessageRepository
}

func (ms *MessageService) GetMessage(userId, withUserId string, offset int) (*[]models.Message, error) {
	return ms.MessageRepo.LoadMessages(userId, withUserId, offset)
}

func (ms *MessageService) InsertMessage(message *models.Message) error {
	return ms.MessageRepo.InsertMessage(message)
}

func (ms *MessageService) GetUnreadMessagesCount(userId string) (int, error) {
	return ms.MessageRepo.GetUnreadMessagesCount(userId)
}

func (ms *MessageService) MarkMessagesAsRead(senderId, receiverId string) error {
	return ms.MessageRepo.MarkMessagesAsRead(senderId, receiverId)
}
