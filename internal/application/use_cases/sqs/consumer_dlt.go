package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
)

func (c *sqs) ConsumerDLT(ctx context.Context, ms dtos.ConsumerDto) error {
	id := ms.GetId()
	if err := c.client.DeleteMessage(ctx, id); err != nil {
		return err
	}
	c.memory.Delete(id) 
	return  nil
}