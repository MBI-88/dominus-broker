package input

import (
	"dominus/app/interactors"
	"mime/multipart"

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

func (r *restContext) FormFile(key string) (*multipart.FileHeader, error) {
	return r.context.FormFile(key)
}

func (r *restContext) FormValue(key string) []byte {
	return r.context.FormValue(key)
}

func NewRestContext(ctx *fasthttp.RequestCtx) interactors.RestContextInt {
	return &restContext{
		context: ctx,
		js:      jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}
