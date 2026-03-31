package sqs_flow_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"dominus-broker/internal/domain/repositories"
	"dominus-broker/internal/infrastructure/event"
	"dominus-broker/internal/infrastructure/redis/cmemory"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const buffSize = 1024 * 1024

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.Dial()
	}
}

func newMemoryClient(t *testing.T, s *miniredis.Miniredis, streamID string, lg event.Event) repositories.MemoryClient {
	t.Helper()
	host, portStr, err := net.SplitHostPort(s.Addr())
	if err != nil {
		t.Fatalf("SplitHostPort: %v", err)
	}
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		t.Fatalf("ParseInt port: %v", err)
	}
	return cmemory.NewMemoryClient(
		10,
		0,
		int(time.Second),
		int(time.Second),
		int(time.Second),
		port,
		0,
		host,
		"",
		false,
		"",
		streamID,
		lg,
	)
}

func testRedisClient(t *testing.T, s *miniredis.Miniredis) *redis.Client {
	t.Helper()
	return redis.NewClient(&redis.Options{
		Addr:         s.Addr(),
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
}

func newSqsAPIClient(t *testing.T, lis *bufconn.Listener) pb.SqsAPIClient {
	t.Helper()
	opts := []grpc.DialOption{
		grpc.WithContextDialer(bufDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("passthrough:///bufnet", opts...)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return pb.NewSqsAPIClient(conn)
}
