package queue

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
	"sync"

	"github.com/google/uuid"
)

type Queue interface {
	Provider(ctx context.Context, ms adapters.ProviderDto) error 
	Consumer(ctx context.Context) (*entities.Queue, error)
	ConsumerDLT(ctx context.Context, ms adapters.ConsumerDto) error
}

type queue struct {
	memory   adapters.MemoryClient
	lg       adapters.Logs
	idList  []uuid.UUID
	lock    *sync.Mutex
}

func NewQueue(lg adapters.Logs, memory adapters.MemoryClient) Queue {
	return &queue{
		lg:     lg,
		memory: memory,
		idList: make([]uuid.UUID, 100),
		lock: new(sync.Mutex),
	}
}
