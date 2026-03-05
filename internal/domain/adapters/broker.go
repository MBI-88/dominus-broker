package adapters

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

// Broker outputs
type BrokerClient interface {
	ClientStream(urls []string, msg <-chan []byte, ctx context.Context)
	ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, done chan<- struct{})
	BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{})
}

type BrokerServerDto interface {
	Send(payload []byte) error
	Context() context.Context
}
