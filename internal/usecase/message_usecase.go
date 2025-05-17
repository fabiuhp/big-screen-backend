package usecase

import (
	"errors"
	"time"

	"github.com/abiopereira/sw-criciuma/internal/domain"
	"github.com/google/uuid"
)

type messageUseCase struct {
	messageRepo domain.MessageRepository
}

func NewMessageUseCase(repo domain.MessageRepository) domain.MessageUseCase {
	return &messageUseCase{
		messageRepo: repo,
	}
}

func (uc *messageUseCase) Create(message *domain.Message) error {
	if message.Duration > 30 {
		return errors.New("duration cannot be greater than 30 seconds")
	}

	if message.Type != "text" && message.Type != "image" && message.Type != "video" {
		return errors.New("invalid message type")
	}

	message.ID = uuid.New().String()
	message.CreatedAt = time.Now()

	return uc.messageRepo.Create(message)
}

func (uc *messageUseCase) GetAll() ([]domain.Message, error) {
	return uc.messageRepo.GetAll()
}

func (uc *messageUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.messageRepo.Delete(id)
}
