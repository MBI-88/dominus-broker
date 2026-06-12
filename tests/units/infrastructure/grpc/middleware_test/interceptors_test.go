package middleware_test

import (
	"context"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/grpc/middlewares"
	"dominus-broker/mocks"
	"errors"
	"sync"
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
			WriteLog(gomock.Any(), enum.DEBUG, "interceptor.UnaryAuthInterceptor", enum.DEBUG_DESCRIPTION).
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
			WriteLog(gomock.Any(), enum.DEBUG, "interceptor.StreamAuthInterceptor", enum.DEBUG_DESCRIPTION).
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
		var mu sync.Mutex
		seen := make(map[string]int)

		lgMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), "interceptor.StreamAuthInterceptor", gomock.Any()).
			Times(2).
			Do(func(_ context.Context, level, op, dsc string) {
				mu.Lock()
				seen[level+"|"+dsc]++
				mu.Unlock()
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

		deadline := time.After(200 * time.Millisecond)
		for {
			mu.Lock()
			ok := seen[enum.DEBUG+"|"+enum.DEBUG_DESCRIPTION] == 1 && seen[enum.ERROR+"|boom"] == 1
			mu.Unlock()
			if ok {
				break
			}
			select {
			case <-deadline:
				mu.Lock()
				got := seen
				mu.Unlock()
				t.Fatalf("expected debug+error StreamAuthInterceptor logs, got %#v", got)
			case <-time.After(5 * time.Millisecond):
			}
		}
	})
}
