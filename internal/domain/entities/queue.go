package entities

import (
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	Message   []byte    `json:"message"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
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

func (q *Queue) SetMessage(data []byte) {
	q.Message = data
}
func (q *Queue) GetMessage() []byte {
	return q.Message
}

func (q *Queue) SetHidden(ok bool) {
	q.Hidden = ok
}
func (q *Queue) GetHidden() bool {
	return q.Hidden
}

func (q *Queue) SetID(id uuid.UUID) {
	q.ID = id
}
func (q *Queue) GetID() string {
	return q.ID.String()
}

func (q *Queue) GetCreatedAt() time.Time {
	return q.CreatedAt
}
func (q *Queue) SetCreateAt(date time.Time) {
	q.CreatedAt = date
}
