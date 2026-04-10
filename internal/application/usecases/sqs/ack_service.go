package sqs

import (
	"context"
	"dominus-broker/internal/domain/entities"
	"fmt"
)

func (s *sqs) Ack(ctx context.Context, ms ConsumerDto) error {
	qs := entities.NewMessage()
	if !qs.SetMessageId(ms.GetMessageId()) {
		return fmt.Errorf("sqs.Ack invalid messageId")
	}
	return s.client.AckMessage(ctx, ms.GetMessageId(), ms.GetGroupId())
}
