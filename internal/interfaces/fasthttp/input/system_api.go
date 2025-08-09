package input

import (
	st "dominus-project/internal/interactors/system"
	"dominus-project/internal/interfaces/fasthttp/errors"
	"fmt"

	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"

	"github.com/valyala/fasthttp"
)

type system struct {
	router *router.Router
	uc     st.SystemService
	js     jsoniter.API
}

// @Tags Logs
// @Description <h3>gets all logs fron the database using page and size (required)</h3>
// @Param page query int true "page"
// @Param size query int true "size"
// @Security ApiKeyAuth
// @Success 200 {object} map[string][]string "Success response"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /logs [get]
func (r *system) getLogs(ctx *fasthttp.RequestCtx) {
	result, err := r.uc.GetLogs()
	if err != nil {
		b := errors.SetErrorResponse(ctx, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		ctx.Response.SetBody(b)
		return
	}
	msg := make(map[string]any, 1)
	msg["logs"] = result
	body, _ := r.js.Marshal(msg)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

// @Tags Logs
// @Description <h3>Gets all selected items from the database for making a backup</h3>
// @Security ApiKeyAuth
// @Success 200 {object} map[string][]string "Success response"
// @Failure 406 {object} map[string]string "Error response"
// @Router /logs-backup [get]
func (r *system) getBackup(ctx *fasthttp.RequestCtx) {
	result, err := r.uc.GetLogs()
	if err != nil {
		b := errors.SetErrorResponse(ctx, "application/json", fasthttp.StatusNotAcceptable, err.Error())
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

func (r *system) path() {
	r.router.GET("/logs", r.getLogs)
	r.router.GET("/logs-backup", r.getBackup)
}

func NewSystemAPI(r *router.Router, uc st.SystemService) {
	re := &system{router: r, uc: uc, js: jsoniter.ConfigCompatibleWithStandardLibrary}
	re.path()
}
