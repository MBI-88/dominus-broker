package inbound

import (
	"dominus-project/internal/infrastructure/event"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/valyala/fasthttp"
)

var (
	cpumetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_usage_percentage",
		Help: "Current cpu usage in percentage",
	}, []string{"cpu_used"})
	memorymetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "memory_usage_percentage",
		Help: "Current memory usage in percentage",
	}, []string{"mem_used"})
)

type adapter struct {
	ctx *fasthttp.RequestCtx
}

func (a *adapter) Header() http.Header {
	return http.Header{}
}

func (a *adapter) Write(data []byte) (int, error) {
	return a.ctx.Write(data)
}

func (a *adapter) WriteHeader(status int) {
	a.ctx.SetStatusCode(status)
}

type monitor struct {
	router *router.Router
	opts   promhttp.HandlerOpts
	reg    *prometheus.Registry
	log    event.Event
}

func NewMonitorAPI(r *router.Router, reg *prometheus.Registry, log event.Event) {
	reg.MustRegister(cpumetrics, memorymetrics)
	m := &monitor{
		router: r,
		reg:    reg,
		opts:   promhttp.HandlerOpts{EnableOpenMetrics: true, DisableCompression: true},
		log:  log,
	}
	m.path()
}

func (m *monitor) convertToHTTP(ctx *fasthttp.RequestCtx) (*http.Request, error) {
	req := &http.Request{
		Method: string(ctx.Method()),
		Header: make(http.Header),
	}
	for k,v := range ctx.Request.Header.All() {
		req.Header.Set(string(k), string(v))
	}
	return req, nil
}

func (m *monitor) collectMetrics() {
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil {
		cpumetrics.WithLabelValues("cpu_used").Set(cpuPercent[0])
	}
	mem, err := mem.VirtualMemory()
	if err == nil {
		memorymetrics.WithLabelValues("mem_used").Set(mem.UsedPercent)
	}
}

// @Tags Monitor
// @Description <h3>gets metrics</h3>
// @Security ApiKeyAuth
// @Success 200 {object} map[string][]string "Success response"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /metrics [get]
func (m *monitor) getMetrics(ctx *fasthttp.RequestCtx) {
	handler := promhttp.HandlerFor(m.reg, m.opts)
	resp := &adapter{ctx}
	req, err := m.convertToHTTP(ctx)
	if err != nil {
		ctx.Error("Error converting request", fasthttp.StatusNotAcceptable)
		return
	}
	m.collectMetrics()
	handler.ServeHTTP(resp, req)
}

// @Tags Monitor
// @Description <h3>get healthCeck</h3>
// @Security ApiKeyAuth
// @Success 200 {object} map[string][]string "Success response"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /health [get]
func (*monitor) getHealthCheck(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody([]byte("Health ok"))
}

func (m *monitor) path() {
	m.router.GET("/metrics", m.getMetrics)
	m.router.GET("/health", m.getHealthCheck)
}
