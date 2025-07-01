package grpcconn

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

type grpcService struct {
	client repos.GrpClientInt
	lg     repos.LogsInt
	topics entities.TopicsInt
}



type GrpcServiceInt interface {
	SimpleConn(ms repos.GrpRequestMessageInt) error
	StreamClientConn(st repos.StreamClientInt) error
	StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error
	StreamBiConn(st repos.StreamBiInt) error
	RunQueue(close <-chan struct{}) error
}

func NewGrpcService(lclient repos.LogsInt, gclient repos.GrpClientInt, topics entities.TopicsInt) GrpcServiceInt {
	return &grpcService{
		lg:     lclient,
		client: gclient,
		topics: topics,
	}
}
