package broker_flow_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"dominus-project/internal/application/use_cases/broker"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/grpc/inbound"
	"dominus-project/internal/infrastructure/grpc/middlewares"
	"dominus-project/internal/infrastructure/grpc/outbound"
	"dominus-project/mocks"

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

func TestBidirectionalStreamFlow(t *testing.T) {
	t.Run("BidirectionalStream flow", func(t *testing.T) {
		echo := &bidirEchoPeer{}
		addr1 := startTCPGRPCServer(t, func(s *grpc.Server) {
			pb.RegisterBrokerAPIServer(s, echo)
		})
		addr2 := startTCPGRPCServer(t, func(s *grpc.Server) {
			pb.RegisterBrokerAPIServer(s, echo)
		})
		subscribers := []string{addr1, addr2}

		mainLis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		eventMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		eventMock.EXPECT().
			CheckID(gomock.Any()).
			DoAndReturn(func(ctx context.Context) context.Context { return ctx }).
			AnyTimes()
		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		const apiKey = "integration-token"
		mid := middlewares.NewMiddleware(apiKey, eventMock, checkerMock)
		inter := middlewares.NewInterceptor(apiKey, eventMock)

		// Dial real peers (extremo final); sin bufconn.
		subDialOpts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithStreamInterceptor(inter.StreamAuthInterceptor),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
			grpc.WithDefaultServiceConfig(enum.ROUTER_CLIENT),
		}
		outClient := outbound.NewGrpClient(subDialOpts, eventMock)
		brApp := broker.NewBroker(outClient)

		server := grpc.NewServer(grpc.ChainStreamInterceptor(
			auth.StreamServerInterceptor(mid.ApiToken),
		))
		inbound.NewBrokerAPI(server, brApp, eventMock)
		go func() {
			if err := server.Serve(mainLis); err != nil {
				t.Logf("main Serve: %v", err)
			}
		}()
		t.Cleanup(server.Stop)

		ingressDial := []grpc.DialOption{
			grpc.WithContextDialer(bufDialer(mainLis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithStreamInterceptor(inter.StreamAuthInterceptor),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
			grpc.WithDefaultServiceConfig(enum.ROUTER_CLIENT),
		}

		rootCtx, rootCancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer rootCancel()
		ctx, cancel := context.WithCancel(rootCtx)
		defer cancel()

		conn, err := grpc.NewClient("passthrough:///buf", ingressDial...)
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })

		cli := pb.NewBrokerAPIClient(conn)
		stream, err := cli.BidirectionalStream(ctx)
		if err != nil {
			t.Fatalf("BidirectionalStream: %v", err)
		}

		if err := stream.Send(&pb.StreamRequestMessage{
			Subscribers: subscribers,
			Payload:     []byte("hello"),
		}); err != nil {
			t.Fatalf("Send first: %v", err)
		}
		if err := stream.Send(&pb.StreamRequestMessage{Payload: []byte("world")}); err != nil {
			t.Fatalf("Send second: %v", err)
		}

		// 2 suscriptores × 2 mensajes provider = 4 respuestas (orden no garantizado).
		want := map[string]int{
			"echo:hello": 2,
			"echo:world": 2,
		}
		got := map[string]int{}
		for i := 0; i < 4; i++ {
			resp, err := stream.Recv()
			if err != nil {
				t.Fatalf("Recv %d: %v", i, err)
			}
			got[string(resp.GetPayload())]++
		}
		for k, w := range want {
			if got[k] != w {
				t.Fatalf("payload counts: got %+v want %+v", got, want)
			}
		}

		// Los suscriptores TCP (echo) no cierran el stream solos: si solo hiciéramos CloseSend(),
		// streamProv se cerraría pero los workers del outbound seguirían bloqueados en Recv hacia
		// el echo y nunca llegaría el tx a `closed` (deadlock con <-closed y close(closed) en el use case).
		// Cancelar el ctx del RPC desbloquea BidirectionalStream (workers + Wait + señal closed).
		cancel()

		_, err = stream.Recv()
		if err == nil {
			t.Fatal("expected terminal error after cancel")
		}
		if errors.Is(err, io.EOF) {
			return
		}
		switch status.Code(err) {
		case codes.Canceled:
		case codes.Aborted:
			// p. ej. "connection closed" si el handler termina antes que el ctx cancele el stream
		case codes.Unavailable:
		default:
			t.Fatalf("final Recv: want Canceled/Aborted/Unavailable/EOF, got %v: %v", status.Code(err), err)
		}
	})
}
