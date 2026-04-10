package broker

import "context"

type BrokerRequestDto interface {
	GetPayload() []byte
	GetSubscribers() []string
}

type BrokerBidirectionalDto interface {
	Recv() (BrokerRequestDto, error)
	Send(msg []byte) error
	Context() context.Context
}

type BrokerClientDto interface {
	Recv() (BrokerRequestDto, error)
	Context() context.Context
}

type BrokerServerDto interface {
	Send(payload []byte) error
	Context() context.Context
}
