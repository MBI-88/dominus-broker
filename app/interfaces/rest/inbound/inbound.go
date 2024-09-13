package inbound

import (
	"dominus/app/interactors"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type rest struct {
	router *router.Router
}

func (r rest) local(ctx *fasthttp.RequestCtx)  {
	inter := interactors.NewInteractor()
	inter.RestToRest(ctx)
}

func (r rest) remote(ctx *fasthttp.RequestCtx) {
	inter := interactors.NewInteractor()
	inter.RestToGrp(ctx)
}

func (r rest) managerCreate(ctx *fasthttp.RequestCtx) {
	inter := interactors.NewInteractor()
	inter.CreateTopic(ctx)
}

func (r rest) managerUpdate(ctx *fasthttp.RequestCtx) {
	inter := interactors.NewInteractor()
	inter.UpdateTopic(ctx)
}

func (r rest) managerGet(ctx *fasthttp.RequestCtx) {
	inter := interactors.NewInteractor()
	inter.GetTopic(ctx)
}


func (r rest) path() {
	r.router.POST("/rest-local", r.local)
	r.router.POST("/rest-remote", r.remote)
	r.router.POST("/manager", r.managerCreate)
	r.router.PATCH("/manager", r.managerUpdate)
	r.router.GET("/manager", r.managerGet)
}


func NewRestApi(r *router.Router) {
	re := &rest{router: r}
	re.path()
}
