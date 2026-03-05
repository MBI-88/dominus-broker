package middlewares

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/enum"
	"fmt"
	"net"

	"github.com/valyala/fasthttp"
)

type hostAllowed struct {
	cidr string
	log adapters.Logs
}

func NewMiddlewareHost(c string, log adapters.Logs) Middlewares {
	return &hostAllowed{
		cidr: c,
		log: log,
	}
}

func (h *hostAllowed) CheckMiddleware(ctx *fasthttp.RequestCtx) error {
	_, allowNet, err := net.ParseCIDR(h.cidr)
	if err != nil && !allowNet.Contains(ctx.RemoteIP()) {
		go h.log.WriteLog(ctx, enum.ERROR, "CheckMiddleware", "host not allowed")
		return fmt.Errorf("host not allowed")
	}
	return nil
}
