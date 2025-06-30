package input

import (
	"dominus-project/internal/domain/repos"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type restContext struct {
	context *fasthttp.RequestCtx
	js      jsoniter.API
}

func (r *restContext) BodyParser(obj any) error {
	body := r.context.Request.Body()
	if err := r.js.Unmarshal(body, obj); err != nil {
		return err
	}
	return nil
}

func (r *restContext) Param(key string) string {
	return r.context.UserValue(key).(string)
}

func NewRestContext(ctx *fasthttp.RequestCtx) repos.RestContextInt {
	return &restContext{
		context: ctx,
		js:      jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}
