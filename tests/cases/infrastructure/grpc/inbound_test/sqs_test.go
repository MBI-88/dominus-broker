package inbound_test

import (
	"context"
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/infrastructure/enum"
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
	conn, err := grpc.NewClient("bufnet", opts...)
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
			WriteLog(gomock.All(), enum.DEBUG, "Producer", enum.REQUEST_OK)

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
			WriteLog(gomock.All(), enum.DEBUG, "Producer", enum.REQUEST_OK)
		evnetMock.EXPECT().
			WriteLog(gomock.All(), enum.ERROR, "Producer", enum.INVALID_PAYLOAD)

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
			Return(fmt.Errorf("connection error"))

		evnetMock.EXPECT().
			WriteLog(gomock.All(), enum.DEBUG, "Producer", enum.REQUEST_OK)
		evnetMock.EXPECT().
			WriteLog(gomock.All(), enum.ERROR, "Producer", "connection error")

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
				Message: []byte("test-done"),
				MeesageId: "123456789",
				CreatedAt: time.Now(),
			}, nil)

		evnetMock.EXPECT().
			WriteLog(gomock.All(), enum.DEBUG, "Consumer", enum.DEBUG_DESCRIPTION)

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
			WriteLog(gomock.All(), enum.DEBUG, "Consumer", enum.DEBUG_DESCRIPTION)
		evnetMock.EXPECT().
			WriteLog(gomock.All(), enum.ERROR, "Consumer", enum.GROUP_ID)

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
			WriteLog(gomock.All(), enum.DEBUG, "Consumer", enum.DEBUG_DESCRIPTION)
		evnetMock.EXPECT().
			WriteLog(gomock.All(), enum.ERROR, "Consumer", enum.WORKER_ID)

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
			WriteLog(gomock.All(), enum.DEBUG, "Consumer", enum.DEBUG_DESCRIPTION)
		evnetMock.EXPECT().
			WriteLog(gomock.All(), enum.ERROR, "Consumer", enum.WORKER_ID)
		
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
			WriteLog(gomock.All(), enum.DEBUG, "Ack", enum.DEBUG_DESCRIPTION)
		
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
			WorkerId: "consumer-1",
			MessageId: "message-1",
			GroupId: "consumer-g",
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
			WriteLog(gomock.All(), enum.DEBUG, "Ack", enum.DEBUG_DESCRIPTION)
		
		evnetMock.EXPECT(). 
		WriteLog(gomock.All(), enum.ERROR, "Ack", enum.INVALID_ID)
		

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId: "consumer-g",
			WorkerId: "consumer-1",
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
			WriteLog(gomock.All(), enum.DEBUG, "Ack", enum.DEBUG_DESCRIPTION)
		
		evnetMock.EXPECT(). 
		WriteLog(gomock.All(), enum.ERROR, "Ack", enum.GROUP_ID)
		

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId: "",
			WorkerId: "consumer-1",
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
			WriteLog(gomock.All(), enum.DEBUG, "Ack", enum.DEBUG_DESCRIPTION)
		
		evnetMock.EXPECT(). 
		WriteLog(gomock.All(), enum.ERROR, "Ack", enum.WORKER_ID)
		

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewSqsAPI(server, sqsMock, evnetMock)
		client := newSqsAPITestClient(t, lis)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		_, err := client.Ack(context.Background(), &pb.ConsumerRequest{
			GroupId: "consumer-g",
			WorkerId: "",
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
			WriteLog(gomock.All(), enum.DEBUG, "Ack", enum.DEBUG_DESCRIPTION)
		
		evnetMock.EXPECT(). 
		WriteLog(gomock.All(), enum.ERROR, "Ack", enum.WORKER_ID)

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
			GroupId: "consumer-g",
			WorkerId: "consumer-1",
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