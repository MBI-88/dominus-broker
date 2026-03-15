package repositories

import (
	"context"
	"dominus-project/internal/domain/entities"
)


type MemoryClient interface {
	SendMessage(ctx context.Context, q *entities.Queue) error
	DeleteMessage(ctx context.Context, ID string) error
	GetMessage(ctx context.Context, key string) (*entities.Queue, error)
	GetKeys(ctx context.Context, mem entities.Memory) error
}

type CheckerClient interface {
	SaveConsumer(ctx context.Context, key string ) error
	CheckConsumer(ctx context.Context, key string) bool
}