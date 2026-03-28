package outbound_test

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"dominus-project/internal/infrastructure/grpc/outbound"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type stubEvent struct{}

func (stubEvent) WriteLog(context.Context, string, string, string) {}

func (stubEvent) CheckID(ctx context.Context) context.Context { return ctx }

type clientStreamSubscriber struct {
	pb.UnimplementedBrokerAPIServer
}

func (clientStreamSubscriber) ClientStream(stream pb.BrokerAPI_ClientStreamServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponseMessage{Status: int64(codes.OK)})
		}
		if err != nil {
			return err
		}
	}
}

func TestBrokerClient_ClientStream_smoke(t *testing.T) {
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	pb.RegisterBrokerAPIServer(srv, clientStreamSubscriber{})
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("Serve: %v", err)
		}
	}()
	t.Cleanup(srv.Stop)

	opts := []grpc.DialOption{
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	client := outbound.NewGrpClient(opts, stubEvent{})
	msgCh := make(chan []byte, 2)
	msgCh <- []byte("a")
	msgCh <- []byte("b")
	close(msgCh)

	done := make(chan struct{})
	go func() {
		client.ClientStream([]string{"passthrough:///peer-a", "passthrough:///peer-b"}, msgCh, context.Background())
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("ClientStream did not finish")
	}
}
