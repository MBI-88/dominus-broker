package entities

import (
	"fmt"
	"sync"
)


type topics struct {
	ar    []*Topic
	limit int
	next int
	mutex *sync.Mutex
}

func (ts *topics) Append(t *Topic) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	if ts.checkLenght() {
		ts.ar = append(ts.ar, t)
		return nil
	}
	return fmt.Errorf("Limit reached")
}

func (ts *topics) Update(t *Topic) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	p, err := ts.findTopic(t.Name)
	if err != nil {
		return err
	}
	ts.ar[p].Subscribers = t.Subscribers
	return nil
}

func (ts *topics) Delete(t *Topic) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	p, err := ts.findTopic(t.Name)
	if err != nil {
		return err
	}
	left := ts.ar[:p]
	right := ts.ar[p+1:]
	ts.ar = make([]*Topic, 0, ts.limit)
	ts.ar = append(ts.ar, left...)
	ts.ar = append(ts.ar, right...)
	return nil
}

func (ts *topics) Find(topic string) (*Topic, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	p, err := ts.findTopic(topic)
	if err != nil {
		return nil, err
	}
	return ts.ar[p], nil
}

func (ts *topics) findTopic(name string) (int, error) {
	for i, j := 0, len(ts.ar)-1; i <= j; i, j = i+1, j-1 {
		if ts.ar[i].Name == name {
			return i, nil
		} else if ts.ar[j].Name == name {
			return j, nil
		}
	}
	return 0, fmt.Errorf("Nod found")
}

func (ts *topics) checkLenght() bool {
	if len(ts.ar) >= ts.limit {
		return false
	}
	return true
}

func (ts *topics) Next() (Topic, error) {
	if ts.next > len(ts.ar) - 1 {
		ts.next = 0
		return Topic{}, fmt.Errorf("End")
	}
	ts.mutex.Lock()
	topic := ts.ar[ts.next]
	ts.mutex.Unlock()
	ts.next++
	return *topic, nil 
}

type TopicsInt interface {
	Append(t *Topic) error
	Update(t *Topic) error
	Delete(t *Topic) error
	Find(topic string) (*Topic, error)
	Next() (Topic, error) 
}

func NewTopics(limit int) TopicsInt {
	return &topics{
		ar: make([]*Topic, 0, limit),
		limit: limit,
		mutex: new(sync.Mutex),
	}
}
