package streambidirectional

import "dominus-broker/internal/domain/repositories"

type StreamBidirectionalUseCase interface {
	StreamBidirectional(st BrokerBidirectionalDto) error
}

type streamBidirectionalUseCase struct {
	client repositories.BrokerClient
}

func New(client repositories.BrokerClient) StreamBidirectionalUseCase {
	return &streamBidirectionalUseCase{client: client}
}
