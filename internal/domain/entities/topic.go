package entities

import (
	"sync"
)

type Topic struct {
	Name        string   `json:"name" validate:"omitempty,alpha,lowercase"`
	Subscribers []string `json:"subscribers" validate:"required,dive,hostname_port"`
	Message     []byte
	Lck         *sync.Mutex
}

func (t *Topic) SetMessage(data []byte)  {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	t.Message = data
}

func (t *Topic) GetMessage() []byte {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Message
}

func (t *Topic) GetSubscribers() []string {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Subscribers
}

func (t *Topic) SetSubscribers(sb []string) {
	t.Lck.Lock()
	defer t.Lck.Unlock() 
	t.Subscribers = sb
}

func (t *Topic) GetName() string {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Name
}

func (t *Topic) SetName(name string) {
	t.Lck.Lock()
	defer t.Lck.Unlock() 
	t.Name = name
}


type TopicInt interface {
	SetMessage(data []byte)
	GetMessage() []byte
	SetSubscribers(sb []string)
	GetSubscribers() []string
	GetName() string
	SetName(name string)
}