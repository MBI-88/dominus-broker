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
		go h.log.WriteLog(ctx, enum.ERROR, "CheckMiddleware", enum.NO_HOST_ALLOW)
		return fmt.Errorf(enum.NO_HOST_ALLOW)
	}
	return nil
}
