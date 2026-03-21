package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repositories"
	"dominus-project/internal/domain/services"
	"os"
)

type SQS interface {
	Producer(ctx context.Context, ms dtos.ProducerDto) error 
	Consumer(ctx context.Context) (*entities.Message, error)
	ConsumerDLT(ctx context.Context, ms dtos.ConsumerDto) error
	CheckMemory()
	ReactivateMessage(ch <-chan os.Signal)
}

type sqs struct {
	client   repositories.MemoryClient
	memory  services.Memory

}

func NewSQS(client repositories.MemoryClient, m services.Memory) SQS {
	return &sqs{
		client: client,
		memory: m,
	}
}
