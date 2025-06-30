package entities


type node struct {
	data []byte
	next *node
}

func (n *node) push(data []byte) error {
	if n.data == nil {
		n.data = data
		n.next = new(node)
		return nil 
	}
	return n.next.push(data)
}

func (n *node) pop() []byte {
	if n == nil {
		return nil
	}
	if n.data != nil {
		data := n.data
		n = n.next
		return data
	}
	return n.next.pop()
}


type QueueInt interface {
	pop() []byte
	push(data []byte) error
}

func NewQueue() QueueInt {
	return new(node)
}