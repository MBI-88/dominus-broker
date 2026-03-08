package adapters

import (
	"context"
	"dominus-project/internal/domain/entities"
)

type ProviderDto interface {
	GetPayload() []byte
}

type ConsumerDto interface {
	GetId() string
}


type MemoryClient interface {
	SendMessage(ctx context.Context, q *entities.Queue) error
	DeleteMessage(ctx context.Context, ID string) error
	GetMessage(ctx context.Context, key string) (*entities.Queue, error)
	GetKeys(ctx context.Context, mem entities.Memory) error
}
