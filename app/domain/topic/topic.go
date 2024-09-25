package topic

import (
	"fmt"
	"sync"
)


type topic struct {
	tp map[string][]string
	mu sync.RWMutex
}


// CreateSubcribers adds new topic with subscribers 
func (t *topic) CreateTopic(key string, sub []string) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.tp[key] = sub
}

// GetSubcribers returns subcribers related to a topic
func (t *topic) GetSubcribers(key string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if subs, ok := t.tp[key]; ok {
		return subs
	}else {
		return []string{}
	}
}


func (t *topic) GetTopic() map[string][]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.tp
}

func (t *topic) DeleteTopic(key string) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if _,ok := t.tp[key]; !ok {
		return fmt.Errorf("Not found key")

	}
	delete(t.tp, key)
	return nil
}


type TopicInt interface {
	//Create a new topic in memory
	CreateTopic(key string, sub []string)
	//Return subcribers based on a key given
	GetSubcribers(key string) []string
	//Returns  every topic in memory
	GetTopic() map[string][]string
	//Delete a topic using a selected key
	DeleteTopic(key string) error 
}


func NewTopic() TopicInt {
	return &topic{tp: make(map[string][]string)}
}