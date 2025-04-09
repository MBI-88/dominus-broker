package entities

import "fmt"

type topics struct {
	ar []*Topic
}

func (ts *topics) Append(t *Topic) {
	mutex.RLock()
	defer mutex.RUnlock()
	ts.ar = append(ts.ar, t)
}

func (ts *topics) Update(t *Topic) error {
	if _, err := ts.findTopic(t.Name); err != nil {
		return err
	}
	return nil
}

func (ts *topics) Delete(t *Topic) error {
	if _, err := ts.findTopic(t.Name); err != nil {
		return err
	}
	return nil
}

func (ts *topics) Find(topic string) (*Topic, error) {
	return ts.findTopic(topic)
}

func (ts *topics) findTopic(name string) (*Topic, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	for i,j := 0, len(ts.ar)-1; i < j; i,j = i+1, j-1 {
		if ts.ar[i].Name == name {
			return ts.ar[i], nil 
		}else if ts.ar[j].Name == name {
			return ts.ar[j], nil
		}
	}
	return nil, fmt.Errorf("Nod found")
}

type topicsInt interface {
	Append(t *Topic)
	Update(t *Topic) error
	Delete(t *Topic) error
	Find(topic string) (*Topic, error)
}

func NewTopics() topicsInt {
	return &topics{
		ar: make([]*Topic, 0, 100),
	}
}