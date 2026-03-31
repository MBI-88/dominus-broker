package inbound_test

import (
	"context"
	"dominus-broker/internal/application/dtos"
	"dominus-broker/internal/infrastructure/grpc/inbound"
	"dominus-broker/internal/infrastructure/grpc/outbound"
	"dominus-broker/mocks"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const buffSize = 1024 * 1024

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestClientStream(t *testing.T) {

	t.Run("ClientStream Ok", func(t *testing.T) {
		c := make(chan struct{})
		msg := make(chan []byte)
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		brokerMock := mocks.NewMockBroker(ctrl)
		eventMock := mocks.NewMockEvent(ctrl)

		brokerMock.EXPECT().
			StreamClientConn(gomock.All()).
			Return(nil).AnyTimes()

		eventMock.EXPECT().
			WriteLog(context.Background(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewBrokerAPI(server, brokerMock, eventMock)
		client := outbound.NewGrpClient(opts, eventMock)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		go func() {
			for range 5 {
				msg <- []byte("testing")
			}
			close(msg)
		}()

		client.ClientStream([]string{"server1.api.com", "server2.api.com"}, msg, context.Background())

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})

	t.Run("ClientStream error", func(t *testing.T) {
		c := make(chan struct{})
		msg := make(chan []byte)
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		eventMock := mocks.NewMockEvent(ctrl)
		brokerMock := mocks.NewMockBroker(ctrl)

		brokerMock.EXPECT().
			StreamClientConn(gomock.All()).
			Return(fmt.Errorf("error")).AnyTimes()

		eventMock.EXPECT().
			WriteLog(context.Background(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewBrokerAPI(server, brokerMock, eventMock)
		client := outbound.NewGrpClient(opts, eventMock)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		go func() {
			for range 1 {
				msg <- []byte("testing")
			}
			close(msg)
		}()

		client.ClientStream([]string{"server1.api.com", "server2.api.com"}, msg, context.Background())

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})

	t.Run("ClientStream ingress CloseAndRecv on io.EOF", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		brokerMock := mocks.NewMockBroker(ctrl)
		eventMock := mocks.NewMockEvent(ctrl)

		brokerMock.EXPECT().
			StreamClientConn(gomock.Any()).
			DoAndReturn(func(st any) error {
				bc := st.(dtos.BrokerClientDto)
				if _, err := bc.Recv(); err != nil {
					return err
				}
				for {
					if _, err := bc.Recv(); err != nil {
						return err
					}
				}
			}).
			Times(1)

		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		srv := grpc.NewServer()
		inbound.NewBrokerAPI(srv, brokerMock, eventMock)
		go func() {
			if err := srv.Serve(lis); err != nil {
				t.Logf("Serve: %v", err)
			}
		}()
		t.Cleanup(srv.Stop)

		conn, err := grpc.NewClient("passthrough:///buf",
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })

		cli := pb.NewBrokerAPIClient(conn)
		stream, err := cli.ClientStream(context.Background())
		if err != nil {
			t.Fatalf("ClientStream: %v", err)
		}
		if err := stream.Send(&pb.StreamRequestMessage{
			Subscribers: []string{"s1.example"},
			Payload:     []byte("one"),
		}); err != nil {
			t.Fatalf("Send: %v", err)
		}
		if err := stream.CloseSend(); err != nil {
			t.Fatalf("CloseSend: %v", err)
		}
		resp, err := stream.CloseAndRecv()
		if err != nil {
			t.Fatalf("CloseAndRecv: %v", err)
		}
		if got, want := resp.GetStatus(), int64(codes.OK); got != want {
			t.Fatalf("status: got %d want %d", got, want)
		}
	})

	t.Run("ClientStream ingress broker returns error", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		brokerMock := mocks.NewMockBroker(ctrl)
		eventMock := mocks.NewMockEvent(ctrl)

		brokerMock.EXPECT().
			StreamClientConn(gomock.Any()).
			Return(fmt.Errorf("stream failure")).
			Times(1)

		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		srv := grpc.NewServer()
		inbound.NewBrokerAPI(srv, brokerMock, eventMock)
		go func() {
			if err := srv.Serve(lis); err != nil {
				t.Logf("Serve: %v", err)
			}
		}()
		t.Cleanup(srv.Stop)

		conn, err := grpc.NewClient("passthrough:///buf",
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })

		cli := pb.NewBrokerAPIClient(conn)
		stream, err := cli.ClientStream(context.Background())
		if err != nil {
			t.Fatalf("ClientStream: %v", err)
		}
		_, err = stream.CloseAndRecv()
		if err == nil {
			t.Fatal("expected error from CloseAndRecv")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Aborted || st.Message() != "stream failure" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ClientStream ingress broker returns nil", func(t *testing.T) {
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		brokerMock := mocks.NewMockBroker(ctrl)
		eventMock := mocks.NewMockEvent(ctrl)

		brokerMock.EXPECT().
			StreamClientConn(gomock.Any()).
			Return(nil).
			Times(1)

		eventMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		srv := grpc.NewServer()
		inbound.NewBrokerAPI(srv, brokerMock, eventMock)
		go func() {
			if err := srv.Serve(lis); err != nil {
				t.Logf("Serve: %v", err)
			}
		}()
		t.Cleanup(srv.Stop)

		conn, err := grpc.NewClient("passthrough:///buf",
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })

		cli := pb.NewBrokerAPIClient(conn)
		stream, err := cli.ClientStream(context.Background())
		if err != nil {
			t.Fatalf("ClientStream: %v", err)
		}
		resp, err := stream.CloseAndRecv()
		if err != nil {
			t.Fatalf("CloseAndRecv: %v", err)
		}
		if got, want := resp.GetStatus(), int64(codes.OK); got != want {
			t.Fatalf("status: got %d want %d", got, want)
		}
	})
}

func TestServerStream(t *testing.T) {
	t.Run("ServerStream Ok", func(t *testing.T) {
		c := make(chan struct{})
		msg := make(chan []byte)
		done := make(chan struct{})
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		eventMock := mocks.NewMockEvent(ctrl)
		brokerMock := mocks.NewMockBroker(ctrl)

		brokerMock.EXPECT().
			StreamServerConn(gomock.All(), gomock.All()).
			DoAndReturn(func(req, st any) error {
				tem := st.(dtos.BrokerServerDto)
				for range 10 {
					if err := tem.Send([]byte("test-body")); err != nil {
						log.Panicln(err)
					}
				}
				return nil
			}).AnyTimes()

		eventMock.EXPECT().
			WriteLog(context.Background(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewBrokerAPI(server, brokerMock, eventMock)
		client := outbound.NewGrpClient(opts, eventMock)
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		go func() {
			time.Sleep(5 * time.Second)
			cancel()
		}()

		go func() {
			for range msg {
			}
		}()

		go func() {
			<-done
		}()

		client.ServerStream([]string{"server1.api.com", "server2.api.com"}, []byte("test-initial"), msg, ctx, done)

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})

	t.Run("ServerStream error", func(t *testing.T) {
		c := make(chan struct{})
		msg := make(chan []byte)
		done := make(chan struct{})
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		eventMock := mocks.NewMockEvent(ctrl)
		brokerMock := mocks.NewMockBroker(ctrl)

		brokerMock.EXPECT().
			StreamServerConn(gomock.All(), gomock.All()).
			DoAndReturn(func(req, st any) error {
				tem := st.(dtos.BrokerServerDto)
				for range 10 {
					if err := tem.Send([]byte("test-body")); err != nil {
						log.Panicln(err)
					}
				}
				return fmt.Errorf("error")
			}).AnyTimes()

		eventMock.EXPECT().
			WriteLog(context.Background(), gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewBrokerAPI(server, brokerMock, eventMock)
		client := outbound.NewGrpClient(opts, eventMock)
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		go func() {
			time.Sleep(5 * time.Second)
			cancel()
		}()

		go func() {
			for range msg {
			}
		}()

		go func() {
			<-done
		}()

		client.ServerStream([]string{"server1.api.com", "server2.api.com"}, []byte("test-initial"), msg, ctx, done)

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})
}

func TestBidirectionalStream(t *testing.T) {

	t.Run("BidirectionalStream Ok", func(t *testing.T) {
		c := make(chan struct{})
		provMsg := make(chan []byte, 2)
		subMsg := make(chan []byte, 2)
		tx := make(chan struct{})
		done := make(chan struct{}, 8)
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		eventMock := mocks.NewMockEvent(ctrl)
		brokerMock := mocks.NewMockBroker(ctrl)
		ctx, cancel := context.WithCancel(context.Background())

		brokerMock.EXPECT().
			StreamBiConn(gomock.All()).
			DoAndReturn(func(st any) error {
				tem := st.(dtos.BrokerBidirectionalDto)
				go func() {
					for {
						resp, err := tem.Recv()
						if err != nil {
							return
						}
						log.Println(resp)
					}
				}()

				for range 10 {
					if err := tem.Send([]byte("test-body")); err != nil {
						log.Println(err)
					}
				}

				return nil
			}).AnyTimes()

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		eventMock.EXPECT().
			WriteLog(ctx, gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewBrokerAPI(server, brokerMock, eventMock)
		client := outbound.NewGrpClient(opts, eventMock)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		go func() {
			time.Sleep(5 * time.Second)
			cancel()
			done <- struct{}{}
		}()

		go func() {
			for range 10 {
				provMsg <- []byte("test-body")
			}
			close(provMsg)
		}()

		go func() {
			<-tx
			close(tx)
		}()

		go func() {
			for range subMsg {
			}
		}()

		client.BidirectionalStream([]string{"server1.api.com", "server2.api.com"},
			provMsg,
			subMsg,
			tx,
			ctx,
		)

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})

	t.Run("BidirectionalStream error", func(t *testing.T) {
		c := make(chan struct{})
		provMsg := make(chan []byte, 2)
		subMsg := make(chan []byte, 2)
		tx := make(chan struct{})
		done := make(chan struct{}, 8)
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		eventMock := mocks.NewMockEvent(ctrl)
		brokerMock := mocks.NewMockBroker(ctrl)
		ctx, cancel := context.WithCancel(context.Background())

		brokerMock.EXPECT().
			StreamBiConn(gomock.All()).
			DoAndReturn(func(st any) error {
				tem := st.(dtos.BrokerBidirectionalDto)
				go func() {
					for {
						resp, err := tem.Recv()
						if err != nil {
							return
						}
						log.Println(resp)
					}
				}()

				for range 10 {
					if err := tem.Send([]byte("test-body")); err != nil {
						log.Println(err)
					}
				}

				return fmt.Errorf("error")
			}).AnyTimes()

		eventMock.EXPECT().
			WriteLog(ctx, gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := grpc.NewServer([]grpc.ServerOption{}...)
		inbound.NewBrokerAPI(server, brokerMock, eventMock)
		client := outbound.NewGrpClient(opts, eventMock)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		go func() {
			time.Sleep(5 * time.Second)
			cancel()
			done <- struct{}{}
		}()

		go func() {
			for range 10 {
				provMsg <- []byte("test-body")
			}
			close(provMsg)
		}()

		go func() {
			<-tx
			close(tx)
		}()

		go func() {
			for range subMsg {
			}
		}()

		client.BidirectionalStream([]string{"server1.api.com", "server2.api.com"},
			provMsg,
			subMsg,
			tx,
			ctx,
		)

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})
}
