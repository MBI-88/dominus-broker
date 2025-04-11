package entities

import (
	"fmt"
	"sync"
)

var (
	mutex = new(sync.Mutex)
)

type node struct {
	data []byte
	next *node
	limit int64
	total int64
}

func (n *node) Push(data []byte) error {
	mutex.Lock()
	if n.data == nil {
		if n.isFull() {
			return fmt.Errorf("Full queue")
		}
		n.data = data
		n.next = new(node)
		n.total++
		return nil 
	}
	mutex.Unlock()
	n.next.Push(data)
	return nil 
}

func (n *node) Pop() []byte {
	mutex.Lock()
	if n == nil {
		return nil
	}
	if n.data != nil {
		data := n.data
		n = n.next
		return data
	}
	mutex.Unlock()
	return n.next.Pop()
}

func (n *node) isFull() bool {
	if n.total >= n.limit {
		return true
	}
	return false
}

type QueueInt interface {
	Pop() []byte
	Push(data []byte) error
}

func NewQueue(limit int64) QueueInt {
	return &node{
		limit: limit,
	}
}