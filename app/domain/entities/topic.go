package entities

import (
	"fmt"
	"sync"
)

type Topic struct {
	Name       string   `json:"name" validate:"omitempty,alpha,lowercase"`
	Partitions []string `json:"partitions" validate:"required,dive,uri"`
	Queue      QueueInt
	Lck        *sync.Mutex
	Limit      int64
	total      int64
}

func (t *Topic) Push(data []byte) error {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	if t.isFull() {
		return fmt.Errorf("Queue full")
	}
	t.total++
	return t.Queue.push(data)
}

func (t *Topic) Pop() []byte {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return t.Queue.pop()
}

func (t *Topic) isFull() bool {
	if t.total >= t.Limit {
		return true
	}
	return false
}