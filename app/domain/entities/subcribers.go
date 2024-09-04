package entities 

import (
	"sync"
)


type subcribers struct {
	subs map[string][]string
	mu sync.RWMutex
}


// FeedSubcribers adds new topic with subscribers 
func (s *subcribers) FeedSubcribers(key string, sub []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.subs[key] = sub
}

// GetSubcribers returns subcribers related to a topic
func (s *subcribers) GetSubcribers(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if subs, ok := s.subs[key]; ok {
		return subs
	}else {
		return []string{}
	}
}



type subcribersInt interface {
	FeedSubcribers(key string, sub []string)
	GetSubcribers(key string) []string
}


func NewSubcribers() subcribersInt {
	return &subcribers{subs: make(map[string][]string)}
}