package inbound_test

import (
	"strings"
	"testing"

	"dominus-project/internal/infrastructure/fasthttp/inbound"
	"dominus-project/mocks"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

// Smoke: GET /metrics goes through monitor.getMetrics → promhttp.HandlerFor → inbound.adapter
// (Write / WriteHeader on *fasthttp.RequestCtx). Same style as grpc outbound_test against public API only.
func TestMonitorAPI_MetricsPrometheusAdapterSmoke(t *testing.T) {
	r := router.New()
	reg := prometheus.NewRegistry()
	ctrl := gomock.NewController(t)
	eventMock := mocks.NewMockEvent(ctrl)
	inbound.NewMonitorAPI(r, reg, eventMock)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/metrics")
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.Header.Set("X-Adapter-Smoke", "1")
	r.Handler(ctx)

	if got := ctx.Response.StatusCode(); got != fasthttp.StatusOK {
		t.Fatalf("status: got %d want %d", got, fasthttp.StatusOK)
	}
	body := string(ctx.Response.Body())
	if len(body) == 0 {
		t.Fatal("metrics body is empty")
	}
	// Gauges registered in NewMonitorAPI (collectMetrics fills them on scrape).
	if !strings.Contains(body, "cpu_usage") && !strings.Contains(body, "memory_usage") {
		t.Fatalf("expected monitor gauge names in prometheus output, got prefix: %.200q", body)
	}
}
