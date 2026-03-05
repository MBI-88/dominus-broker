package middlewares

import (
	"github.com/valyala/fasthttp"
)

type Middleware interface {
	Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler
	AddMiddleware(mid ...Middlewares)
}

type Middlewares interface {
	CheckMiddleware(ctx *fasthttp.RequestCtx) error
}

type middleware struct {
	mids []Middlewares
}

func NewMiddleware() Middleware {
	return &middleware{
		mids: []Middlewares{},
	}
}

func (m *middleware) AddMiddleware(mid ...Middlewares) {
	m.mids = append(m.mids, mid...)
}

func (m *middleware) Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		for _, mid := range m.mids {
			if err := mid.CheckMiddleware(ctx); err != nil {
				ctx.Response.Header.SetStatusCode(fasthttp.StatusForbidden)
				ctx.Response.Header.SetContentType("application/text")
				ctx.Response.SetBodyString(err.Error())
				return
			}
		}
		handler(ctx)
	}
}
