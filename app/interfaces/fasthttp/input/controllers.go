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

// getLogs 
// @Summary gets logs
// @Tag Manager
// @Description gets all logs fron the database using page and size (required)
// @Param page query int true "page"
// @Param size query int true "size"
// @Param API_TOKEN header string true "api token to connect with Dominus"
// @securityDefinitions.apikey API_TOKEN
// @in header
// @name API_TOKEN
// @Success 200 {object} entities.Logs "Respose body {id:string,desc:string,create_at:time, stage:string, subscriber:string}"
// @Failure 404 {object} map[string]string "Response body {message: error}"
// @Router /manager [get]
func (r *rest) getLogs(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	result, err := manger.GetLogs(context)
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

// getPages
// @Summary gets pages
// @Tag Manager
// @Description gets total pages in the database
// @Param API_TOKEN header string true "api token to connect with Dominus"
// @securityDefinitions.apikey API_TOKEN
// @in header
// @name API_TOKEN
// @Success 200 {object} map[string]int "Respose body {pages:int}"
// @Failure 404 {object} map[string]string "Response body {message: error}"
// @Router /manager-pages [get]
func (r *rest) getPages(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	result, err := manger.GetTotalPages(context)
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

// deleteAll
// @Summary delete logs
// @Tag Manager
// @Description delete all logs in the database
// @Param API_TOKEN header string true "api token to connect with Dominus"
// @securityDefinitions.apikey API_TOKEN
// @in header
// @name API_TOKEN
// @Success 200 {object} map[string]string "Respose body {message:string}"
// @Failure 404 {object} map[string]string "Response body {message: error}"
// @Router /manager [delete]
func (r *rest) deleteAll(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	if err := manger.DelectLogs(context); err != nil {
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

// getStats
// @Summary gets stats
// @Tag Manager
// @Description gets all statistic from the database
// @Param API_TOKEN header string true "api token to connect with Dominus"
// @securityDefinitions.apikey API_TOKEN
// @in header
// @name API_TOKEN
// @Success 200 {object} map[string]any "Respose body { avgObjSize: int,collections: int, dataSize: int, db: string, fsTotalSize: int, fsUsedSize: int, indexSize: int, indexes: int, objects: int,ok: int, scaleFactor: int, storageSize: int, totalSize: int, views: int}"
// @Failure 404 {object} map[string]string "Response body {message: error}"
// @Router /manager-stats [get]
func (r *rest) getStats(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	result, err := manger.GetStats()
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

// getBackup
// @Summary gets backup
// @Tag Manager
// @Description gets all selected items from the database for making a backup
// @Param API_TOKEN header string true "api token to connect with Dominus"
// @securityDefinitions.apikey API_TOKEN
// @in header
// @name API_TOKEN
// @Success 200 {object} map[string]any "Respose body {logs:[{id:string,desc:string,create_at:time, stage:string, subscriber:string}]}"
// @Failure 404 {object} map[string]string "Response body {message: error}"
// @Router /manager-backup [get]
func (r *rest) getBackup(ctx *fasthttp.RequestCtx) {
	manger := r.inter.NewManager()
	context := NewRestContext(ctx)

	result, err := manger.GetBackup(context)
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
