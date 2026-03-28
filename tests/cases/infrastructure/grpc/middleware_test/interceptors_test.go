package middleware_test

import (
	"context"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/grpc/middlewares"
	"dominus-project/mocks"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryAuthInterceptor(t *testing.T) {
	t.Run("adds api token and calls invoker", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		it := middlewares.NewInterceptor("test-token", lgMock)
		done := make(chan struct{})

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.DEBUG, "UnaryAuthInterceptor", enum.DEBUG_DESCRIPTION).
			Do(func(context.Context, string, string, string) {
				close(done)
			})

		invokerCalled := false
		invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			invokerCalled = true
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				t.Fatalf("expected outgoing metadata")
			}
			got := md.Get(enum.X_API_KEY)
			if len(got) == 0 || got[0] != "test-token" {
				t.Fatalf("expected %s=test-token, got %v", enum.X_API_KEY, got)
			}
			return nil
		}

		err := it.UnaryAuthInterceptor(context.Background(), "/dominus.Api/Publish", "req", "reply", nil, invoker)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !invokerCalled {
			t.Fatalf("expected invoker to be called")
		}

		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("expected debug log call")
		}
	})

}

func TestStreamAuthInterceptor(t *testing.T) {
	t.Run("success adds token and calls streamer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		it := middlewares.NewInterceptor("test-token", lgMock)
		done := make(chan struct{})

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.DEBUG, "StreamAuthInterceptor", enum.DEBUG_DESCRIPTION).
			Do(func(context.Context, string, string, string) {
				close(done)
			})

		streamerCalled := false
		streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			streamerCalled = true
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				t.Fatalf("expected outgoing metadata")
			}
			got := md.Get(enum.X_API_KEY)
			if len(got) == 0 || got[0] != "test-token" {
				t.Fatalf("expected %s=test-token, got %v", enum.X_API_KEY, got)
			}
			return nil, nil
		}

		s, err := it.StreamAuthInterceptor(context.Background(), &grpc.StreamDesc{}, nil, "/dominus.Api/PublishStream", streamer)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if s != nil {
			t.Fatalf("expected nil stream in this test")
		}
		if !streamerCalled {
			t.Fatalf("expected streamer to be called")
		}

		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("expected debug log call")
		}
	})

	t.Run("error from streamer is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		it := middlewares.NewInterceptor("test-token", lgMock)
		errBoom := errors.New("boom")
		debugDone := make(chan struct{})
		errorDone := make(chan struct{})

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.DEBUG, "StreamAuthInterceptor", enum.DEBUG_DESCRIPTION).
			Do(func(context.Context, string, string, string) {
				close(debugDone)
			})
		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.ERROR, "StreamAuthInterceptor", "boom").
			Do(func(context.Context, string, string, string) {
				close(errorDone)
			})

		streamer := func(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error) {
			return nil, errBoom
		}

		s, err := it.StreamAuthInterceptor(context.Background(), &grpc.StreamDesc{}, nil, "/dominus.Api/PublishStream", streamer)
		if err == nil {
			t.Fatalf("expected error")
		}
		if err.Error() != "boom" {
			t.Fatalf("expected boom error, got %v", err)
		}
		if s != nil {
			t.Fatalf("expected nil stream on error")
		}

		select {
		case <-debugDone:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("expected debug log call")
		}
		select {
		case <-errorDone:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("expected error log call")
		}
	})
}
