package middleware_test

import (
	"context"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/grpc/middlewares"
	"dominus-project/mocks"
	"errors"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func TestIdPotency(t *testing.T) {
	t.Run("ok schedules SaveConsumer for new key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)

		baseCtx := context.Background()
		ctx := metadata.NewIncomingContext(baseCtx, metadata.Pairs(enum.ID_POTENCY_HEADER, "idem-key-1"))

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.DEBUG, "IdPotency", enum.DEBUG_DESCRIPTION).
			AnyTimes()

		checkerMock.EXPECT().
			CheckConsumer(gomock.Any(), "idem-key-1").
			Return(false)

		saveDone := make(chan struct{})
		checkerMock.EXPECT().
			SaveConsumer(gomock.Any(), "idem-key-1").
			DoAndReturn(func(context.Context, string) error {
				close(saveDone)
				return nil
			}).
			Times(1)

		out, err := mid.IdPotency(ctx)
		if err != nil {
			t.Fatalf("IdPotency: %v", err)
		}
		if out != ctx {
			t.Fatalf("expected same context, got different pointer")
		}

		select {
		case <-saveDone:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("SaveConsumer was not invoked")
		}
	})

	t.Run("aborted when key already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(enum.ID_POTENCY_HEADER, "idem-dup"))

		lgMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		checkerMock.EXPECT().
			CheckConsumer(gomock.Any(), "idem-dup").
			Return(true)

		_, err := mid.IdPotency(ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("not a status error: %v", err)
		}
		if st.Code() != codes.Aborted || st.Message() != "id potency found" {
			t.Fatalf("got code=%v msg=%q want Aborted / id potency found", st.Code(), st.Message())
		}
	})

	t.Run("data loss when incoming metadata missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)

		ctx := context.Background()

		lgMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		_, err := mid.IdPotency(ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("not a status error: %v", err)
		}
		if st.Code() != codes.DataLoss || st.Message() != enum.NOT_FOUND {
			t.Fatalf("got code=%v msg=%q want DataLoss / %q", st.Code(), st.Message(), enum.NOT_FOUND)
		}
	})

	t.Run("data loss when idempotency header empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(enum.ID_POTENCY_HEADER, ""))

		lgMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		_, err := mid.IdPotency(ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("not a status error: %v", err)
		}
		if st.Code() != codes.DataLoss || st.Message() != enum.NOT_FOUND {
			t.Fatalf("got code=%v msg=%q want DataLoss / %q", st.Code(), st.Message(), enum.NOT_FOUND)
		}
	})

	t.Run("SaveConsumer error is logged", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		checkerMock := mocks.NewMockCheckerClient(ctrl)
		mid := middlewares.NewMiddleware("test-token", lgMock, checkerMock)

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(enum.ID_POTENCY_HEADER, "idem-fail-save"))

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.DEBUG, "IdPotency", enum.DEBUG_DESCRIPTION).
			AnyTimes()

		checkerMock.EXPECT().
			CheckConsumer(gomock.Any(), "idem-fail-save").
			Return(false)

		saveErr := errors.New("redis down")
		logDone := make(chan struct{})
		checkerMock.EXPECT().
			SaveConsumer(gomock.Any(), "idem-fail-save").
			Return(saveErr)

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.ERROR, "IdPotency.SaveConsumer", saveErr.Error()).
			Do(func(context.Context, string, string, string) { close(logDone) })

		out, err := mid.IdPotency(ctx)
		if err != nil {
			t.Fatalf("IdPotency: %v", err)
		}
		if out != ctx {
			t.Fatalf("expected same context")
		}

		select {
		case <-logDone:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("WriteLog for SaveConsumer error was not called")
		}
	})
}
