package input

import (
	"context"
	grpcconn "dominus-project/app/interactors/grpc_conn"
	pb "dominus-project/app/interfaces/grpc/proto/builder"
	"io"

	"google.golang.org/grpc"
)

type grpcController struct {
	pb.UnimplementedGrpcServer
	uc grpcconn.GrpcServiceInt
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
	ctx := newClientStreamContext(stream)
	err := s.uc.StreamClientConn(ctx)
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
	return s.uc.StreamServerConn(ms, ctx)
}

// Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.Grpc_BidirectionalStreamServer) error {
	ctx := newBiStreamConn(stream)
	return s.uc.StreamBiConn(ctx)
}

func (s *grpcController) runQueue() {
	go s.uc.RunQueue()
}

func NewGrpcAPI(opts []grpc.ServerOption, uc grpcconn.GrpcServiceInt) *grpc.Server {
	s := grpc.NewServer(opts...)
	gsrv := &grpcController{uc: uc}
	pb.RegisterGrpcServer(s, gsrv)
	gsrv.runQueue()
	return s
}
