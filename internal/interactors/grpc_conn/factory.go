package grpcconn

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

type GrpcService interface {
	SimpleConn(ms adapters.GrpcDto) error
	StreamClientConn(st adapters.StreamClient) error
	StreamServerConn(req adapters.GrpcDto, st adapters.StreamServer) error
	StreamBiConn(st adapters.StreamBi) error
	RunQueue(close <-chan struct{}) error
}

type grpcService struct {
	client adapters.GrpcClient
	lg     adapters.Logs
	topics entities.Topics
}

func NewGrpcService(lclient adapters.Logs, gclient adapters.GrpcClient, topics entities.Topics) GrpcService {
	return &grpcService{
		lg:     lclient,
		client: gclient,
		topics: topics,
	}
}
