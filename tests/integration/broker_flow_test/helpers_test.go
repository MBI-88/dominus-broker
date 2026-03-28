package broker_flow_test

import (
	"context"
	"io"
	"net"
	"testing"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)


const buffSize = 1024 * 1024

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.Dial()
	}
}
// bidirEchoPeer simula un broker remoto: eco por cada StreamRequest recibido.
type bidirEchoPeer struct {
	pb.UnimplementedBrokerAPIServer
}

func (e *bidirEchoPeer) BidirectionalStream(stream pb.BrokerAPI_BidirectionalStreamServer) error {
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

func startTCPGRPCServer(t *testing.T, register func(*grpc.Server)) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	register(srv)
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("tcp peer serve: %v", err)
		}
	}()
	t.Cleanup(func() {
		srv.Stop()
		_ = lis.Close()
	})
	return lis.Addr().String()
}
