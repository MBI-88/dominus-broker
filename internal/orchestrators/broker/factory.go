package broker

import (
	"dominus-project/internal/domain/adapters"
)

type Broker interface {
	StreamClientConn(st adapters.BrokerClientDto) error
	StreamServerConn(req adapters.BrokerRequestDto, st adapters.BrokerServerDto) error
	StreamBiConn(st adapters.BrokerBidirectionalDto) error
}

type broker struct {
	client adapters.BrokerClient
	lg     adapters.Logs
}

func NewBroker(lclient adapters.Logs, gclient adapters.BrokerClient) Broker {
	return &broker{
		lg:     lclient,
		client: gclient,
	}
}
