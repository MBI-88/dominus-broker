package broker

import (
	"dominus-broker/internal/application/dtos"
	"dominus-broker/internal/domain/repositories"
)

type Broker interface {
	StreamClientConn(st dtos.BrokerClientDto) error
	StreamServerConn(req dtos.BrokerRequestDto, st dtos.BrokerServerDto) error
	StreamBiConn(st dtos.BrokerBidirectionalDto) error
}

type broker struct {
	client repositories.BrokerClient
}

func NewBroker(gclient repositories.BrokerClient) Broker {
	return &broker{
		client: gclient,
	}
}
