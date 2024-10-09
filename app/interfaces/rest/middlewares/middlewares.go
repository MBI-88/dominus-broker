package middlewares

import (
	"github.com/valyala/fasthttp"
)


type middleware struct {
	mids  []middlewaresInt
}

func (m middleware) AddMiddleware(mid ...middlewaresInt) {
	m.mids = append(m.mids, mid...)
}

func (m middleware) Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
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
	//Middleware wrapper for any kind of middlewares that implement middlewaresInt interface
	//
	//middlewareInt has an unic method CheckMiddleware
	//
	//Parameters
	//
	//-> handler: fasthttp Request handler
	//
	//Returns
	//
	//-> error
	Middlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler
	//AddMiddleware adds every instances that implement middlwaresInt interface
	//
	//Parameters
	//
	//-> mi: instance of middlewaresInt
	AddMiddleware(mid ...middlewaresInt)
}

type middlewaresInt interface {
	//CheckMiddleware method for checkin any kind of condicion. It's a method declaration 
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	CheckMiddleware(ctx *fasthttp.RequestCtx) error
}

func NewMiddleware() middlewareInt {
	return &middleware{
		mids: []middlewaresInt{},
	}
}
