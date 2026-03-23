package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/entities"
)

func (s *sqs) Consumer(ctx context.Context, ms dtos.ConsumerDto) (*entities.Message, error) {
	return s.client.GetMessage(ctx, ms.GetWorker(), ms.GetGroupId())
}
