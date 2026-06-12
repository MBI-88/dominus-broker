package streamserver

import "dominus-broker/internal/domain/repositories"

type StreamServerUseCase interface {
	StreamServer(req BrokerRequestDto, st BrokerServerDto) error
}

type streamServerUseCase struct {
	client repositories.BrokerClient
}

func New(client repositories.BrokerClient) StreamServerUseCase {
	return &streamServerUseCase{client: client}
}
