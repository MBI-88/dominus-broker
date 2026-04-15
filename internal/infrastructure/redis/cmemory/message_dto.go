package cmemory

import "time"

type MessageDto struct {
	Message   []byte    `json:"message"`
	MessageId string    `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}
