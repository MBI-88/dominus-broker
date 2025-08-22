package input

import (
	"context"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/internal/interfaces/grpc/dto"
	"fmt"
	"io"

	pb "github.com/PR0C0D3-MBI/dominus-proto-definition/src/dominus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s", err))
	}
	return &pb.Response{Status: 200, Message: "[+] Accepted"}, nil
}

// Receives array messages from client
func (s *grpcController) ClientStream(stream pb.Grpc_ClientStreamServer) error {
	ctx := dto.NewClientStreamContext(stream)
	err := s.uc.StreamClientConn(ctx)
	if err != io.EOF {
		return status.Error(codes.Canceled, fmt.Sprintf("%s", err))
	}
	return stream.SendAndClose(&pb.Response{
		Status:  200,
		Message: "[*]Connection closed",
	})
}

// Sends array messages to client
func (s *grpcController) ServerStream(ms *pb.RequestMessage, stream pb.Grpc_ServerStreamServer) error {
	ctx := dto.NewServerStreamContext(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.uc.StreamServerConn(ms, ctx))) 
}

// Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.Grpc_BidirectionalStreamServer) error {
	ctx := dto.NewBiStreamConn(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.uc.StreamBiConn(ctx))) 
}

func (s *grpcController) runQueue(close <-chan struct{}) {
	go s.uc.RunQueue(close)
}
