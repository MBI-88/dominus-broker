package input

import (
	"dominus-project/app/domain/repos"

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

func (r *restContext) Queries() map[string]string {
	parameters := make(map[string]string)
	args := r.context.QueryArgs()
	args.VisitAll(func(key, value []byte) {
		parameters[string(key)] = string(value)
	})
	return parameters
}

func (r *restContext) Params(key string) string {
	return r.context.UserValue(key).(string)
}

func (r *restContext) QueryInt(key string) int {
	args := r.context.QueryArgs()
	result, err := args.GetUint(key)
	if err != nil {
		return 0
	}
	return result
}

func NewRestContext(ctx *fasthttp.RequestCtx) repos.RestContextInt {
	return &restContext{
		context: ctx,
		js:      jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}
