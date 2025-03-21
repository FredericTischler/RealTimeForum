package services

import "forum/repositories"

type MessageService struct {
	MessageRepo *repositories.MessageRepository
}
