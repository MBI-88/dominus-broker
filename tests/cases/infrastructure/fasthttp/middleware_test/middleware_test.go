package middleware_test

import (
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/fasthttp/middlewares"
	"dominus-broker/mocks"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

type alwaysForbidden struct{}

func (alwaysForbidden) CheckMiddleware(*fasthttp.RequestCtx) error {
	return fmt.Errorf("forbidden")
}

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

	t.Run("disallowed when IP not in CIDR", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		mid := middlewares.NewMiddlewareHost("10.0.0.0/8", lgMock)
		ctx := newFastHTTPCtx("/api/publish", "test-token", "127.0.0.1")
		done := make(chan struct{})

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.ERROR, "CheckMiddleware", enum.NO_HOST_ALLOW).
			Do(func(any, string, string, string) { close(done) })

		err := mid.CheckMiddleware(ctx)
		if err == nil {
			t.Fatal("expected error for IP outside CIDR")
		}

		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("expected error log")
		}
	})

	t.Run("invalid CIDR", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		mid := middlewares.NewMiddlewareHost("not-a-valid-cidr", lgMock)
		ctx := newFastHTTPCtx("/api/publish", "test-token", "127.0.0.1")
		done := make(chan struct{})

		lgMock.EXPECT().
			WriteLog(gomock.Any(), enum.ERROR, "CheckMiddleware", enum.NO_HOST_ALLOW).
			Do(func(any, string, string, string) { close(done) })

		err := mid.CheckMiddleware(ctx)
		if err == nil {
			t.Fatal("expected error for invalid CIDR")
		}

		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("expected error log")
		}
	})
}

func TestMiddlewareChain(t *testing.T) {
	t.Run("forbidden short-circuits before handler", func(t *testing.T) {
		mw := middlewares.NewMiddleware()
		mw.AddMiddleware(alwaysForbidden{})
		handlerCalled := false
		h := mw.Middlewares(func(ctx *fasthttp.RequestCtx) {
			handlerCalled = true
		})
		ctx := newFastHTTPCtx("/api/x", "", "127.0.0.1")
		h(ctx)
		if handlerCalled {
			t.Fatal("handler should not run when middleware fails")
		}
		if ctx.Response.StatusCode() != fasthttp.StatusForbidden {
			t.Fatalf("status: got %d want %d", ctx.Response.StatusCode(), fasthttp.StatusForbidden)
		}
		if string(ctx.Response.Header.ContentType()) != enum.CONTENT_TYPE_TEXT {
			t.Fatalf("content-type: got %q", ctx.Response.Header.ContentType())
		}
	})

	t.Run("ok runs handler after passing middlewares", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lgMock := mocks.NewMockEvent(ctrl)
		apiMid := middlewares.NewMiddlewareApiToken("ok-secret", lgMock)
		ctx := newFastHTTPCtx("/api/x", "ok-secret", "127.0.0.1")

		lgMock.EXPECT().CheckID(gomock.Any()).Return(ctx)

		mw := middlewares.NewMiddleware()
		mw.AddMiddleware(apiMid)
		handlerCalled := false
		h := mw.Middlewares(func(c *fasthttp.RequestCtx) {
			handlerCalled = true
			if c != ctx {
				t.Fatal("expected same ctx")
			}
		})
		h(ctx)
		if !handlerCalled {
			t.Fatal("handler should run")
		}
	})
}
