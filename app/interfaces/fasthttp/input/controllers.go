package input

import (
	"dominus/app/interactors"
	"fmt"

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
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	result, err := manger.GetLogs(context)
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusExpectationFailed, err.Error())
		ctx.Response.SetBody(b)
		return
	}

	msg := make(map[string]any)
	msg["logs"] = result
	body, _ := r.js.Marshal(msg)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

func (r *rest) getPages(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	result, err := manger.GetTotalPages(context)
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusExpectationFailed, err.Error())
		ctx.Response.SetBody(b)
		return
	}

	msg := make(map[string]any)
	msg["total"] = result
	body, _ := r.js.Marshal(msg)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

func (r *rest) deleteAll(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	if err := manger.DelectLogs(context); err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusExpectationFailed, err.Error())
		ctx.Response.SetBody(b)
		return
	}

	msg := make(map[string]string)
	msg["message"] = "Operation successful!"
	body, _ := r.js.Marshal(msg)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(body)
}

func (r *rest) getStats(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	result, err := manger.GetStats()
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusExpectationFailed, err.Error())
		ctx.Response.SetBody(b)
		return
	}

	body, _ := r.js.Marshal(result)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

func (r *rest) getBackup(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	result, err := manger.GetBackup(context)
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusExpectationFailed, err.Error())
		ctx.Response.SetBody(b)
		return
	}

	msg := make(map[string]any)
	msg["logs"] = result
	body, _ := r.js.Marshal(msg)
	ctx.Response.Header.Set("Content-Type", "application/octet-stream")
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=logs.json")
	ctx.Response.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

func (r *rest) path() {
	r.router.GET("/manager", r.getLogs)
	r.router.GET("/manager-pages", r.getPages)
	r.router.DELETE("/manager", r.deleteAll)
	r.router.GET("/manager-stats", r.getStats)
	r.router.GET("/manager-backup", r.getBackup)
}

func (r *rest) setErrorResponse(ctx *fasthttp.RequestCtx,  contenType string, statusCode int, er string) []byte {
	ctx.Response.Header.Set("Content-Type", contenType)
	ctx.Response.Header.SetStatusCode(statusCode)
	msg := make(map[string]string)
	msg["message"] = er
	b, _ := r.js.Marshal(msg)
	return b
}

func NewRestController(r *router.Router, i interactors.InteractorInt) {
	re := &rest{router: r, inter: i, js: jsoniter.ConfigCompatibleWithStandardLibrary}
	re.path()
}
