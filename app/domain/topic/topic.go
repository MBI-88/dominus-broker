package topic

import (
	"fmt"
	"sync"
)


type topic struct {
	tp map[string][]string
	mu sync.RWMutex
}

 
func (t *topic) CreateTopic(key string, sub []string) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.tp[key] = sub
}

func (t *topic) GetSubcribers(key string) []string {
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
	if _,ok := t.tp[key]; !ok {
		return fmt.Errorf("Not found key")

	}
	delete(t.tp, key)
	return nil
}


type TopicInt interface {
	//Create a new topic in memory
	//
	//Parameters
	//
	//* key: topic name
	//
	//* sub: subscribers
	CreateTopic(key string, sub []string)
	//Return subcribers based on a key given
	//
	//Parameters
	//
	//* key: topic name
	GetSubcribers(key string) []string
	//Delete a topic using a selected key
	//
	//Parameters
	//
	//* key: topic name
	DeleteTopic(key string) error 
}


func NewTopic() TopicInt {
	return &topic{tp: make(map[string][]string)}
}