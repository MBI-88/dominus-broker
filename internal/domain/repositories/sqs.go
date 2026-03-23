package repositories

import (
	"context"
	"dominus-project/internal/domain/entities"
)

type MemoryClient interface {
	SendMessage(ctx context.Context, q *entities.Message) error
	AckMessage(ctx context.Context, ID string) error
	GetMessage(ctx context.Context, key string) (*entities.Message, error)
	Group(ctx context.Context) error
}

type CheckerClient interface {
	SaveConsumer(ctx context.Context, key string) error
	CheckConsumer(ctx context.Context, key string) bool
}
