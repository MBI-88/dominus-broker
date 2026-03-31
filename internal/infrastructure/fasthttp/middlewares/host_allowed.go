package middlewares

import (
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"
	"fmt"
	"net"

	"github.com/valyala/fasthttp"
)

type hostAllowed struct {
	cidr string
	log  event.Event
}

func NewMiddlewareHost(c string, log event.Event) Middlewares {
	return &hostAllowed{
		cidr: c,
		log:  log,
	}
}

func (h *hostAllowed) CheckMiddleware(ctx *fasthttp.RequestCtx) error {
	_, allowNet, err := net.ParseCIDR(h.cidr)
	if err != nil {
		go h.log.WriteLog(ctx, enum.ERROR, "CheckMiddleware", enum.NO_HOST_ALLOW)
		return fmt.Errorf(enum.NO_HOST_ALLOW)
	}
	if !allowNet.Contains(ctx.RemoteIP()) {
		go h.log.WriteLog(ctx, enum.ERROR, "CheckMiddleware", enum.NO_HOST_ALLOW)
		return fmt.Errorf(enum.NO_HOST_ALLOW)
	}
	return nil
}
