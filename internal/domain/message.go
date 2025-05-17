package domain

import "time"

type Message struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`
}

type MessageRepository interface {
	Create(message *Message) error
	GetAll() ([]Message, error)
}

type MessageUseCase interface {
	Create(message *Message) error
	GetAll() ([]Message, error)
}
