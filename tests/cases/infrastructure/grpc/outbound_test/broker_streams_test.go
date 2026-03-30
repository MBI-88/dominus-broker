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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type serverStreamStub struct {
	pb.UnimplementedBrokerAPIServer
}

func (serverStreamStub) ServerStream(_ *pb.StreamRequestMessage, stream pb.BrokerAPI_ServerStreamServer) error {
	for i := 0; i < 2; i++ {
		if err := stream.Send(&pb.StreamResponseMessage{Payload: []byte{byte('0' + i)}}); err != nil {
			return err
		}
	}
	return nil
}

type bidirEchoStub struct {
	pb.UnimplementedBrokerAPIServer
}

func (bidirEchoStub) BidirectionalStream(stream pb.BrokerAPI_BidirectionalStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		out := append([]byte("echo:"), req.GetPayload()...)
		if err := stream.Send(&pb.StreamResponseMessage{Payload: out}); err != nil {
			return err
		}
	}
}

func TestBrokerClient_ServerStream_smoke(t *testing.T) {
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	pb.RegisterBrokerAPIServer(srv, serverStreamStub{})
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
	msg := make(chan []byte, 4)
	tx := make(chan struct{}, 1)

	go func() {
		for range msg {
		}
	}()

	done := make(chan struct{})
	go func() {
		client.ServerStream(
			[]string{"passthrough:///peer"},
			[]byte("init"),
			msg,
			context.Background(),
			tx,
		)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("ServerStream did not finish")
	}

	select {
	case <-tx:
	default:
		t.Fatal("expected tx signal after workers complete")
	}
}

func TestBrokerClient_ServerStream_nilTx(t *testing.T) {
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	pb.RegisterBrokerAPIServer(srv, serverStreamStub{})
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
	msg := make(chan []byte, 4)
	go func() {
		for range msg {
		}
	}()

	done := make(chan struct{})
	go func() {
		client.ServerStream(
			[]string{"passthrough:///peer"},
			[]byte("init"),
			msg,
			context.Background(),
			nil,
		)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("ServerStream (nil tx) did not finish")
	}
}

func TestBrokerClient_BidirectionalStream_smoke(t *testing.T) {
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	pb.RegisterBrokerAPIServer(srv, bidirEchoStub{})
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
	prov := make(chan []byte, 2)
	sub := make(chan []byte, 8)
	tx := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runDone := make(chan struct{})
	go func() {
		client.BidirectionalStream(
			[]string{"passthrough:///peer"},
			prov,
			sub,
			tx,
			ctx,
		)
		close(runDone)
	}()

	go func() {
		prov <- []byte("ping")
		close(prov)
	}()

	select {
	case p := <-sub:
		if string(p) != "echo:ping" {
			t.Fatalf("subscriber payload: got %q want %q", p, "echo:ping")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for echoed payload on sub")
	}

	// Workers block on a second Recv while the echo server waits for another client message; cancel unblocks them.
	cancel()

	select {
	case <-tx:
	case <-time.After(15 * time.Second):
		t.Fatal("BidirectionalStream did not signal tx after cancel")
	}

	select {
	case <-runDone:
	case <-time.After(10 * time.Second):
		t.Fatal("BidirectionalStream did not return after cancel")
	}
}
