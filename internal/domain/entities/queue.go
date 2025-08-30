package entities

type Queue interface {
	Enqueue(ms []byte)
	Length() int
	Dequeue() ([]byte, bool)
}

type queue [][]byte

func NewQueue() Queue {
	return new(queue)
}

func (q *queue) Enqueue(ms []byte) {
	*q = append(*q, ms)
}

func (q *queue) Dequeue() ([]byte, bool) {
	if len(*q) == 0 {
		return nil, false
	}
	ms := (*q)[0]
	*q = (*q)[1:]
	return ms, true
}

func (q *queue) Length() int {
	return len(*q)
}
