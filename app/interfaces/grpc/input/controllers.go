package input

import (
	"context"
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"
	"io"

	"google.golang.org/grpc"
)

type grpcController struct {
	pb.UnimplementedGrpcServer
	inter interactors.InteractorInt
}

// Receives simple messages from client
func (s *grpcController) Simple(_ context.Context, ms *pb.RequestMessage) (*pb.Response, error) {
	conn := s.inter.NewConnection()
	if err := conn.SimpleConn(ms); err != nil {
		return &pb.Response{Status: uint32(500), Message: err.Error()}, err
	}
	return &pb.Response{Status: uint32(202), Message: "[+]Accepted"}, nil
}

// Receives array messages from client
func (s *grpcController) ClientStream(stream pb.Grpc_ClientStreamServer) error {
	conn := s.inter.NewConnection()
	ctx := newClientStreamContext(stream)
	err := conn.StreamClientConn(ctx)
	if err != io.EOF {
		return stream.SendAndClose(&pb.Response{
			Status:  uint32(500),
			Message: err.Error(),
		})
	}
	return stream.SendAndClose(&pb.Response{
		Status: uint32(202),
		Message: "[+]Connection closed",
	})
}

// Sends array messages to client
func (s *grpcController) ServerStream(ms *pb.RequestMessage, stream pb.Grpc_ServerStreamServer) error {
	conn := s.inter.NewConnection()
	ctx := newServerStreamContext(stream)
	return conn.StreamServerConn(ms, ctx)
}

// Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.Grpc_BidirectionalStreamServer) error {
	conn := s.inter.NewConnection()
	ctx := newBiStreamConn(stream)
	return conn.StreamBiConn(ctx)
}

func NewGrpcController(opts []grpc.ServerOption, i interactors.InteractorInt) *grpc.Server {
	s := grpc.NewServer(opts...)
	pb.RegisterGrpcServer(s, &grpcController{inter: i})
	return s
}
