package inbound

import (
	"context"
	"dominus-project/internal/application/use_cases/sqs"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/event"

	pb "github.com/PR0C0D3-MBI/dominus-proto-definition/dominus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sqsAPI struct {
	pb.UnimplementedSqsAPIServer
	q   sqs.SQS
	log event.Event
}

func NewSqsAPI(server *grpc.Server, q sqs.SQS, log event.Event) {
	gsrv := &sqsAPI{log: log, q: q}
	pb.RegisterSqsAPIServer(server, gsrv)
}

// Receives simple messages from client
func (s *sqsAPI) Producer(ctx context.Context, ms *pb.ProducerRequest) (*pb.ProducerResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "Producer", enum.REQUEST_OK)

	if len(ms.GetPayload()) == 0 {
		go s.log.WriteLog(ctx, enum.ERROR, "Producer", enum.INVALID_PAYLOAD)
		return nil, status.Error(codes.OutOfRange, enum.INVALID_PAYLOAD)
	}

	if err := s.q.Producer(ctx, ms); err != nil {
		go s.log.WriteLog(ctx, enum.ERROR, "Producer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.ProducerResponse{Status: int64(codes.OK)}, nil
}

func (s *sqsAPI) Consumer(ctx context.Context, ms *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "Consumer", enum.REQUEST_OK)
	response, err := s.q.Consumer(ctx)
	if err != nil {
		go s.log.WriteLog(ctx, enum.ERROR, "Consumer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.ConsumerResponse{
		Id:      response.GetID(),
		Message: response.GetMessage(),
		Hidden:  response.GetHidden(),
		Date:    timestamppb.New(response.GetCreatedAt()),
	}, nil
}

func (s *sqsAPI) ConsumerDLT(ctx context.Context, ms *pb.ConsumerDeleteRequest) (*pb.ConsumerDeleteResponse, error) {
	go s.log.WriteLog(ctx, enum.DEBUG, "ConsumerDLT", enum.REQUEST_OK)

	if ms.GetId() == "" {
		go s.log.WriteLog(ctx, enum.ERROR, "ConsumerDLT", enum.INVALID_ID)
		return nil, status.Error(codes.NotFound, enum.INVALID_ID)
	}
	if err := s.q.ConsumerDLT(ctx, ms); err != nil {
		go s.log.WriteLog(ctx, enum.ERROR, "ConsumerDLT", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.ConsumerDeleteResponse{
		Status: int64(codes.OK),
		Description: enum.DESCRIPTION_DELETE_MESSAGE,
	}, nil
}