package entities

import (
	"fmt"
	"sync"
)

type Topic struct {
	Limit       int64
	total       int64
	Name        string   `json:"name" validate:"omitempty,alpha,lowercase"`
	Subscribers []string `json:"subscribers" validate:"required,dive,hostname"`
	Queue       QueueInt
	Lck         *sync.Mutex
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
