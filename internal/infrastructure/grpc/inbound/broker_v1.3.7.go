package inbound

import (
	"dominus-project/internal/infrastructure/event"
	"dominus-project/internal/infrastructure/grpc/mappers"
	"dominus-project/internal/application/use_cases/broker"
	"dominus-project/internal/infrastructure/enum"
	"fmt"
	"io"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type brokerAPI struct {
	pb.UnimplementedBrokerAPIServer
	br  broker.Broker
	log event.Event
}

func NewBrokerAPI(server *grpc.Server, br broker.Broker, log event.Event)  {
	gsrv := &brokerAPI{br: br, log: log}
	pb.RegisterBrokerAPIServer(server, gsrv)
}



// Receives array messages from client
func (s *brokerAPI) ClientStream(stream pb.BrokerAPI_ClientStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "ClientStream", enum.REQUEST_OK)
	ctx := mappers.NewClientStreamContext(stream)
	err := s.br.StreamClientConn(ctx)
	if err != io.EOF {
		go s.log.WriteLog(stream.Context(), enum.ERROR, "ClientStream", err.Error())
		return status.Error(codes.Aborted, err.Error())
	}
	return stream.SendAndClose(&pb.StreamResponseMessage{
		Status: int64(codes.OK),
	})
}

// Sends array messages to client
func (s *brokerAPI) ServerStream(ms *pb.StreamRequestMessage, stream pb.BrokerAPI_ServerStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "ServerStream", enum.REQUEST_OK)
	ctx := mappers.NewServerStreamContext(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.br.StreamServerConn(ms, ctx)))
}

// Receives and sends messages from server to client
func (s *brokerAPI) BidirectionalStream(stream pb.BrokerAPI_BidirectionalStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "BidirectionalStream", enum.REQUEST_OK)
	ctx := mappers.NewBiStreamConn(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.br.StreamBiConn(ctx)))
}
