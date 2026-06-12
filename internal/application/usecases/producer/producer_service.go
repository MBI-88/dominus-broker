package producer

import (
	"context"
	"dominus-broker/internal/domain/entities"
	"fmt"
)

func (p *producerUseCase) Producer(ctx context.Context, ms ProducerDto) error {
	payload := ms.GetPayload()
	if len(payload) == 0 {
		return fmt.Errorf("producer.Producer empty payload")
	}
	q := entities.NewMessageWithID(payload)
	return p.client.SendMessage(ctx, q)
}
