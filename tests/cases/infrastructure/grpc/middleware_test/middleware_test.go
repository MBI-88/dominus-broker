package middleware_test

import (
	"context"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/grpc/middlewares"
	"dominus-project/mocks"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type testServerStream struct {
	ctx context.Context
}

func (t *testServerStream) SetHeader(metadata.MD) error  { return nil }
func (t *testServerStream) SendHeader(metadata.MD) error { return nil }
func (t *testServerStream) SetTrailer(metadata.MD)       {}
func (t *testServerStream) Context() context.Context     { return t.ctx }
func (t *testServerStream) SendMsg(any) error            { return nil }
func (t *testServerStream) RecvMsg(any) error            { return nil }

func TestApiToken(t *testing.T) {

	t.Run("ApiToken ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)
		baseCtx := context.Background()
		ctx := metadata.NewIncomingContext(baseCtx, metadata.Pairs(enum.X_API_KEY, "test-token"))

		lgMock.EXPECT().
			CheckID(ctx).
			Return(ctx)

		_, err := mid.ApiToken(ctx)
		if err != nil {
			t.Fatalf("Expected error nil got %s\n", err)
		}

	})

	t.Run("ApiToken empty context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)
		ctx := context.Background()

		lgMock.EXPECT().
			CheckID(ctx).
			Return(ctx)

		lgMock.EXPECT().
			WriteLog(ctx, gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		_, err := mid.ApiToken(ctx)
		if err == nil {
			t.Fatalf("Expected error got %v\n", err)
		}

	})

	t.Run("ApiToken different token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)
		baseCtx := context.Background()
		ctx := metadata.NewIncomingContext(baseCtx, metadata.Pairs(enum.X_API_KEY, "test-T"))

		lgMock.EXPECT().
			CheckID(ctx).
			Return(ctx)

		lgMock.EXPECT().
			WriteLog(ctx, gomock.All(), gomock.All(), gomock.All()).
			AnyTimes()

		_, err := mid.ApiToken(ctx)
		if err == nil {
			t.Fatalf("Expected error nil got %v\n", err)
		}

	})

}

func TestUnaryLog(t *testing.T) {

	t.Run("UnaryLog ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{FullMethod: "/dominus.Api/Publish"}
		req := "request-body"
		handlerCalled := false

		// UnaryLog writes logs in a goroutine; allow it without making the test timing-sensitive.
		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.DEBUG, info.FullMethod, enum.DEBUG_DESCRIPTION).
			AnyTimes()

		handler := func(hCtx context.Context, hReq any) (any, error) {
			handlerCalled = true
			if hCtx != ctx {
				t.Fatalf("Expected same ctx")
			}
			if hReq != req {
				t.Fatalf("Expected req %v got %v", req, hReq)
			}
			return "ok", nil
		}

		res, err := mid.UnaryLog(ctx, req, info, handler)
		if err != nil {
			t.Fatalf("Expected error nil got %v", err)
		}
		if res != "ok" {
			t.Fatalf("Expected response ok got %v", res)
		}
		if !handlerCalled {
			t.Fatalf("Expected handler to be called")
		}

	})

}

func TestStreamLog(t *testing.T) {
	t.Run("StreamLog ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)
		ctx := context.Background()
		ss := &testServerStream{ctx: ctx}
		info := &grpc.StreamServerInfo{FullMethod: "/dominus.Api/PublishStream"}
		handlerCalled := false

		// StreamLog writes logs in a goroutine; allow async invocation.
		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.DEBUG, info.FullMethod, enum.DEBUG_DESCRIPTION).
			AnyTimes()

		handler := func(srv any, stream grpc.ServerStream) error {
			handlerCalled = true
			if srv != "srv" {
				t.Fatalf("Expected srv got %v", srv)
			}
			if stream != ss {
				t.Fatalf("Expected same stream instance")
			}
			return nil
		}

		err := mid.StreamLog("srv", ss, info, handler)
		if err != nil {
			t.Fatalf("Expected nil error got %v", err)
		}
		if !handlerCalled {
			t.Fatalf("Expected handler to be called")
		}
	})
}

func TestLogErrors(t *testing.T) {
	t.Run("LogErrors maps levels and writes logs", func(t *testing.T) {
		tests := []struct {
			name          string
			level         logging.Level
			expectedLevel string
		}{
			{name: "debug", level: logging.LevelDebug, expectedLevel: enum.DEBUG},
			{name: "info", level: logging.LevelInfo, expectedLevel: enum.INFO},
			{name: "warn", level: logging.LevelWarn, expectedLevel: enum.WARN},
			{name: "error default", level: logging.LevelError, expectedLevel: enum.ERROR},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				lgMock := mocks.NewMockEvent(ctrl)
				checkerMock := mocks.NewMockCheckerClient(ctrl)
				mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)
				logger := mid.LogErrors()
				done := make(chan struct{})

				lgMock.EXPECT().
					WriteLog(gomock.Any(), tc.expectedLevel, "LogErrors", "test-message").
					Do(func(context.Context, string, string, string) {
						close(done)
					})

				logger.Log(context.Background(), tc.level, "test-message")

				select {
				case <-done:
				case <-time.After(200 * time.Millisecond):
					t.Fatalf("WriteLog was not called for level %v", tc.level)
				}
			})
		}
	})
}
