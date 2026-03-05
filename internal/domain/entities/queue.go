package entities

import (
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	Message   []byte    `json:"message"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Hidden    bool      `json:"hidden"`
}

func NewQueue(message []byte) *Queue {
	return &Queue{
		Message:   message,
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Hidden:    false,
	}
}
