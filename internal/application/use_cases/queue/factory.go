package queue

import (
	"context"
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repositories"
)

type Queue interface {
	Producer(ctx context.Context, ms dtos.ProducerDto) error 
	Consumer(ctx context.Context) (*entities.Queue, error)
	ConsumerDLT(ctx context.Context, ms dtos.ConsumerDto) error
	CheckMemory()
}

type queue struct {
	client   repositories.MemoryClient
	lg       repositories.Logs
	memory  entities.Memory

}

func NewQueue(lg repositories.Logs, client repositories.MemoryClient, m entities.Memory) Queue {
	return &queue{
		lg:     lg,
		client: client,
		memory: m,
	}
}
