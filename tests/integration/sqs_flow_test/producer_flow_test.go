package sqs_flow_test

import (
	"context"
	"testing"

	"dominus-broker/internal/application/usecases/sqs"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/grpc/inbound"
	"dominus-broker/internal/infrastructure/redis/cmemory"
	"dominus-broker/mocks"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"github.com/alicebob/miniredis/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func payloadFromStreamValues(values []string) (string, bool) {
	for i := 0; i+1 < len(values); i += 2 {
		if values[i] == enum.PAYLOAD {
			return values[i+1], true
		}
	}
	return "", false
}

func TestProducerFlow(t *testing.T) {
	const streamID = "stream-producer-flow"

	t.Run("Producer end-to-end via gRPC", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		eventMock := mocks.NewMockEvent(ctrl)
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		mem := newMemoryClient(t, s, streamID, eventMock)
		ctx := context.Background()

		svc := sqs.NewSqs(mem)
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
		want := []byte("producer-flow-payload")
		resp, err := client.Producer(ctx, &pb.ProducerRequest{Payload: want})
		if err != nil {
			t.Fatalf("Producer: %v", err)
		}
		if resp.GetStatus() != 0 {
			t.Fatalf("Status: got %d want 0", resp.GetStatus())
		}

		entries, err := s.Stream(streamID)
		if err != nil {
			t.Fatalf("Stream: %v", err)
		}
		if len(entries) != 1 {
			t.Fatalf("expected 1 stream entry, got %d", len(entries))
		}
		raw, ok := payloadFromStreamValues(entries[0].Values)
		if !ok {
			t.Fatalf("missing %q in values: %v", enum.PAYLOAD, entries[0].Values)
		}
		var got cmemory.MessageDto
		if err := jsoniter.Unmarshal([]byte(raw), &got); err != nil {
			t.Fatalf("Unmarshal payload: %v", err)
		}

		entity := entities.NewMessage()
		entity.SetCreateAt(got.CreatedAt)
		entity.SetMessage(got.Message)

		if !entity.SetMessageId(got.MessageId) {
			t.Fatalf("invalid id %s", got.MessageId)
		}

		if string(entity.GetMessage()) != string(want) {
			t.Fatalf("message body: got %q want %q", entity.GetMessage(), want)
		}
		streamEntryID := entries[0].ID
		if entity.GetMessageId() != streamEntryID {
			t.Fatalf("message_id in JSON %q != Redis stream entry id %q", entity.GetMessageId(), streamEntryID)
		}
	})

	t.Run("Producer rejects empty payload", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		eventMock := mocks.NewMockEvent(ctrl)
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		mem := newMemoryClient(t, s, streamID, eventMock)
		svc := sqs.NewSqs(mem)
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

		_, err := client.Producer(context.Background(), &pb.ProducerRequest{Payload: nil})
		if err == nil {
			t.Fatal("expected error for nil payload")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.OutOfRange {
			t.Fatalf("expected OutOfRange, got %v (%v)", st.Code(), err)
		}
	})

	t.Run("Producer fails when Redis rejects XADD", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		eventMock := mocks.NewMockEvent(ctrl)
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		mem := newMemoryClient(t, s, streamID, eventMock)
		svc := sqs.NewSqs(mem)
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

		s.SetError("simulated XADD failure")
		_, err := client.Producer(context.Background(), &pb.ProducerRequest{Payload: []byte("x")})
		if err == nil {
			t.Fatal("expected error when Redis returns failure")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Aborted {
			t.Fatalf("expected Aborted, got %v (%v)", st.Code(), err)
		}
	})
}
