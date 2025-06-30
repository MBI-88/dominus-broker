package middlewares

import (
	"github.com/valyala/fasthttp"
)


type middleware struct {
	mids  []middlewaresInt
}

func (m *middleware) AddMiddleware(mid ...middlewaresInt) {
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

type middlewareInt interface {
	Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler
	AddMiddleware(mid ...middlewaresInt)
}

type middlewaresInt interface {
	CheckMiddleware(ctx *fasthttp.RequestCtx) error
}

func NewMiddleware() middlewareInt {
	return &middleware{
		mids: []middlewaresInt{},
	}
}
