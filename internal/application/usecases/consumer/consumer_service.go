package consumer

import (
	"context"
	"dominus-broker/internal/domain/entities"
)

func (c *consumerUseCase) Consumer(ctx context.Context, ms ConsumerDto) (*entities.Message, error) {
	return c.client.GetMessage(ctx, ms.GetWorkerId(), ms.GetGroupId())
}
