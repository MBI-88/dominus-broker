package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type apiMiddleware struct {
	token []byte
}

func (a *apiMiddleware) CheckMiddleware(ctx *fasthttp.RequestCtx) error {
	path := string(ctx.RequestURI())
	if strings.HasPrefix(path, "/swagger") {
		return nil 
	}
	token := ctx.Request.Header.Peek("x-api-key")
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