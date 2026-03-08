package entities

import (
	"slices"
	"sync"

	"github.com/google/uuid"
)

// Referency to ids
type Memory interface {
	Set(id uuid.UUID) string
	Delete(id uuid.UUID) string
	Find(id uuid.UUID) string
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

func (m *memory) Set(id uuid.UUID) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.idList = append(m.idList, id.String())
	return id.String()
}

func (m *memory) Delete(id uuid.UUID) string {
	m.lock.Lock()
	defer m.lock.Unlock()

	pos, ok := slices.BinarySearch(m.idList, id.String())
	if ok {
		m.idList = slices.Delete(m.idList, pos, pos+1)
		return id.String()
	}
	return ""
}

func (m *memory) Find(id uuid.UUID) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	pos, ok := slices.BinarySearch(m.idList, id.String())
	if ok {
		return m.idList[pos]
	}
	return ""
}
