package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
)

func (s *sqs) Ack(ctx context.Context, ms dtos.ConsumerDto) error {
	return s.client.AckMessage(ctx, ms.GetMessageId(), ms.GetGroupId())
}