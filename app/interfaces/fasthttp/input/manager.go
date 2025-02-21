package input

import (
	"dominus/app/interactors"
	"fmt"

	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	
	"github.com/valyala/fasthttp"
)

type manager struct {
	router *router.Router
	service  interactors.ManagerInt
	js     jsoniter.API
}

 

// @Tags Manager
// @Description <h3>gets all logs fron the database using page and size (required)</h3>
// @Param page query int true "page"
// @Param size query int true "size"
// @Security ApiKeyAuth
// @Success 200 {object} []entities.Logs "Success response"
// @Failure 404 {object} map[string]string "Response body {message: error}"
// @Router /manager [get]
func (r *manager) getLogs(ctx *fasthttp.RequestCtx) {
	context := NewRestContext(ctx)
	result, err := r.service.GetLogs(context)
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusNotFound, err.Error())
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


// @Tags Manager
// @Description <h3>Gets total pages in the database</h3>
// @Security ApiKeyAuth
// @Success 200 {object} map[string]int "Success response: {pages:int}"
// @Failure 404 {object} map[string]string "Response body {message: error}"
// @Router /manager-pages [get]
func (r *manager) getPages(ctx *fasthttp.RequestCtx) {
	context := NewRestContext(ctx)
	result, err := r.service.GetTotalPages(context)
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusNotFound, err.Error())
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



// @Tags Manager
// @Description <h3>Delete all logs in the database</h3>
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Respose response: {message:string}"
// @Failure 404 {object} map[string]string "Error response {message: error}"
// @Router /manager [delete]
func (r *manager) deleteAll(ctx *fasthttp.RequestCtx) {
	context := NewRestContext(ctx)
	if err := r.service.DelectLogs(context); err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusNotFound, err.Error())
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


// @Tags Manager
// @Description <h3>Gets all statistic from the database</h3>
// @Security ApiKeyAuth
// @Success 200 {object} map[string]any "Success response: { avgObjSize: int,collections: int, dataSize: int, db: string, fsTotalSize: int, fsUsedSize: int, indexSize: int, indexes: int, objects: int,ok: int, scaleFactor: int, storageSize: int, totalSize: int, views: int}"
// @Failure 404 {object} map[string]string "Error response {message: error}"
// @Router /manager-stats [get]
func (r *manager) getStats(ctx *fasthttp.RequestCtx) {
	result, err := r.service.GetStats()
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusNotFound, err.Error())
		ctx.Response.SetBody(b)
		return
	}

	body, _ := r.js.Marshal(result)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}


// @Tags Manager
// @Description <h3>Gets all selected items from the database for making a backup</h3>
// @Security ApiKeyAuth
// @Success 200 {object} []entities.Logs "Success response"
// @Failure 404 {object} map[string]string "Error response"
// @Router /manager-backup [get]
func (r *manager) getBackup(ctx *fasthttp.RequestCtx) {
	context := NewRestContext(ctx)
	result, err := r.service.GetBackup(context)
	if err != nil {
		b := r.setErrorResponse(ctx, "application/json", fasthttp.StatusNotFound, err.Error())
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



func (r *manager) path() {
	r.router.GET("/manager", r.getLogs)
	r.router.GET("/manager-pages", r.getPages)
	r.router.DELETE("/manager", r.deleteAll)
	r.router.GET("/manager-stats", r.getStats)
	r.router.GET("/manager-backup", r.getBackup)
}

func (r *manager) setErrorResponse(ctx *fasthttp.RequestCtx,  contenType string, statusCode int, er string) []byte {
	ctx.Response.Header.Set("Content-Type", contenType)
	ctx.Response.Header.SetStatusCode(statusCode)
	msg := make(map[string]string)
	msg["message"] = er
	b, _ := r.js.Marshal(msg)
	return b
}

func NewManager(r *router.Router, i interactors.ManagerInt) {
	re := &manager{router: r, service: i, js: jsoniter.ConfigCompatibleWithStandardLibrary}
	re.path()
}
