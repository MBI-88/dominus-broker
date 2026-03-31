package inbound_test

import (
	"context"
	appsqs "dominus-broker/internal/application/use_cases/sqs"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/infrastructure/grpc/inbound"
	"dominus-broker/mocks"
	"fmt"
	"testing"
	"time"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// newSqsAPITestClient returns a gRPC client that talks to a server registered on lis.
// The address is ignored when using WithContextDialer(bufDialer(lis)).
func newSqsAPITestClient(t *testing.T, lis *bufconn.Listener) pb.SqsAPIClient {
	t.Helper()
	opts := []grpc.DialOption{
		grpc.WithContextDialer(bufDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	// passthrough avoids DNS resolution; WithContextDialer handles the real connection via bufconn.
	conn, err := grpc.NewClient("passthrough:///bufnet", opts...)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return pb.NewSqsAPIClient(conn)
}

func TestProducer(t *testing.T) {

	t.Run("Producer OK", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		sqsMock.EXPECT().
			Producer(gomock.All(), gomock.All()).
			Return(nil)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		resp, err := client.Producer(context.Background(), &pb.ProducerRequest{
			Payload: []byte("test-1-ok"),
		})

		if err != nil {
			t.Fatal(err)
		}

		if resp.Status != 0 {
			t.Fatalf("Expected %d go %d\n", 0, resp.Status)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})

	t.Run("Procuder invalid payload", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Producer(context.Background(), &pb.ProducerRequest{
			Payload: []byte(""),
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})

	t.Run("Producer error", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		sqsMock.EXPECT().
			Producer(gomock.All(), gomock.All()).
			Return(fmt.Errorf("connection error")).
			AnyTimes()

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Producer(context.Background(), &pb.ProducerRequest{
			Payload: []byte("test-1-ok"),
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})
}

func TestConsumer(t *testing.T) {

	t.Run("Consumer oK", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		wantID := "1730000000000-0"
		sqsMock.EXPECT().
			Consumer(gomock.All(), gomock.All()).
			Return(&entities.Message{
				Message:   []byte("test-done"),
				MeesageId: wantID,
				CreatedAt: time.Now(),
			}, nil).
			AnyTimes()

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		resp, err := client.Consumer(context.Background(), &pb.ConsumerRequest{
			WorkerId: "consumer-1",
			GroupId:  "consumer-g",
		})

		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected different from nil got %v\n", resp)
		}
		if string(resp.GetMessage()) != "test-done" {
			t.Fatalf("message body: got %q want %q", resp.GetMessage(), "test-done")
		}
		if resp.GetMessageId() != wantID {
			t.Fatalf("message id: got %q want %q", resp.GetMessageId(), wantID)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})

	t.Run("Consumer empty group id", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Consumer(context.Background(), &pb.ConsumerRequest{
			WorkerId: "consumer-1",
			GroupId:  "",
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})

	})

	t.Run("Consumer empty work id", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Consumer(context.Background(), &pb.ConsumerRequest{
			WorkerId: "",
			GroupId:  "consumer-g",
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})

	t.Run("Consumer error", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		sqsMock.EXPECT().
			Consumer(gomock.All(), gomock.All()).
			Return(nil, fmt.Errorf("connection error"))

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Consumer(context.Background(), &pb.ConsumerRequest{
			WorkerId: "consumer-1",
			GroupId:  "consumer-g",
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})
}

func TestAck(t *testing.T) {

	t.Run("Ack ok", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		sqsMock.EXPECT().
			Ack(gomock.All(), gomock.All()).
			Return(nil)

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		const ackMsgID = "1700000000001-0"
		resp, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			WorkerId:  "consumer-1",
			MessageId: ackMsgID,
			GroupId:   "consumer-g",
		})

		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected %v got %v\n", nil, resp)
		}
		if resp.GetMessageId() != ackMsgID {
			t.Fatalf("Ack response MessageId: got %q want %q", resp.GetMessageId(), ackMsgID)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})

	t.Run("Ack empty message id", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId:   "consumer-g",
			WorkerId:  "consumer-1",
			MessageId: "",
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})

	})

	t.Run("Ack empty group id", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId:   "",
			WorkerId:  "consumer-1",
			MessageId: "1700000000002-0",
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})

	})

	t.Run("Ack empty worker id", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId:   "consumer-g",
			WorkerId:  "",
			MessageId: "1700000000003-0",
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})

	t.Run("Ack error", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		sqsMock := mocks.NewMockSQS(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		sqsMock.EXPECT().
			Ack(gomock.All(), gomock.All()).
			Return(fmt.Errorf("connection error"))

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId:   "consumer-g",
			WorkerId:  "consumer-1",
			MessageId: "1700000000004-0",
		})

		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})

	t.Run("Ack invalid message id with real use case", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		evnetMock := mocks.NewMockEvent(ctrl)
		memMock := mocks.NewMockMemoryClient(ctrl)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		svc := appsqs.NewSQS(memMock)
		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, svc, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId:   "consumer-g",
			WorkerId:  "consumer-1",
			MessageId: "not-a-redis-stream-id",
		})
		if err == nil {
			t.Fatal("expected error for invalid MessageId")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Aborted {
			t.Fatalf("expected Aborted, got %v (%v)", st.Code(), err)
		}
		if st.Message() != "invalid messageId" {
			t.Fatalf("expected invalid messageId, got %q", st.Message())
		}

		t.Cleanup(func() {
			server.GracefulStop()
			ctrl.Finish()
		})
	})
}
