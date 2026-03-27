package inbound_test

import (
	"context"
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/infrastructure/grpc/inbound"
	"dominus-project/mocks"
	"fmt"
	"testing"
	"time"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

		sqsMock.EXPECT().
			Consumer(gomock.All(), gomock.All()).
			Return(&entities.Message{
				Message:   []byte("test-done"),
				MeesageId: "123456789",
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

		resp, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			WorkerId:  "consumer-1",
			MessageId: "message-1",
			GroupId:   "consumer-g",
		})

		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected %v got %v\n", nil, resp)
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
			MessageId: "messag-1",
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
			MessageId: "message-1",
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
			MessageId: "message-1",
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
