package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net"

	"github.com/valyala/fasthttp"
)

type Middlewares func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx

type middleware struct {
	token    []byte
	cidr     string
}


func (m middleware) apiMiddleware(ctx *fasthttp.RequestCtx) error {
	token := ctx.Request.Header.Peek("API_TOKEN")
	hashedToken := sha256.Sum256(token)
	hashedKey := sha256.Sum256(m.token)
	if subtle.ConstantTimeCompare(hashedKey[:], hashedToken[:]) == 0 {
		return fmt.Errorf("Invalid token")
	}
	return nil 
}

func (m middleware) allowedHosts(ctx *fasthttp.RequestCtx) error {
	_, allowNet, err := net.ParseCIDR(m.cidr)
	if err != nil && !allowNet.Contains(ctx.RemoteIP()) {
		return fmt.Errorf("Host not allowed")
	}
	return nil

}

func (m middleware) Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if m.token != nil {
			if err := m.apiMiddleware(ctx); err != nil {
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
				ctx.Response.SetBody([]byte(err.Error()))
				return
			}
		}
		if m.cidr != "" {
			if err := m.allowedHosts(ctx); err != nil {
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
				ctx.Response.SetBody([]byte(err.Error()))
				return
			}
		}	
		handler(ctx)
	}
}


type middlewareInt interface {
	Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler
}

func NewMiddleware(token, cidr string) middlewareInt {
	return &middleware{
		token: []byte(token),
		cidr: cidr,
	}
}
