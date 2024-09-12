package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"

	"github.com/valyala/fasthttp"
)

type middleware struct {
	ApiToken string
}

func (m *middleware) SetApiToken(k string) {
	m.ApiToken = k
}

func (m middleware) APImiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token := ctx.Request.Header.Peek("API_TOKEN")
		hashedToken := sha256.Sum256(token)
		hashedKey := sha256.Sum256([]byte(m.ApiToken))

		if subtle.ConstantTimeCompare(hashedKey[:], hashedToken[:]) == 0 {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		next(ctx)
	}

}

type middlewareInt interface {
	SetApiToken(k string)
	APImiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler 
}


func NewMiddleware() middlewareInt {
	return new(middleware)
}