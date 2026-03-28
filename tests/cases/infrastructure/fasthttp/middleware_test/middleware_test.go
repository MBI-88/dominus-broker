package middleware_test

import (
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/fasthttp/middlewares"
	"dominus-project/mocks"
	"net"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

func newFastHTTPCtx(uri, token, ip string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	req := &fasthttp.Request{}
	req.SetRequestURI(uri)
	if token != "" {
		req.Header.Set(enum.X_API_KEY, token)
	}
	ctx.Init(req, &net.TCPAddr{IP: net.ParseIP(ip), Port: 12345}, nil)
	return ctx
}

func TestCheckMiddleware(t *testing.T) {
	t.Run("api token ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		mid := middlewares.NewMiddlewareApiToken("test-token", lgMock)
		ctx := newFastHTTPCtx("/api/publish", "test-token", "127.0.0.1")

		lgMock.EXPECT().
			CheckID(gomock.Any()).
			Return(ctx)

		err := mid.CheckMiddleware(ctx)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}
	})

	t.Run("api token invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		mid := middlewares.NewMiddlewareApiToken("test-token", lgMock)
		ctx := newFastHTTPCtx("/api/publish", "bad-token", "127.0.0.1")
		done := make(chan struct{})

		lgMock.EXPECT().
			CheckID(gomock.Any()).
			Return(ctx)
		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.ERROR, "CheckMiddleware", enum.MATCH_TOKEN).
			Do(func(any, string, string, string) {
				close(done)
			})

		err := mid.CheckMiddleware(ctx)
		if err == nil {
			t.Fatalf("expected error")
		}

		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("expected error log call")
		}
	})

	t.Run("swagger bypass", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		mid := middlewares.NewMiddlewareApiToken("test-token", lgMock)
		ctx := newFastHTTPCtx("/swagger/index.html", "", "127.0.0.1")

		lgMock.EXPECT().
			CheckID(gomock.Any()).
			Return(ctx)

		err := mid.CheckMiddleware(ctx)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}
	})

}

func TestAllowedHost(t *testing.T) {
	t.Run("allowed host ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		mid := middlewares.NewMiddlewareHost("127.0.0.0/8", lgMock)
		ctx := newFastHTTPCtx("/api/publish", "test-token", "127.0.0.1")

		err := mid.CheckMiddleware(ctx)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}
	})
}
