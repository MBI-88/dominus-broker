package broker

import (
	"dominus-broker/internal/domain/repositories"
)

type Broker interface {
	StreamClientConn(st BrokerClientDto) error
	StreamServerConn(req BrokerRequestDto, st BrokerServerDto) error
	StreamBiConn(st BrokerBidirectionalDto) error
}

type broker struct {
	client repositories.BrokerClient
}

func NewBroker(gclient repositories.BrokerClient) Broker {
	return &broker{
		client: gclient,
	}
}
