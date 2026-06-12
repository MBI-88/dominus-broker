package mappers

import (
	"context"
	"dominus-broker/internal/application/dto"
	streambidirectional "dominus-broker/internal/application/usecases/stream_bidirectional"
	streamclient "dominus-broker/internal/application/usecases/stream_client"
	streamserver "dominus-broker/internal/application/usecases/stream_server"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
)

type grpcBiContextStream struct {
	sr pb.BrokerAPI_BidirectionalStreamServer
}

func NewBiStreamConn(sr pb.BrokerAPI_BidirectionalStreamServer) streambidirectional.BrokerBidirectionalDto {
	return &grpcBiContextStream{
		sr: sr,
	}
}

func (g *grpcBiContextStream) Recv() (dto.BrokerRequestDto, error) {
	return g.sr.Recv()
}

func (g *grpcBiContextStream) Send(msg []byte) error {
	return g.sr.Send(&pb.StreamResponseMessage{
		Payload: msg,
		Status:  0,
	})
}

func (g *grpcBiContextStream) Context() context.Context {
	return g.sr.Context()
}

type grpcContextClientStream struct {
	sr pb.BrokerAPI_ClientStreamServer
}

func NewClientStreamContext(sr pb.BrokerAPI_ClientStreamServer) streamclient.BrokerClientDto {
	return &grpcContextClientStream{
		sr: sr,
	}
}

func (g *grpcContextClientStream) Recv() (dto.BrokerRequestDto, error) {
	return g.sr.Recv()
}

func (g *grpcContextClientStream) Context() context.Context {
	return g.sr.Context()
}

type grpcContextServerStream struct {
	sr pb.BrokerAPI_ServerStreamServer
}

func NewServerStreamContext(sr pb.BrokerAPI_ServerStreamServer) streamserver.BrokerServerDto {
	return &grpcContextServerStream{
		sr: sr,
	}
}

func (g *grpcContextServerStream) Send(payload []byte) error {
	return g.sr.Send(&pb.StreamResponseMessage{
		Payload: payload,
		Status:  0,
	})
}

func (g *grpcContextServerStream) Context() context.Context {
	return g.sr.Context()
}
