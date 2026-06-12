package streamclient

import "dominus-broker/internal/domain/repositories"

type StreamClientUseCase interface {
	StreamClient(st BrokerClientDto) error
}

type streamClientUseCase struct {
	client repositories.BrokerClient
}

func New(client repositories.BrokerClient) StreamClientUseCase {
	return &streamClientUseCase{client: client}
}
