package producer

import (
	"context"
	"dominus-broker/internal/domain/repositories"
)

type ProducerUseCase interface {
	Producer(ctx context.Context, ms ProducerDto) error
}

type producerUseCase struct {
	client repositories.MemoryClient
}

func New(client repositories.MemoryClient) ProducerUseCase {
	return &producerUseCase{client: client}
}
