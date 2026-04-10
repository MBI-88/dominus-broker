package broker_flow_test

import (
	"context"
	"testing"
	"time"

	"dominus-broker/internal/application/usecases/broker"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/grpc/inbound"
	"dominus-broker/internal/infrastructure/grpc/middlewares"
	"dominus-broker/mocks"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestServerStreamFlow(t *testing.T) {
	t.Run("ServerStream flow", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		eventMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mockBrokerClient := mocks.NewMockBrokerClient(ctrl)

		eventMock.EXPECT().
			CheckID(gomock.Any()).
			DoAndReturn(func(ctx context.Context) context.Context { return ctx }).
			AnyTimes()
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		subscribers := []string{"subscriber-a.example", "subscriber-b.example"}
		initialPayload := []byte("seed-from-client")

		mockBrokerClient.EXPECT().
			ServerStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(urlsArg, initArg, msgArg, ctxArg, txArg any) {
				urls := urlsArg.([]string)
				init := initArg.([]byte)
				if string(init) != string(initialPayload) {
					t.Errorf("initial payload: got %q want %q", init, initialPayload)
				}
				if len(urls) != len(subscribers) {
					t.Errorf("subscribers: got %v want %v", urls, subscribers)
				}
				msg := msgArg.(chan<- []byte)
				_ = ctxArg.(context.Context)
				tx := txArg.(chan<- struct{})

				msg <- []byte("chunk-1")
				msg <- []byte("chunk-2")
				// stream is buffered; StreamServerConn select can pick <-tx before draining msg unless we yield.
				time.Sleep(50 * time.Millisecond)
				tx <- struct{}{}
			}).
			AnyTimes()

		const apiKey = "integration-token"
		mid := middlewares.NewMiddleware(apiKey, eventMock, checkerMock)
		inter := middlewares.NewInterceptor(apiKey, eventMock)

		dialOpts := []grpc.DialOption{
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithStreamInterceptor(inter.StreamAuthInterceptor),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
			grpc.WithDefaultServiceConfig(enum.ROUTER_CLIENT),
		}

		serverOpts := []grpc.ServerOption{
			grpc.ChainStreamInterceptor(
				auth.StreamServerInterceptor(mid.ApiToken),
			),
		}

		server := grpc.NewServer(serverOpts...)
		brApp := broker.NewBroker(mockBrokerClient)
		inbound.NewBrokerAPI(server, brApp, eventMock)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Logf("Serve: %v", err)
			}
		}()
		t.Cleanup(server.Stop)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := grpc.NewClient("passthrough:///buf", dialOpts...)
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })

		cli := pb.NewBrokerAPIClient(conn)
		stream, err := cli.ServerStream(ctx, &pb.StreamRequestMessage{
			Subscribers: subscribers,
			Payload:     initialPayload,
		})
		if err != nil {
			t.Fatalf("ServerStream: %v", err)
		}

		var payloads [][]byte
		for {
			resp, err := stream.Recv()
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("Recv final error: %v", err)
				}
				if st.Code() != codes.Aborted {
					t.Fatalf("want code Aborted, got %v: %v", st.Code(), err)
				}
				if st.Message() != "broker.StreamServerConn connection closed" {
					t.Fatalf("want message %q, got %q", "broker.StreamServerConn connection closed", st.Message())
				}
				break
			}
			payloads = append(payloads, append([]byte(nil), resp.GetPayload()...))
		}

		if len(payloads) != 2 {
			t.Fatalf("want 2 stream payloads, got %d: %q", len(payloads), payloads)
		}
		if string(payloads[0]) != "chunk-1" || string(payloads[1]) != "chunk-2" {
			t.Fatalf("unexpected payloads: %q", payloads)
		}
	})
}
