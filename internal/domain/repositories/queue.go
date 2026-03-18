package repositories

import (
	"context"
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/services"
)


type MemoryClient interface {
	SendMessage(ctx context.Context, q *entities.Queue) error
	DeleteMessage(ctx context.Context, ID string) error
	GetMessage(ctx context.Context, key string) (*entities.Queue, error)
	GetKeys(ctx context.Context, mem services.Memory) error
}

type CheckerClient interface {
	SaveConsumer(ctx context.Context, key string ) error
	CheckConsumer(ctx context.Context, key string) bool
}