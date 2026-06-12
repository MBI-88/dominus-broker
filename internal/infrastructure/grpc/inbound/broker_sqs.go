package inbound

import (
	"context"
	"dominus-broker/internal/application/usecases/ack"
	"dominus-broker/internal/application/usecases/consumer"
	"dominus-broker/internal/application/usecases/producer"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type brokerSqs struct {
	pb.UnimplementedSqsAPIServer
	p   producer.ProducerUseCase
	c   consumer.ConsumerUseCase
	a   ack.AckUseCase
	log event.Event
}

func NewBrokerSqs(
	server *grpc.Server,
	p producer.ProducerUseCase,
	c consumer.ConsumerUseCase,
	a ack.AckUseCase,
	log event.Event,
) {
	gsrv := &brokerSqs{
		p:   p,
		c:   c,
		a:   a,
		log: log,
	}
	pb.RegisterSqsAPIServer(server, gsrv)
}

// Receives simple messages from client
func (s *brokerSqs) Producer(ctx context.Context, ms *pb.ProducerRequest) (*pb.ProducerResponse, error) {
	s.log.WriteLog(ctx, enum.DEBUG, "inbound.Producer", enum.REQUEST_OK)

	if len(ms.GetPayload()) == 0 {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Producer.GetPayload", enum.INVALID_PAYLOAD)
		return nil, status.Error(codes.OutOfRange, enum.INVALID_PAYLOAD)
	}

	if err := s.p.Producer(ctx, ms); err != nil {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Producer.Producer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.ProducerResponse{Status: 0}, nil
}

func (s *brokerSqs) Consumer(ctx context.Context, ms *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	s.log.WriteLog(ctx, enum.DEBUG, "inbound.Consumer", enum.REQUEST_OK)

	if ms.GetGroupId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Consumer.GetGroupId", enum.GROUP_ID)
		return nil, status.Error(codes.NotFound, enum.GROUP_ID)
	}

	if ms.GetWorkerId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Consumer.GetWorkerId", enum.WORKER_ID)
		return nil, status.Error(codes.NotFound, enum.WORKER_ID)
	}

	response, err := s.c.Consumer(ctx, ms)
	if err != nil {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Consumer.Consumer", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.ConsumerResponse{
		MessageId: response.GetMessageId(),
		Message:   response.GetMessage(),
		Date:      timestamppb.New(response.GetCreatedAt()),
	}, nil
}

func (s *brokerSqs) Ack(ctx context.Context, ms *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	s.log.WriteLog(ctx, enum.DEBUG, "inbound.Ack", enum.REQUEST_OK)

	if ms.GetMessageId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Ack.GetMessageId", enum.INVALID_ID)
		return nil, status.Error(codes.NotFound, enum.INVALID_ID)
	}

	if ms.GetGroupId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Ack.GetGroupId", enum.GROUP_ID)
		return nil, status.Error(codes.NotFound, enum.GROUP_ID)
	}

	if ms.GetWorkerId() == "" {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Ack.GetWorkerId", enum.WORKER_ID)
		return nil, status.Error(codes.NotFound, enum.WORKER_ID)
	}

	if err := s.a.Ack(ctx, ms); err != nil {
		s.log.WriteLog(ctx, enum.ERROR, "inbound.Ack.Ack", err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.ConsumerResponse{
		MessageId: ms.GetMessageId(),
		Date:      timestamppb.Now(),
	}, nil
}
