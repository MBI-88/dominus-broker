package entities

import "fmt"

type topics struct {
	ar []*Topic
}

func (ts *topics) Append(t *Topic) {
	mutex.Lock()
	defer mutex.Unlock()
	ts.ar = append(ts.ar, t)
}

func (ts *topics) Update(t *Topic) error {
	mutex.Lock()
	defer mutex.Unlock()
	p, err := ts.findTopic(t.Name)
	if err != nil {
		return err
	}
	ts.ar[p].Partitions = t.Partitions
	return nil
}

func (ts *topics) Delete(t *Topic) error {
	mutex.Lock()
	defer mutex.Lock()
	p, err := ts.findTopic(t.Name)
	if err != nil {
		return err
	}
	left := ts.ar[:p]
	right := ts.ar[p+1:]
	ts.ar = make([]*Topic, 0, 100)
	ts.ar = append(ts.ar, left...)
	ts.ar = append(ts.ar, right...) 
	return nil
}

func (ts *topics) Find(topic string) (*Topic, error) {
	mutex.Lock()
	defer mutex.Unlock()
	p, err := ts.findTopic(topic)
	if err != nil {
		return nil, err
	}
	return ts.ar[p], nil
}

func (ts *topics) findTopic(name string) (int, error) {
	for i,j := 0, len(ts.ar)-1; i < j; i,j = i+1, j-1 {
		if ts.ar[i].Name == name {
			return i, nil 
		}else if ts.ar[j].Name == name {
			return j, nil
		}
	}
	return 0, fmt.Errorf("Nod found")
}

type TopicsInt interface {
	Append(t *Topic)
	Update(t *Topic) error
	Delete(t *Topic) error
	Find(topic string) (*Topic, error)
}

func NewTopics() TopicsInt {
	return &topics{
		ar: make([]*Topic, 0, 100),
	}
}