package sqs

import (
	"context"

	"dominus-broker/internal/domain/entities"
)

func (s *sqs) Consumer(ctx context.Context, ms ConsumerDto) (*entities.Message, error) {
	return s.client.GetMessage(ctx, ms.GetWorkerId(), ms.GetGroupId())
}
