package grpcconn

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

type IGrpcService interface {
	SimpleConn(ms adapters.IGrpcDto) error
	StreamClientConn(st adapters.IStreamClient) error
	StreamServerConn(req adapters.IGrpcDto, st adapters.IStreamServer) error
	StreamBiConn(st adapters.IStreamBi) error
	RunQueue(close <-chan struct{}) error
}

type grpcService struct {
	client adapters.IGrpcClient
	lg     adapters.ILogs
	topics entities.ITopics
}

func NewGrpcService(lclient adapters.ILogs, gclient adapters.IGrpcClient, topics entities.ITopics) IGrpcService {
	return &grpcService{
		lg:     lclient,
		client: gclient,
		topics: topics,
	}
}
