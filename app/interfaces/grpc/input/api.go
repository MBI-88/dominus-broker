package input

import (
	"context"
	"dominus-project/app/interactors"
	pb "dominus-project/app/interfaces/grpc/proto/builder"
	"io"

	"google.golang.org/grpc"
)

type grpcController struct {
	pb.UnimplementedGrpcServer
	service interactors.GrpcServiceInt
}

// Receives simple messages from client
func (s *grpcController) Simple(_ context.Context, ms *pb.RequestMessage) (*pb.Response, error) {
	if err := s.service.SimpleConn(ms); err != nil {
		return &pb.Response{Status: uint32(500), Message: err.Error()}, err
	}
	return &pb.Response{Status: uint32(200), Message: "[*]Accepted"}, nil
}

// Receives array messages from client
func (s *grpcController) ClientStream(stream pb.Grpc_ClientStreamServer) error {
	ctx := newClientStreamContext(stream)
	err := s.service.StreamClientConn(ctx)
	if err != io.EOF {
		return stream.SendAndClose(&pb.Response{
			Status:  uint32(500),
			Message: err.Error(),
		})
	}
	return stream.SendAndClose(&pb.Response{
		Status: uint32(200),
		Message: "[*]Connection closed",
	})
}

// Sends array messages to client
func (s *grpcController) ServerStream(ms *pb.RequestMessage, stream pb.Grpc_ServerStreamServer) error {
	ctx := newServerStreamContext(stream)
	return s.service.StreamServerConn(ms, ctx)
}

// Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.Grpc_BidirectionalStreamServer) error {
	ctx := newBiStreamConn(stream)
	return s.service.StreamBiConn(ctx)
}

func NewGrpcAPI(opts []grpc.ServerOption, i interactors.GrpcServiceInt) *grpc.Server {
	s := grpc.NewServer(opts...)
	pb.RegisterGrpcServer(s, &grpcController{service: i})
	return s
}
