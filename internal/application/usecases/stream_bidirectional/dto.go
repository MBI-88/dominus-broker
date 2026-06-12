package streambidirectional

import (
	"context"
	"dominus-broker/internal/application/dto"
)

type BrokerBidirectionalDto interface {
	Recv() (dto.BrokerRequestDto, error)
	Send(msg []byte) error
	Context() context.Context
}
