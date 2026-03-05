package dto

import (
	"context"
	"dominus-project/internal/domain/adapters"

	pb "github.com/PR0C0D3-MBI/dominus-proto-definition/dominus"
)

type grpcBiContextStream struct {
	sr pb.API_BidirectionalStreamServer
}

func NewBiStreamConn(sr pb.API_BidirectionalStreamServer) adapters.BrokerBidirectionalDto {
	return &grpcBiContextStream{
		sr: sr,
	}
}

func (g *grpcBiContextStream) Recv() (adapters.BrokerRequestDto, error) {
	return g.sr.Recv()
}

func (g *grpcBiContextStream) Send(msg []byte) error {
	return g.sr.Send(&pb.StreamResponseMessage{
		Payload: msg,
		Status: 0,
	})
}

func (g *grpcBiContextStream) Context() context.Context {
	return  g.sr.Context()
}





type grpcContextClientStream struct {
	sr pb.API_ClientStreamServer
}

func NewClientStreamContext(sr pb.API_ClientStreamServer) adapters.BrokerClientDto {
	return &grpcContextClientStream{
		sr: sr,
	}
}

func (g *grpcContextClientStream) Recv() (adapters.BrokerRequestDto, error) {
	return g.sr.Recv()
}

func (g *grpcContextClientStream) Context() context.Context {
	return g.sr.Context()
}




type grpcContextServerStream struct {
	sr pb.API_ServerStreamServer
}

func NewServerStreamContext(sr pb.API_ServerStreamServer) adapters.BrokerServerDto {
	return &grpcContextServerStream{
		sr: sr,
	}
}

func (g *grpcContextServerStream) Send(payload []byte) error {
	return g.sr.Send(&pb.StreamResponseMessage{
		Payload: payload,
		Status: 0,
	})
}

func (g *grpcContextServerStream) Context() context.Context {
	return g.sr.Context()
}
