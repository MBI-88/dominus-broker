package input

import (
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
}

func (m *monitor) convertToHTTP(ctx *fasthttp.RequestCtx) (*http.Request, error) {
	req := &http.Request{
		Method: string(ctx.Method()),
		Header: make(http.Header),
	}
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		req.Header.Set(string(k), string(v))
	})
	return req, nil
}

func (m *monitor) getMetrics(ctx *fasthttp.RequestCtx) {
	handler := promhttp.HandlerFor(m.reg, m.opts)
	resp := &adapter{ctx}
	req, err := m.convertToHTTP(ctx)
	if err != nil {
		ctx.Error("Error converting request", fasthttp.StatusInternalServerError)
		return
	}
	m.collectMetrics()
	handler.ServeHTTP(resp, req)
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

func (m *monitor) path() {
	m.router.GET("/metrics", m.getMetrics)
}

func NewMonitor(r *router.Router, reg *prometheus.Registry) {
	reg.MustRegister(cpumetrics, memorymetrics)
	m := &monitor{
		router: r,
		reg:    reg,
		opts:   promhttp.HandlerOpts{EnableOpenMetrics: true, DisableCompression: true},
	}
	m.path()
}
