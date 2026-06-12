package streamserver

import "context"

type BrokerRequestDto interface {
	GetPayload() []byte
	GetSubscribers() []string
}

type BrokerServerDto interface {
	Send(payload []byte) error
	Context() context.Context
}
