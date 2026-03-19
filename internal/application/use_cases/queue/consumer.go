package queue

import (
	"context"
	"dominus-project/internal/domain/entities"
	"fmt"
)

func (c *queue) Consumer(ctx context.Context) (*entities.Queue, error) {
	for i := 0; i < c.memory.Len(); i++ {
		id := c.memory.Iter(i)
		payload, err := c.client.GetMessage(ctx, id)
		if err != nil {
			return nil, err
		}

		if !payload.Hidden {
			return payload, nil
		}
	}
	return nil, fmt.Errorf("not messages found")
}
