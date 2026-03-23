package sqs

import (
	"context"
	"dominus-project/internal/application/dtos"
)

func (c *sqs) Ack(ctx context.Context, ms dtos.ConsumerDto) error {
	return c.client.AckMessage(ctx, ms.GetId())
}