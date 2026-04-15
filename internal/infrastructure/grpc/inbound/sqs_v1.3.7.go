package inbound

import (
	"context"
	"dominus-broker/internal/application/usecases/sqs"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sqsAPI struct {
	pb.UnimplementedSqsAPIServer
	q   sqs.Sqs
	log event.Event
}

func NewSqsAPI(server *grpc.Server, q sqs.Sqs, log event.Event) {
	gsrv := &sqsAPI{log: log, q: q}
	pb.RegisterSqsAPIServer(server, gsrv)
}

// Receives simple messages from client
func (s *sqsAPI) Producer(ctx context.Context, ms *pb.ProducerRequest) (*pb.ProducerResponse, error) {
	s.log.WriteLog(ctx, enum.DEBUG, "sqsAPI.Producer", enum.REQUEST_OK)

	if len(ms.GetPayload()) == 0 {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Producer.GetPayload", enum.INVALID_PAYLOAD)
		return nil, status.Error(codes.OutOfRange, enum.INVALID_PAYLOAD)
	}

	if err := s.q.Producer(ctx, ms); err != nil {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Producer.Producer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.ProducerResponse{Status: 0}, nil
}

func (s *sqsAPI) Consumer(ctx context.Context, ms *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	s.log.WriteLog(ctx, enum.DEBUG, "sqsAPI.Consumer", enum.REQUEST_OK)

	if ms.GetGroupId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Consumer.GetGroupId", enum.GROUP_ID)
		return nil, status.Error(codes.NotFound, enum.GROUP_ID)
	}

	if ms.GetWorkerId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Consumer.GetWorkerId", enum.WORKER_ID)
		return nil, status.Error(codes.NotFound, enum.WORKER_ID)
	}

	response, err := s.q.Consumer(ctx, ms)
	if err != nil {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Consumer.Consumer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.ConsumerResponse{
		MessageId: response.GetMessageId(),
		Message:   response.GetMessage(),
		Date:      timestamppb.New(response.GetCreatedAt()),
	}, nil
}

func (s *sqsAPI) Ack(ctx context.Context, ms *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	s.log.WriteLog(ctx, enum.DEBUG, "sqsAPI.Ack", enum.REQUEST_OK)

	if ms.GetMessageId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Ack.GetMessageId", enum.INVALID_ID)
		return nil, status.Error(codes.NotFound, enum.INVALID_ID)
	}

	if ms.GetGroupId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Ack.GetGroupId", enum.GROUP_ID)
		return nil, status.Error(codes.NotFound, enum.GROUP_ID)
	}

	if ms.GetWorkerId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Ack.GetWorkerId", enum.WORKER_ID)
		return nil, status.Error(codes.NotFound, enum.WORKER_ID)
	}

	if err := s.q.Ack(ctx, ms); err != nil {
		s.log.WriteLog(ctx, enum.ERROR, "sqsAPI.Ack.Ack", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.ConsumerResponse{
		MessageId: ms.GetMessageId(),
		Date:      timestamppb.Now(),
	}, nil
}
