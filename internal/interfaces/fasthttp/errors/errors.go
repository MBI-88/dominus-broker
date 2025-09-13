package errors

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var (
	js = jsoniter.ConfigCompatibleWithStandardLibrary
)

func SetErrorResponse(ctx *fasthttp.RequestCtx, contenType string, statusCode int, er string) []byte {
	ctx.Response.Header.Set("Content-Type", contenType)
	ctx.Response.Header.SetStatusCode(statusCode)
	msg := make(map[string]string)
	msg["message"] = er
	b, _ := js.Marshal(msg)
	return b
}
