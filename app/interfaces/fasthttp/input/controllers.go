package input

import (
	"dominus/app/interactors"

	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type rest struct {
	router *router.Router
	inter  interactors.InteractorInt
	js     jsoniter.API
}



func (r *rest) getLogs(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	context := NewRestContext(ctx)
	
	result, err := inter.GetLogs(context)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		msg := make(map[string]string)
		msg["message"] = err.Error()
		b, _ := r.js.Marshal(msg)
		ctx.Response.SetBody(b)
		return
	}

	msg := make(map[string]any)
	msg["logs"] = result
	body, err := r.js.Marshal(msg)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		msg := make(map[string]string)
		msg["message"] = err.Error()
		b, _ := r.js.Marshal(msg)
		ctx.Response.SetBody(b)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(body)
}

func (r *rest) getPages(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	context := NewRestContext(ctx)

	result, err := inter.GetTotalPages(context)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		msg := make(map[string]string)
		msg["message"] = err.Error()
		b, _ := r.js.Marshal(msg)
		ctx.Response.SetBody(b)
		return
	}

	msg := make(map[string]any)
	msg["total"] = result
	body, err := r.js.Marshal(msg)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		msg := make(map[string]string)
		msg["message"] = err.Error()
		b, _ := r.js.Marshal(msg)
		ctx.Response.SetBody(b)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(body)
}

func (r *rest) deleteAll(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	context := NewRestContext(ctx)

	if err := inter.DelectLogs(context); err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		msg := make(map[string]string)
		msg["message"] = err.Error()
		b, _ := r.js.Marshal(msg)
		ctx.Response.SetBody(b)
		return
	}

	msg := make(map[string]string)
	msg["message"] = "Operation successful!"
	body, err := r.js.Marshal(msg)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		msg := make(map[string]string)
		msg["message"] = err.Error()
		b, _ := r.js.Marshal(msg)
		ctx.Response.SetBody(b)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(body)
}

func (r *rest) path() {
	r.router.GET("/manager", r.getLogs)
	r.router.GET("/manager-pages", r.getPages)
	r.router.DELETE("/manager", r.deleteAll)
}


func NewRestApi(r *router.Router, i interactors.InteractorInt) {
	re := &rest{router: r, inter: i, js: jsoniter.ConfigCompatibleWithStandardLibrary}
	re.path()
}
