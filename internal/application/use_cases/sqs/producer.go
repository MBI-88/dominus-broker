package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/entities"
	"fmt"
)

func (s *sqs) Producer(ctx context.Context, ms dtos.ProducerDto) error {
	payload := ms.GetPayload()
	if len(payload) == 0 {
		return  fmt.Errorf("empty payload")
	}
	q := entities.NewMessageWithID(payload)
	return s.client.SendMessage(ctx, q)
}
