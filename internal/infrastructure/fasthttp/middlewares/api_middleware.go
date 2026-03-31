package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"
	"fmt"

	"strings"

	"github.com/valyala/fasthttp"
)

type apMiddleware struct {
	token []byte
	log   event.Event
}

func NewMiddlewareApiToken(t string, log event.Event) Middlewares {
	return &apMiddleware{
		token: []byte(t),
		log:   log,
	}
}

func (a *apMiddleware) CheckMiddleware(ctx *fasthttp.RequestCtx) error {
	if t, ok := a.log.CheckID(ctx).(*fasthttp.RequestCtx); ok {
		ctx = t
	}
	path := string(ctx.RequestURI())
	if strings.HasPrefix(path, "/swagger") {
		return nil
	}
	token := ctx.Request.Header.Peek(enum.X_API_KEY)
	hashedToken := sha256.Sum256(token)
	hashedKey := sha256.Sum256(a.token)
	if subtle.ConstantTimeCompare(hashedKey[:], hashedToken[:]) == 0 {
		go a.log.WriteLog(ctx, enum.ERROR, "CheckMiddleware", enum.MATCH_TOKEN)
		return fmt.Errorf(enum.MATCH_TOKEN)
	}
	return nil
}
