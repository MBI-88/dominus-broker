package input

import (
	"context"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/internal/interfaces/grpc/dto"
	pb "dominus-project/internal/interfaces/grpc/proto/builder"
	"io"

	"google.golang.org/grpc"
)

type grpcController struct {
	pb.UnimplementedGrpcServer
	uc grpcconn.GrpcService
}

func NewGrpcAPI(opts []grpc.ServerOption, uc grpcconn.GrpcService, close <-chan struct{}) *grpc.Server {
	s := grpc.NewServer(opts...)
	gsrv := &grpcController{uc: uc}
	pb.RegisterGrpcServer(s, gsrv)
	gsrv.runQueue(close)
	return s
}

// Receives simple messages from client
func (s *grpcController) Simple(_ context.Context, ms *pb.RequestMessage) (*pb.Response, error) {
	if err := s.uc.SimpleConn(ms); err != nil {
		return &pb.Response{Status: uint32(500), Message: err.Error()}, err
	}
	return &pb.Response{Status: uint32(200), Message: "[+] Accepted"}, nil
}

// Receives array messages from client
func (s *grpcController) ClientStream(stream pb.Grpc_ClientStreamServer) error {
	ctx := dto.NewClientStreamContext(stream)
	err := s.uc.StreamClientConn(ctx)
	if err != io.EOF {
		return stream.SendAndClose(&pb.Response{
			Status:  uint32(500),
			Message: err.Error(),
		})
	}
	return stream.SendAndClose(&pb.Response{
		Status:  uint32(200),
		Message: "[*]Connection closed",
	})
}

// Sends array messages to client
func (s *grpcController) ServerStream(ms *pb.RequestMessage, stream pb.Grpc_ServerStreamServer) error {
	ctx := dto.NewServerStreamContext(stream)
	return s.uc.StreamServerConn(ms, ctx)
}

// Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.Grpc_BidirectionalStreamServer) error {
	ctx := dto.NewBiStreamConn(stream)
	return s.uc.StreamBiConn(ctx)
}

func (s *grpcController) runQueue(close <-chan struct{}) {
	go s.uc.RunQueue(close)
}
