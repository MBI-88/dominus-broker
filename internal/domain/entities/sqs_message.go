package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Message   []byte    `json:"message"`
	MeesageId string    `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessage(message []byte) *Message {
	return &Message{
		Message:   message,
		MeesageId: uuid.New().String(),
		CreatedAt: time.Now(),
	}
}

func (q *Message) SetMessage(data []byte) {
	q.Message = data
}
func (q *Message) GetMessage() []byte {
	return q.Message
}

func (q *Message) SetMessageId(id string) {
	q.MeesageId = id
}
func (q *Message) GetMessageId() string {
	return q.MeesageId
}

func (q *Message) GetCreatedAt() time.Time {
	return q.CreatedAt
}
func (q *Message) SetCreateAt(date time.Time) {
	q.CreatedAt = date
}
