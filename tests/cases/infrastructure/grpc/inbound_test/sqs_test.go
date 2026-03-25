package inbound_test

import (
	"context"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/grpc/inbound"
	"dominus-project/mocks"
	"testing"

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

func TestSqsAPI(t *testing.T) {

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
}
