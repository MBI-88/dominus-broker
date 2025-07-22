package middlewares

import (
	"github.com/valyala/fasthttp"
)

type IMiddleware interface {
	Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler
	AddMiddleware(mid ...IMiddlewares)
}

type IMiddlewares interface {
	CheckMiddleware(ctx *fasthttp.RequestCtx) error
}

type middleware struct {
	mids []IMiddlewares
}

func NewMiddleware() IMiddleware {
	return &middleware{
		mids: []IMiddlewares{},
	}
}

func (m *middleware) AddMiddleware(mid ...IMiddlewares) {
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
