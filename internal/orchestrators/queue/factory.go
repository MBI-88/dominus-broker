package queue

import (
	"context"
	"dominus-project/internal/domain/adapters"
)

type Queue interface {
	Provider(ctx context.Context, ms adapters.ProviderDto) error 
	Consumer(ctx context.Context, ms adapters.ConsumerDto) error
}

type queue struct {
	provider adapters.ProviderClient
	consumer adapters.ConsumerClient
	lg       adapters.Logs
}

func NewQueue(lg adapters.Logs, prov adapters.ProviderClient, consu adapters.ConsumerClient) Queue {
	return &queue{
		lg:     lg,
		provider: prov,
		consumer: consu,
	}
}
