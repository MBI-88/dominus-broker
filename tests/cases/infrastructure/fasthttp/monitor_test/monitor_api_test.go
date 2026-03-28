package monitor_test

import (
	"dominus-project/internal/infrastructure/fasthttp/inbound"
	"dominus-project/mocks"
	"testing"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

func TestMonitorController(t *testing.T) {
	t.Run("Metricts_Ok", func(t *testing.T) {
		router := router.New()
		reg := prometheus.NewRegistry()
		crtl := gomock.NewController(t)
		eventMock := mocks.NewMockEvent(crtl)
		inbound.NewMonitorAPI(router, reg, eventMock)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/metrics")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		ctx.Request.Header.Set("X-Metrics-Probe", "1")
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
	})

	t.Run("Health", func(t *testing.T) {
		router := router.New()
		reg := prometheus.NewRegistry()
		crtl := gomock.NewController(t)
		eventMock := mocks.NewMockEvent(crtl)
		inbound.NewMonitorAPI(router, reg, eventMock)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/health")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
	})
}
