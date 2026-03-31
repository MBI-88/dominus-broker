package repositories

import (
	"context"
	"dominus-broker/internal/domain/entities"
)

type MemoryClient interface {
	SendMessage(ctx context.Context, q *entities.Message) error
	AckMessage(ctx context.Context, messageId, groupId string) error
	GetMessage(ctx context.Context, workerId, groupId string) (*entities.Message, error)
	Group(groupId string) error
}

type CheckerClient interface {
	SaveConsumer(ctx context.Context, key string) error
	CheckConsumer(ctx context.Context, key string) bool
}
