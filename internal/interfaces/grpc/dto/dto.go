package dto

import (
	"context"
	"dominus-project/internal/domain/adapters"
	pb "dominus-project/internal/interfaces/grpc/proto/builder"
)

type grpcBiContextStream struct {
	sr pb.Grpc_BidirectionalStreamServer
}

func NewBiStreamConn(sr pb.Grpc_BidirectionalStreamServer) adapters.IStreamBi {
	return &grpcBiContextStream{
		sr: sr,
	}
}

func (g *grpcBiContextStream) Recv() (adapters.IGrpcDto, error) {
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

func NewClientStreamContext(sr pb.Grpc_ClientStreamServer) adapters.IStreamClient {
	return &grpcContextClientStream{
		sr: sr,
	}
}

func (g *grpcContextClientStream) Recv() (adapters.IGrpcDto, error) {
	return g.sr.Recv()
}

type grpcContextServerStream struct {
	sr pb.Grpc_ServerStreamServer
}

func NewServerStreamContext(sr pb.Grpc_ServerStreamServer) adapters.IStreamServer {
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
