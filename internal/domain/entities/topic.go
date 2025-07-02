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

func (t *Topic) GetSubscribers() []string {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Subscribers
}

func (t *Topic) GetName() string {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Name
}

func (t *Topic) GetMessage() []byte {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Message
}