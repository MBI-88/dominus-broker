package dto

import (
	"context"
	"dominus-project/internal/domain/adapters"
	pb "github.com/PR0C0D3-MBI/dominus-proto-definition/src/dominus"
)

type grpcBiContextStream struct {
	sr pb.Grpc_BidirectionalStreamServer
}

func NewBiStreamConn(sr pb.Grpc_BidirectionalStreamServer) adapters.StreamBi {
	return &grpcBiContextStream{
		sr: sr,
	}
}

func (g *grpcBiContextStream) Recv() (adapters.GrpcDto, error) {
	return g.sr.Recv()
}

func (g *grpcBiContextStream) Send(msg []byte) error {
	return g.sr.Send(&pb.ResponseMessage{
		Payload: msg,
	})
}

type grpcContextClientStream struct {
	sr pb.Grpc_ClientStreamServer
}

func NewClientStreamContext(sr pb.Grpc_ClientStreamServer) adapters.StreamClient {
	return &grpcContextClientStream{
		sr: sr,
	}
}

func (g *grpcContextClientStream) Recv() (adapters.GrpcDto, error) {
	return g.sr.Recv()
}

type grpcContextServerStream struct {
	sr pb.Grpc_ServerStreamServer
}

func NewServerStreamContext(sr pb.Grpc_ServerStreamServer) adapters.StreamServer {
	return &grpcContextServerStream{
		sr: sr,
	}
}

func (g *grpcContextServerStream) Send(payload []byte) error {
	return g.sr.Send(&pb.ResponseMessage{
		Payload: payload,
	})
}

func (g *grpcContextServerStream) Context() context.Context {
	return g.sr.Context()
}
