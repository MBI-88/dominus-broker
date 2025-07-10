package grpcconn

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

type grpcService struct {
	client repos.IGrpClient
	lg     repos.ILogs
	topics entities.ITopics
}



type GrpcServiceInt interface {
	SimpleConn(ms repos.IGrpRequestMessage) error
	StreamClientConn(st repos.IStreamClient) error
	StreamServerConn(req repos.IGrpRequestMessage, st repos.IStreamServer) error
	StreamBiConn(st repos.IStreamBi) error
	RunQueue(close <-chan struct{}) error
}

func NewGrpcService(lclient repos.ILogs, gclient repos.IGrpClient, topics entities.ITopics) GrpcServiceInt {
	return &grpcService{
		lg:     lclient,
		client: gclient,
		topics: topics,
	}
}
