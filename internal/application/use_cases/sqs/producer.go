package sqs

import (
	"context"
	"dominus-broker/internal/application/dtos"
	"dominus-broker/internal/domain/entities"
	"fmt"
)

func (s *sqs) Producer(ctx context.Context, ms dtos.ProducerDto) error {
	payload := ms.GetPayload()
	if len(payload) == 0 {
		return fmt.Errorf("empty payload")
	}
	q := entities.NewMessageWithID(payload)
	return s.client.SendMessage(ctx, q)
}
