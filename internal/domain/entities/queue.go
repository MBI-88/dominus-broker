package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Message   []byte    `json:"message"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Hidden    bool      `json:"hidden"`
}

func NewMessage(message []byte) *Message {
	return &Message{
		Message:   message,
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Hidden:    false,
	}
}

func (q *Message) SetMessage(data []byte) {
	q.Message = data
}
func (q *Message) GetMessage() []byte {
	return q.Message
}

func (q *Message) SetHidden(ok bool) {
	q.Hidden = ok
}
func (q *Message) GetHidden() bool {
	return q.Hidden
}

func (q *Message) SetID(id uuid.UUID) {
	q.ID = id
}
func (q *Message) GetID() string {
	return q.ID.String()
}

func (q *Message) GetCreatedAt() time.Time {
	return q.CreatedAt
}
func (q *Message) SetCreateAt(date time.Time) {
	q.CreatedAt = date
}
