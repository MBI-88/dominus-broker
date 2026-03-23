package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Message   []byte    `json:"message"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessage(message []byte) *Message {
	return &Message{
		Message:   message,
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
	}
}

func (q *Message) SetMessage(data []byte) {
	q.Message = data
}
func (q *Message) GetMessage() []byte {
	return q.Message
}

func (q *Message) SetID(id string) {
	q.ID = id
}
func (q *Message) GetID() string {
	return q.ID
}

func (q *Message) GetCreatedAt() time.Time {
	return q.CreatedAt
}
func (q *Message) SetCreateAt(date time.Time) {
	q.CreatedAt = date
}
