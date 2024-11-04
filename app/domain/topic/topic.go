package topic

import (
	"fmt"
	"sync"
)


type topic struct {
	tp map[string][]string
	count uint8
	mu sync.RWMutex
}

 
func (t *topic) CreateTopic(key string, sub []string) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.count >= 10 {
		return fmt.Errorf("Limit reached")
	}
	t.count += 1
	t.tp[key] = sub
	return nil
}

func (t *topic) UpdateTopic(key string, sub []string) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if _, ok := t.tp[key]; !ok {
		return fmt.Errorf("Key not found")
	}
	t.tp[key] = sub
	return nil
}

func (t *topic) GetSusbcribers(key string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if subs, ok := t.tp[key]; ok {
		return subs
	}else {
		return []string{} 
	}
}

func (t *topic) DeleteTopic(key string) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if _, ok := t.tp[key]; !ok {
		return fmt.Errorf("Not found key")

	}
	delete(t.tp, key)
	t.count -= 1
	return nil
}


type TopicInt interface {
	//Create a new topic in memory
	//
	//Parameters
	//
	//-> key: topic name
	//
	//-> sub: subscribers
	CreateTopic(key string, sub []string) error
	//Return subcribers based on a key given
	//
	//Parameters
	//
	//-> key: topic name
	GetSusbcribers(key string) []string
	//Delete a topic using a selected key
	//
	//Parameters
	//
	//-> key: topic name
	DeleteTopic(key string) error 
	//Update a new topic in memory
	//
	//Parameters
	//
	//-> key: topic name
	//
	//-> sub: subscribers
	UpdateTopic(key string, sub []string) error
}


func NewTopic() TopicInt {
	return &topic{tp: make(map[string][]string)}
}