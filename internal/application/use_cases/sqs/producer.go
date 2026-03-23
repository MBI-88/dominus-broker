package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/entities"
	"fmt"
)

func (c *sqs) Producer(ctx context.Context, ms dtos.ProducerDto) error {
	payload := ms.GetPayload()
	if len(payload) == 0 {
		return  fmt.Errorf("empty payload")
	}
	q := entities.NewMessage(payload)
	return c.client.SendMessage(ctx, q)
}
