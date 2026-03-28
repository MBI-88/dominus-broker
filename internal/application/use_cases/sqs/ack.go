package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/entities"
	"fmt"
)

func (s *sqs) Ack(ctx context.Context, ms dtos.ConsumerDto) error {
	qs := entities.NewMessage()
	if !qs.SetMessageId(ms.GetMessageId()) {
		return fmt.Errorf("invalid messageId")
	}
	return s.client.AckMessage(ctx, ms.GetMessageId(), ms.GetGroupId())
}