package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"github.com/valyala/fasthttp"
)

type apiMiddleware struct {
	token []byte
}

func (a *apiMiddleware) CheckMiddleware(ctx *fasthttp.RequestCtx) error {
	token := ctx.Request.Header.Peek("API_TOKEN")
	hashedToken := sha256.Sum256(token)
	hashedKey := sha256.Sum256(a.token)
	if subtle.ConstantTimeCompare(hashedKey[:], hashedToken[:]) == 0 {
		return fmt.Errorf("Invalid token")
	}
	return nil
}


func NewMiddlewareApiToken(t string) middlewaresInt {
	return &apiMiddleware{
		token: []byte(t),
	}
}