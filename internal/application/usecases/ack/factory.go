package ack

import (
	"context"
	"dominus-broker/internal/domain/repositories"
)

type AckUseCase interface {
	Ack(ctx context.Context, ms AskDto) error
}

type ackUseCase struct {
	client repositories.MemoryClient
}

func New(client repositories.MemoryClient) AckUseCase {

	return &ackUseCase{client: client}
}
