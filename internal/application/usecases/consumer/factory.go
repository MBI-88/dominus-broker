package consumer

import (
	"context"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/domain/repositories"
)

type ConsumerUseCase interface {
	Consumer(ctx context.Context, ms ConsumerDto) (*entities.Message, error)
}

type consumerUseCase struct {
	client repositories.MemoryClient
}

func New(client repositories.MemoryClient) ConsumerUseCase {
	return &consumerUseCase{client: client}
}
