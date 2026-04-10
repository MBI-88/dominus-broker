package sqs_flow_test

import (
	"context"
	"testing"

	"dominus-broker/internal/application/usecases/sqs"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/infrastructure/grpc/inbound"
	"dominus-broker/mocks"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestAckFlow(t *testing.T) {
	const streamID = "stream-ack-flow"
	const groupID = "ack-flow-group"

	t.Run("Ack end-to-end via gRPC after Consumer", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		eventMock := mocks.NewMockEvent(ctrl)
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		mem := newMemoryClient(t, s, streamID, eventMock)
		ctx := context.Background()

		if err := mem.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
		}

		msg := entities.NewMessageWithID([]byte("ack-flow-payload"))
		if err := mem.SendMessage(ctx, msg); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}

		svc := sqs.NewSQS(mem)
		lis := bufconn.Listen(buffSize)
		server := grpc.NewServer()
		inbound.NewSqsAPI(server, svc, eventMock)
		go func() {
			if err := server.Serve(lis); err != nil {
				t.Logf("Serve: %v", err)
			}
		}()
		t.Cleanup(server.GracefulStop)

		client := newSqsAPIClient(t, lis)

		cons, err := client.Consumer(ctx, &pb.ConsumerRequest{
			GroupId:  groupID,
			WorkerId: "worker-ack-1",
		})
		if err != nil {
			t.Fatalf("Consumer: %v", err)
		}
		if string(cons.GetMessage()) != "ack-flow-payload" {
			t.Fatalf("Consumer message: got %q", cons.GetMessage())
		}
		if cons.GetMessageId() == "" {
			t.Fatal("Consumer returned empty MessageId")
		}

		ack, err := client.Ack(ctx, &pb.ConsumerRequest{
			MessageId: cons.GetMessageId(),
			GroupId:   groupID,
			WorkerId:  "worker-ack-1",
		})
		if err != nil {
			t.Fatalf("Ack: %v", err)
		}
		if ack.GetMessageId() != cons.GetMessageId() {
			t.Fatalf("Ack MessageId: got %q want %q", ack.GetMessageId(), cons.GetMessageId())
		}

		rdb := testRedisClient(t, s)
		t.Cleanup(func() { _ = rdb.Close() })

		pending, err := rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
			Stream: streamID,
			Group:  groupID,
			Start:  "-",
			End:    "+",
			Count:  10,
		}).Result()
		if err != nil {
			t.Fatalf("XPendingExt: %v", err)
		}
		for _, p := range pending {
			if p.ID == cons.GetMessageId() {
				t.Fatalf("message %q should be acked, still pending", p.ID)
			}
		}
	})

	t.Run("Ack rejects empty message id", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		eventMock := mocks.NewMockEvent(ctrl)
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		mem := newMemoryClient(t, s, streamID, eventMock)
		svc := sqs.NewSQS(mem)
		lis := bufconn.Listen(buffSize)
		server := grpc.NewServer()
		inbound.NewSqsAPI(server, svc, eventMock)
		go func() {
			if err := server.Serve(lis); err != nil {
				t.Logf("Serve: %v", err)
			}
		}()
		t.Cleanup(server.GracefulStop)

		client := newSqsAPIClient(t, lis)
		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			MessageId: "",
			GroupId:   groupID,
			WorkerId:  "worker-1",
		})
		if err == nil {
			t.Fatal("expected error for empty MessageId")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.NotFound {
			t.Fatalf("expected NotFound, got %v (%v)", st.Code(), err)
		}
	})
}
