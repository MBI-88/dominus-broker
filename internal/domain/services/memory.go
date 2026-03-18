package services

import (
	"slices"
	"sync"
)

// Referency to ids
type Memory interface {
	Set(id string) 
	Delete(id string) string
	Find(id string) string
	CheckMemory() bool
}

type memory struct {
	idList []string
	lock   *sync.Mutex
}

func NewMemory() Memory {
	return &memory{
		idList: make([]string, 0, 100),
		lock:   new(sync.Mutex),
	}
}

func (m *memory) Set(id string)  {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.idList = append(m.idList, id)
}

func (m *memory) Delete(id string) string {
	m.lock.Lock()
	defer m.lock.Unlock()

	pos, ok := slices.BinarySearch(m.idList, id)
	if ok {
		m.idList = slices.Delete(m.idList, pos, pos+1)
		return id
	}
	return ""
}

func (m *memory) Find(id string) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	pos, ok := slices.BinarySearch(m.idList, id)
	if ok {
		return m.idList[pos]
	}
	return ""
}

func (m *memory) CheckMemory() bool {
	m.lock.Lock() 
	defer m.lock.Unlock()
	return len(m.idList) > 0 // true idList is not empty
}