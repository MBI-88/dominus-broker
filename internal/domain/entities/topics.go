package entities

import (
	"fmt"
	"sync"
)

type Topics interface {
	Append(t Topic) error
	Update(t Topic) error
	Delete(name string) error
	Find(topic string) (Topic, error)
	Next() (Topic, error)
	GetTopicsInfo() map[string]any
}

type topics struct {
	ar    []Topic
	limit int
	next  int
	mutex *sync.RWMutex
}

func NewTopics(limit int) Topics {
	return &topics{
		ar:    make([]Topic, 0, limit),
		limit: limit,
		mutex: new(sync.RWMutex),
	}
}

func (ts *topics) Append(t Topic) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	if ts.checkLenght() {
		ts.ar = append(ts.ar, t)
		return nil
	}
	return fmt.Errorf("Limit reached")
}

func (ts *topics) Update(t Topic) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	p, err := ts.findTopic(t.GetName())
	if err != nil {
		return err
	}
	ts.ar[p].SetSubscribers(t.GetSubscribers())
	return nil
}

func (ts *topics) Delete(name string) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	p, err := ts.findTopic(name)
	if err != nil {
		return err
	}
	left := ts.ar[:p]
	right := ts.ar[p+1:]
	ts.ar = make([]Topic, 0, ts.limit)
	ts.ar = append(ts.ar, left...)
	ts.ar = append(ts.ar, right...)
	return nil
}

func (ts *topics) Find(topic string) (Topic, error) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	p, err := ts.findTopic(topic)
	if err != nil {
		return nil, err
	}
	return ts.ar[p], nil
}

func (ts *topics) findTopic(name string) (int, error) {
	for i, j := 0, len(ts.ar)-1; i <= j; i, j = i+1, j-1 {
		if ts.ar[i].GetName() == name {
			return i, nil
		} else if ts.ar[j].GetName() == name {
			return j, nil
		}
	}
	return 0, fmt.Errorf("nod found")
}

func (ts *topics) checkLenght() bool {
	if len(ts.ar) >= ts.limit {
		return false
	}
	return true
}

func (ts *topics) Next() (Topic, error) {
	if ts.next > len(ts.ar)-1 {
		ts.next = 0
		return nil, fmt.Errorf("End")
	}
	ts.mutex.RLock()
	topic := ts.ar[ts.next]
	ts.mutex.RUnlock()
	ts.next++
	return topic, nil
}

func (ts *topics) GetTopicsInfo() map[string]any {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	var (
		arrayTopic = make([]map[string]any, 0, len(ts.ar))
		result     = make(map[string]any, 3)
	)

	for _, tp := range ts.ar {
		topic := map[string]any{
			"name":        tp.GetName(),
			"subscribers": tp.GetSubscribers(),
			"message":     string(tp.GetMessage()),
		}
		arrayTopic = append(arrayTopic, topic)
	}
	result["topics"] = arrayTopic
	result["total"] = len(ts.ar)
	result["limit"] = ts.limit
	return result
}
