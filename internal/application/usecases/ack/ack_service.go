package ack

import (
	"context"
	"dominus-broker/internal/domain/entities"
	"fmt"
)

func (a *ackUseCase) Ack(ctx context.Context, ms AskDto) error {
	qs := entities.NewMessage()
	if !qs.SetMessageId(ms.GetMessageId()) {
		return fmt.Errorf("ack.Ack invalid messageId")
	}
	return a.client.AckMessage(ctx, ms.GetMessageId(), ms.GetGroupId())
}
