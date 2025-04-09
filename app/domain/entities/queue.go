package entities 

import "sync"

var (
	mutex = new(sync.RWMutex)
)

type node struct {
	data []byte
	next *node
}

func (n *node) Add(data []byte) {
	mutex.RLocker()
	if n.data == nil {
		n.data = data
		n.next = new(node)
		return
	}
	mutex.RUnlock()
	n.next.Add(data)
}

func (n *node) Read() []byte {
	mutex.RLocker()
	if n == nil {
		return []byte{}
	}
	if n.data != nil {
		data := n.data
		n = n.next
		return data
	}
	mutex.RUnlock()
	return n.next.Read()
}

type QueueInt interface {
	Read() []byte
	Add([]byte)
}

func NewQueue() QueueInt {
	return new(node)
}