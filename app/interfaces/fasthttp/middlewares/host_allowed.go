package middlewares

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type hostAllowed struct {
	origins string
}

func (h *hostAllowed) CheckMiddleware(ctx *fasthttp.RequestCtx) error {
	origin := string(ctx.Request.Header.Peek("Origin"))
	if h.origins != "*" { 
		if  origin != "" && strings.Contains(h.origins, origin) {
			return nil 
		}
		return fmt.Errorf("Host not allowed")
	}
	return nil 
}


func NewMiddlewareHost(c string) middlewaresInt {
	var og string
	if c == "" {
		og = "*"
	}else {
		og = c
	}
	return &hostAllowed{
		origins: og,
	}
}