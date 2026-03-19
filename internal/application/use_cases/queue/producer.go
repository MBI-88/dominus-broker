package queue

import (
	"context"
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/entities"
	"fmt"
)

func (c *queue) Producer(ctx context.Context, ms dtos.ProducerDto) error {
	payload := ms.GetPayload()
	if len(payload) == 0 {
		return  fmt.Errorf("empty payload")
	}
	q := entities.NewQueue(payload)
	c.memory.Set(q.GetID())
	return c.client.SendMessage(ctx, q)
}
