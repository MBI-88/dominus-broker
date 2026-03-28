package broker_flow_test

import (
	"context"
	"testing"
	"time"

	"dominus-project/internal/application/use_cases/broker"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/grpc/inbound"
	"dominus-project/internal/infrastructure/grpc/middlewares"
	"dominus-project/mocks"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/test/bufconn"
)

func TestClientStreamFlow(t *testing.T) {
	t.Run("ClientStream flow", func(t *testing.T) {
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

		mockBrokerClient.EXPECT().
			ClientStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_, msg, _ any) {
				ch := msg.(<-chan []byte)
				for range ch {
				}
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
		stream, err := cli.ClientStream(ctx)
		if err != nil {
			t.Fatalf("ClientStream: %v", err)
		}

		subscribers := []string{"subscriber-a.example", "subscriber-b.example"}
		if err := stream.Send(&pb.StreamRequestMessage{
			Subscribers: subscribers,
			Payload:     []byte("first"),
		}); err != nil {
			t.Fatalf("Send first: %v", err)
		}
		for i := 0; i < 4; i++ {
			if err := stream.Send(&pb.StreamRequestMessage{Payload: []byte("testing")}); err != nil {
				t.Fatalf("Send: %v", err)
			}
		}

		resp, err := stream.CloseAndRecv()
		if err != nil {
			t.Fatalf("CloseAndRecv: %v", err)
		}
		if got, want := resp.GetStatus(), int64(codes.OK); got != want {
			t.Fatalf("response status: got %d want %d", got, want)
		}
	})
}
