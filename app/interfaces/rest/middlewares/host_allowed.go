package middlewares

import (
	"fmt"
	"net"

	"github.com/valyala/fasthttp"
)

type hostAllowed struct {
	cidr string
}

func (h hostAllowed) CheckMiddleware(ctx *fasthttp.RequestCtx) error {
	_, allowNet, err := net.ParseCIDR(h.cidr)
	if err != nil && !allowNet.Contains(ctx.RemoteIP()) {
		return fmt.Errorf("Host not allowed")
	}
	return nil 
}


func NewMiddlewareHot(c string) middlewaresInt {
	return &hostAllowed{
		cidr: c,
	}
}