package sqs

import (
	"context"

	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/domain/repositories"
)

type SQS interface {
	Producer(ctx context.Context, ms ProducerDto) error
	Consumer(ctx context.Context, ms ConsumerDto) (*entities.Message, error)
	Ack(ctx context.Context, ms ConsumerDto) error
}

type sqs struct {
	client repositories.MemoryClient
}

func NewSQS(client repositories.MemoryClient) SQS {
	return &sqs{
		client: client,
	}
}
