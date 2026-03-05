package input

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/enum"
	"dominus-project/internal/orchestrators/broker"
	"dominus-project/internal/orchestrators/queue"
	"dominus-project/internal/infrastructure/grpc/dto"
	"fmt"
	"io"

	pb "github.com/PR0C0D3-MBI/dominus-proto-definition/dominus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcController struct {
	pb.UnimplementedAPIServer
	br  broker.Broker
	q   queue.Queue
	log adapters.Logs
}

func NewGrpcAPI(opts []grpc.ServerOption, br broker.Broker, q queue.Queue, log adapters.Logs) *grpc.Server {
	s := grpc.NewServer(opts...)
	gsrv := &grpcController{br: br, log: log, q: q}
	pb.RegisterAPIServer(s, gsrv)
	reflection.Register(s)
	return s
}

// Receives simple messages from client
func (s *grpcController) Provider(ctx context.Context, ms *pb.ProviderRequest) (*pb.ProviderResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "Provider", "request accepted")
	if err := s.q.Provider(ctx, ms); err != nil {
		go s.log.WriteLog(ctx, enum.ERROR, "Provider", err.Error())
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s", err))
	}
	return &pb.ProviderResponse{Status: 0}, nil
}

func (s *grpcController) Consumer(ctx context.Context, ms *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "Consumer", "request accepted")
	return &pb.ConsumerResponse{}, nil
}

func (s *grpcController) ConsumerDLT(ctx context.Context, ms *pb.ConsumerDeleteRequest) (*pb.ConsumerDeleteResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "ConsumerDLT", "request accepted")
	return &pb.ConsumerDeleteResponse{}, nil
}

// Receives array messages from client
func (s *grpcController) ClientStream(stream pb.API_ClientStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "ClientStream", "request accepted")
	ctx := dto.NewClientStreamContext(stream)
	err := s.br.StreamClientConn(ctx)
	if err != io.EOF {
		go s.log.WriteLog(stream.Context(), enum.ERROR, "ClientStream", err.Error())
		return status.Error(codes.Canceled, fmt.Sprintf("%s", err))
	}
	return stream.SendAndClose(&pb.StreamResponseMessage{
		Status: 0,
	})
}

// Sends array messages to client
func (s *grpcController) ServerStream(ms *pb.StreamRequestMessage, stream pb.API_ServerStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "ServerStream", "request accepted")
	ctx := dto.NewServerStreamContext(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.br.StreamServerConn(ms, ctx)))
}

// Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.API_BidirectionalStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "BidirectionalStream", "request accepted")
	ctx := dto.NewBiStreamConn(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.br.StreamBiConn(ctx)))
}
