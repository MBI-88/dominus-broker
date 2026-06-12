package inbound

import (
	streambidirectional "dominus-broker/internal/application/usecases/stream_bidirectional"
	streamclient "dominus-broker/internal/application/usecases/stream_client"
	streamserver "dominus-broker/internal/application/usecases/stream_server"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"
	"dominus-broker/internal/infrastructure/grpc/mappers"
	"fmt"
	"io"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type brokerStream struct {
	pb.UnimplementedBrokerAPIServer
	c   streamclient.StreamClientUseCase
	s   streamserver.StreamServerUseCase
	b   streambidirectional.StreamBidirectionalUseCase
	log event.Event
}

func NewBrokerStream(
	server *grpc.Server,
	c streamclient.StreamClientUseCase,
	s streamserver.StreamServerUseCase,
	b streambidirectional.StreamBidirectionalUseCase,
	log event.Event,
) {
	gsrv := &brokerStream{
		c:   c,
		s:   s,
		b:   b,
		log: log,
	}
	pb.RegisterBrokerAPIServer(server, gsrv)
}

// Receives array messages from client
func (s *brokerStream) ClientStream(stream pb.BrokerAPI_ClientStreamServer) error {
	s.log.WriteLog(stream.Context(), enum.DEBUG, "inbound.ClientStream", enum.REQUEST_OK)
	ctx := mappers.NewClientStreamContext(stream)
	err := s.c.StreamClient(ctx)

	if err == io.EOF {
		s.log.WriteLog(stream.Context(), enum.DEBUG, "inbound.ClientStream.EOF", err.Error())
		return stream.SendAndClose(&pb.StreamResponseMessage{
			Status: int64(codes.OK),
		})
	}

	if err != nil {
		s.log.WriteLog(stream.Context(), enum.ERROR, "inbound.ClientStream.StreamClient", err.Error())
		return status.Error(codes.Aborted, err.Error())
	}

	return stream.SendAndClose(&pb.StreamResponseMessage{
		Status: int64(codes.OK),
	})
}

// Sends array messages to client
func (s *brokerStream) ServerStream(ms *pb.StreamRequestMessage, stream pb.BrokerAPI_ServerStreamServer) error {
	s.log.WriteLog(stream.Context(), enum.DEBUG, "inbound.ServerStream", enum.REQUEST_OK)
	ctx := mappers.NewServerStreamContext(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.s.StreamServer(ms, ctx)))
}

// Receives and sends messages from server to client
func (s *brokerStream) BidirectionalStream(stream pb.BrokerAPI_BidirectionalStreamServer) error {
	s.log.WriteLog(stream.Context(), enum.DEBUG, "inbound.BidirectionalStream", enum.REQUEST_OK)
	ctx := mappers.NewBiStreamConn(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.b.StreamBidirectional(ctx)))
}
