package interactors

import (
	"dominus-project/app/domain/repos"
)

type grpcService struct {
	client     repos.GrpClientInt
	lg         repos.LogsInt
}



type GrpcServiceInt interface {
	SimpleConn(ms repos.GrpRequestMessageInt) error
	StreamClientConn(st repos.StreamClientInt) error
	StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error
	StreamBiConn(st repos.StreamBiInt) error
}
