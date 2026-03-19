package inbound

import (
	"context"
	"dominus-project/internal/infrastructure/event"
	"dominus-project/internal/infrastructure/grpc/mappers"
	"dominus-project/internal/application/use_cases/broker"
	"dominus-project/internal/application/use_cases/queue"
	"dominus-project/internal/infrastructure/enum"
	"fmt"
	"io"

	pb "github.com/PR0C0D3-MBI/dominus-proto-definition/dominus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcController struct {
	pb.UnimplementedAPIServer
	br  broker.Broker
	q   queue.Queue
	log event.Event
}

func NewGrpcAPI(opts []grpc.ServerOption, br broker.Broker, q queue.Queue, log event.Event) *grpc.Server {
	s := grpc.NewServer(opts...)
	gsrv := &grpcController{br: br, log: log, q: q}
	pb.RegisterAPIServer(s, gsrv)
	reflection.Register(s)
	return s
}

// Receives simple messages from client
func (s *grpcController) Producer(ctx context.Context, ms *pb.ProducerRequest) (*pb.ProducerResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "Producer", enum.REQUEST_OK)
	if err := s.q.Producer(ctx, ms); err != nil {
		go s.log.WriteLog(ctx, enum.ERROR, "Producer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.ProducerResponse{Status: enum.OK}, nil
}

func (s *grpcController) Consumer(ctx context.Context, ms *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "Consumer", enum.REQUEST_OK)
	response, err := s.q.Consumer(ctx)
	if err != nil {
		go s.log.WriteLog(ctx, enum.ERROR, "Consumer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.ConsumerResponse{
		Id:      response.ID.String(),
		Message: response.Message,
		Hidden:  response.Hidden,
		Date:    timestamppb.New(response.CreatedAt),
	}, nil
}

func (s *grpcController) ConsumerDLT(ctx context.Context, ms *pb.ConsumerDeleteRequest) (*pb.ConsumerDeleteResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "ConsumerDLT", enum.REQUEST_OK)
	if err := s.q.ConsumerDLT(ctx, ms); err != nil {
		go s.log.WriteLog(ctx, enum.ERROR, "ConsumerDLT", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.ConsumerDeleteResponse{
		Status: enum.OK,
		Description: enum.DESCRIPTION_DELETE_MESSAGE,
	}, nil
}

// Receives array messages from client
func (s *grpcController) ClientStream(stream pb.API_ClientStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "ClientStream", enum.REQUEST_OK)
	ctx := mappers.NewClientStreamContext(stream)
	err := s.br.StreamClientConn(ctx)
	if err != io.EOF {
		go s.log.WriteLog(stream.Context(), enum.ERROR, "ClientStream", err.Error())
		return status.Error(codes.Aborted, err.Error())
	}
	return stream.SendAndClose(&pb.StreamResponseMessage{
		Status: enum.OK,
	})
}

// Sends array messages to client
func (s *grpcController) ServerStream(ms *pb.StreamRequestMessage, stream pb.API_ServerStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "ServerStream", enum.REQUEST_OK)
	ctx := mappers.NewServerStreamContext(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.br.StreamServerConn(ms, ctx)))
}

// Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.API_BidirectionalStreamServer) error {
	go s.log.WriteLog(stream.Context(), enum.DEBUG, "BidirectionalStream", enum.REQUEST_OK)
	ctx := mappers.NewBiStreamConn(stream)
	return status.Error(codes.Aborted, fmt.Sprintf("%s", s.br.StreamBiConn(ctx)))
}
