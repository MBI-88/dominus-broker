package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"net"

	"github.com/valyala/fasthttp"
)

type Middlewares func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx

type middleware struct {
	ApiToken string
	Cidr     string
}

func (m *middleware) SetApiToken(token, cidr string) {
	m.ApiToken = token
	m.Cidr = cidr
}

func (m middleware) apiMiddleware(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	token := ctx.Request.Header.Peek("API_TOKEN")
	hashedToken := sha256.Sum256(token)
	hashedKey := sha256.Sum256([]byte(m.ApiToken))

	if subtle.ConstantTimeCompare(hashedKey[:], hashedToken[:]) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return nil
	}
	return ctx

}

func (m middleware) allowedHosts(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	_, allowNet, err := net.ParseCIDR(m.Cidr)
	if err != nil && !allowNet.Contains(ctx.RemoteIP()) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return nil
	}
	return ctx

}

func (m middleware) Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if m.ApiToken != "" {
			if ctx = m.apiMiddleware(ctx); ctx == nil {
				return
			}
		}
		if m.Cidr != "" {
			if ctx = m.allowedHosts(ctx); ctx == nil {
				return
			}
		}	
		handler(ctx)
	}
}


type middlewareInt interface {
	SetApiToken(token, cidr string)
	Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler
}

func NewMiddleware() middlewareInt {
	return new(middleware)
}
