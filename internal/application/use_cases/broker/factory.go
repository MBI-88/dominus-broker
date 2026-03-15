package broker

import (
	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/repositories"
)

type Broker interface {
	StreamClientConn(st dtos.BrokerClientDto) error
	StreamServerConn(req dtos.BrokerRequestDto, st dtos.BrokerServerDto) error
	StreamBiConn(st dtos.BrokerBidirectionalDto) error
}

type broker struct {
	client repositories.BrokerClient
	lg     repositories.Logs
}

func NewBroker(lclient repositories.Logs, gclient repositories.BrokerClient) Broker {
	return &broker{
		lg:     lclient,
		client: gclient,
	}
}
