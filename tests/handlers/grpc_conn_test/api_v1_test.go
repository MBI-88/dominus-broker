package grpcconn_test

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/infrastructure/grpc/input"
	"dominus-project/internal/infrastructure/grpc/output"
	"dominus-project/mocks"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

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

func TestSimple(t *testing.T) {
	t.Run("Simple Ok", func(t *testing.T) {
		c := make(chan struct{})
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			SimpleConn(gomock.All()).
			Return(nil).Times(1)

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		resp, err := client.Simple("127.0.0.1:5050", []byte("test-simple"))

		if err == nil && resp.GetStatus() != 0 {
			t.Fatalf("Expected %d received %d", 0, resp.GetStatus())
		}

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})

	t.Run("Simple error", func(t *testing.T) {
		c := make(chan struct{})
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			SimpleConn(gomock.All()).
			Return(fmt.Errorf("error")).Times(1)

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Log(err)
			}
		}()

		resp, err := client.Simple("127.0.0.1:5050", []byte("test-simple"))

		if err == nil && resp.GetStatus() != 0 {
			t.Fatalf("Expected %d received %d", 0, resp.GetStatus())
		} else if err.Error() != status.Error(codes.InvalidArgument, "error").Error() {
			t.Fatalf("Received %s", err)
		}

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})

}

func TestClientStream(t *testing.T) {

	t.Run("ClientStream Ok", func(t *testing.T) {
		c := make(chan struct{})
		msg := make(chan []byte)
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			StreamClientConn(gomock.All()).
			Return(nil).AnyTimes()

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)

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
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			StreamClientConn(gomock.All()).
			Return(fmt.Errorf("error")).AnyTimes()

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)

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
}

func TestServerStream(t *testing.T) {
	t.Run("ServerStream Ok", func(t *testing.T) {
		c := make(chan struct{})
		msg := make(chan []byte)
		done := make(chan struct{})
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			StreamServerConn(gomock.All(), gomock.All()).
			DoAndReturn(func(req, st any) error {
				tem := st.(adapters.StreamServer)
				for range 10 {
					if err := tem.Send([]byte("test-body")); err != nil {
						log.Panicln(err)
					}
				}
				return nil
			}).AnyTimes()

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)
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
			for _, ok := <-msg; ok; {
			}
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
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			StreamServerConn(gomock.All(), gomock.All()).
			DoAndReturn(func(req, st any) error {
				tem := st.(adapters.StreamServer)
				for range 10 {
					if err := tem.Send([]byte("test-body")); err != nil {
						log.Panicln(err)
					}
				}
				return fmt.Errorf("error")
			}).AnyTimes()

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)
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
			for _, ok := <-msg; ok; {
			}
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
		errMsg := make(chan error, 2)
		tx := make(chan struct{})
		done := make(chan struct{})
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			StreamBiConn(gomock.All()).
			DoAndReturn(func(st any) error {
				tem := st.(adapters.StreamBi)
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

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)
		ctx, cancel := context.WithCancel(context.Background())

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
			for _, ok := <-subMsg; ok; {
			}
			close(subMsg)
		}()

		client.BidirectionalStream([]string{"server1.api.com", "server2.api.com"},
			provMsg,
			subMsg,
			errMsg,
			tx,
			ctx,
			done,
		)

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})

	t.Run("BidirectionalStream error", func(t *testing.T) {
		c := make(chan struct{})
		provMsg := make(chan []byte, 2)
		subMsg := make(chan []byte, 2)
		errMsg := make(chan error, 2)
		tx := make(chan struct{})
		done := make(chan struct{})
		lis := bufconn.Listen(buffSize)
		ctrl := gomock.NewController(t)
		logMock := mocks.NewMockLogs(ctrl)
		serviceMock := mocks.NewMockGrpcService(ctrl)

		serviceMock.EXPECT().
			StreamBiConn(gomock.All()).
			DoAndReturn(func(st any) error {
				tem := st.(adapters.StreamBi)
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

		serviceMock.EXPECT().
			RunQueue(c).
			Return(nil).Times(1)

		opts := append([]grpc.DialOption{},
			grpc.WithContextDialer(bufDialer(lis)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		server := input.NewGrpcAPI([]grpc.ServerOption{}, serviceMock, logMock)
		client := output.NewGrpClient(opts, logMock)
		ctx, cancel := context.WithCancel(context.Background())

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
			for _, ok := <-subMsg; ok; {
			}
			close(subMsg)
		}()

		client.BidirectionalStream([]string{"server1.api.com", "server2.api.com"},
			provMsg,
			subMsg,
			errMsg,
			tx,
			ctx,
			done,
		)

		close(c)
		server.GracefulStop()
		ctrl.Finish()
	})
}
