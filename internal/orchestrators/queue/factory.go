package queue

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

type Queue interface {
	Provider(ctx context.Context, ms adapters.ProviderDto) error 
	Consumer(ctx context.Context) (*entities.Queue, error)
	ConsumerDLT(ctx context.Context, ms adapters.ConsumerDto) error
}

type queue struct {
	client   adapters.MemoryClient
	lg       adapters.Logs
	memory  entities.Memory

}

func NewQueue(lg adapters.Logs, client adapters.MemoryClient, m entities.Memory) Queue {
	return &queue{
		lg:     lg,
		client: client,
		memory: m,
	}
}
