package entities

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Message   []byte
	MessageId string
	CreatedAt time.Time
}

func NewMessageWithID(message []byte) *Message {
	return &Message{
		Message:   message,
		MessageId: fmt.Sprintf("%d-0", time.Now().UnixMilli()),
		CreatedAt: time.Now(),
	}
}

func NewMessage() *Message {
	return new(Message)
}

func (q *Message) SetMessage(data []byte) {
	q.Message = data
}
func (q *Message) GetMessage() []byte {
	return q.Message
}

func (q *Message) SetMessageId(id string) bool {
	if !q.checkValidFormatID(id) {
		return false
	}
	q.MessageId = id
	return true
}

func (q *Message) GetMessageId() string {
	return q.MessageId
}

func (q *Message) GetCreatedAt() time.Time {
	return q.CreatedAt
}
func (q *Message) SetCreateAt(date time.Time) {
	q.CreatedAt = date
}

func (q *Message) checkValidFormatID(id string) bool {
	tem := strings.Split(id, "-")
	if len(tem) != 2 {
		return false
	}
	_, err1 := strconv.ParseUint(tem[0], 10, 64)
	_, err2 := strconv.ParseUint(tem[1], 10, 64)
	return err1 == nil && err2 == nil
}
