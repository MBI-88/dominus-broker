package streamclient

import (
	"context"
	"dominus-broker/internal/application/dto"
)

type BrokerClientDto interface {
	Recv() (dto.BrokerRequestDto, error)
	Context() context.Context
}
