package sqs_flow_test

import (
	"context"
	"testing"
	"time"

	"dominus-broker/internal/application/use_cases/sqs"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/infrastructure/grpc/inbound"
	"dominus-broker/mocks"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"github.com/alicebob/miniredis/v2"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestConsumerFlow(t *testing.T) {
	const streamID = "stream-consumer-flow"
	const groupID = "consumer-flow-group"

	t.Run("Consumer end-to-end via gRPC", func(t *testing.T) {
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

		sent := entities.NewMessageWithID([]byte("consumer-flow-payload"))
		if err := mem.SendMessage(ctx, sent); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}
		streamEntries, err := s.Stream(streamID)
		if err != nil {
			t.Fatalf("Stream: %v", err)
		}
		if len(streamEntries) != 1 {
			t.Fatalf("expected 1 stream entry after SendMessage, got %d", len(streamEntries))
		}
		wantMessageID := streamEntries[0].ID

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
		before := time.Now().Add(-time.Second)

		resp, err := client.Consumer(ctx, &pb.ConsumerRequest{
			GroupId:  groupID,
			WorkerId: "worker-consumer-1",
		})
		if err != nil {
			t.Fatalf("Consumer: %v", err)
		}
		if string(resp.GetMessage()) != "consumer-flow-payload" {
			t.Fatalf("message body: got %q", resp.GetMessage())
		}
		if resp.GetMessageId() != wantMessageID {
			t.Fatalf("message id: got %q want %q", resp.GetMessageId(), wantMessageID)
		}
		if resp.GetDate() == nil {
			t.Fatal("expected non-nil Date")
		}
		gotAt := resp.GetDate().AsTime()
		if gotAt.Before(before) || gotAt.After(time.Now().Add(time.Minute)) {
			t.Fatalf("Date out of range: %v", gotAt)
		}
	})

	t.Run("Consumer second read returns next message", func(t *testing.T) {
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

		// IDs fijos estrictamente crecientes: dos NewMessage() en el mismo ms compartirían
		// el mismo message_id (UnixMilli-0) y el segundo XADD fallaría en Redis.
		for _, id := range []string{"10-0", "11-0"} {
			m := entities.NewMessageWithID([]byte("m-" + id))
			m.SetMessageId(id)
			if err := mem.SendMessage(ctx, m); err != nil {
				t.Fatalf("SendMessage %s: %v", id, err)
			}
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

		r1, err := client.Consumer(ctx, &pb.ConsumerRequest{GroupId: groupID, WorkerId: "w-multi"})
		if err != nil {
			t.Fatalf("Consumer 1: %v", err)
		}
		r2, err := client.Consumer(ctx, &pb.ConsumerRequest{GroupId: groupID, WorkerId: "w-multi"})
		if err != nil {
			t.Fatalf("Consumer 2: %v", err)
		}
		ids := map[string]bool{r1.GetMessageId(): true, r2.GetMessageId(): true}
		if !ids["10-0"] || !ids["11-0"] {
			t.Fatalf("expected ids 10-0 and 11-0, got %q and %q", r1.GetMessageId(), r2.GetMessageId())
		}
	})

	t.Run("Consumer rejects empty group id", func(t *testing.T) {
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

		_, err := client.Consumer(context.Background(), &pb.ConsumerRequest{
			GroupId:  "",
			WorkerId: "w1",
		})
		if err == nil {
			t.Fatal("expected error")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.NotFound {
			t.Fatalf("expected NotFound, got %v (%v)", st.Code(), err)
		}
	})

	t.Run("Consumer rejects empty worker id", func(t *testing.T) {
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

		_, err := client.Consumer(context.Background(), &pb.ConsumerRequest{
			GroupId:  groupID,
			WorkerId: "",
		})
		if err == nil {
			t.Fatal("expected error")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.NotFound {
			t.Fatalf("expected NotFound, got %v (%v)", st.Code(), err)
		}
	})

	t.Run("Consumer fails when stream has no new messages", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		eventMock := mocks.NewMockEvent(ctrl)
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()
		mem := newMemoryClient(t, s, streamID, eventMock)
		if err := mem.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
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

		_, err := client.Consumer(context.Background(), &pb.ConsumerRequest{
			GroupId:  groupID,
			WorkerId: "w-empty",
		})
		if err == nil {
			t.Fatal("expected error when no messages")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Aborted {
			t.Fatalf("expected Aborted, got %v (%v)", st.Code(), err)
		}
	})
}
