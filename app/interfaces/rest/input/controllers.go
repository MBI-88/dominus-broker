package input

import (
	"dominus/app/interactors"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type rest struct {
	router *router.Router
	inter interactors.InteractorInt
}

//Works with connections rest to rest
//
//Parameters
//
//* ctx: context fasthttp
func (r rest) local(ctx *fasthttp.RequestCtx)  {
	inter := r.inter.NewConnection()
	inter.RestToRest(ctx)
}

//Works with connections rest to rpc
//
//Parameters
//
//* ctx: context fasthttp
func (r rest) remote(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewConnection()
	inter.RestToGrp(ctx)
}

//Crate topic in memory and database
//
//Parameters
//
//* ctx: context fasthttp
func (r rest) managerCreate(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	inter.CreateTopic(ctx)
}

//Update topic in memory and database
//
//Parameters
//
//* ctx: context fasthttp
func (r rest) managerUpdate(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	inter.UpdateTopic(ctx)
}

//Returns topics in the database
//
//Parameters
//
//* ctx: context fasthttp
func (r rest) managerGet(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	inter.GetTopic(ctx)
}

//Deletes a topic in memory and database
//
//Parameters
//
//* ctx: context fasthttp
func (r rest) managerDelete(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	inter.DeleteTopic(ctx)
}

//Path connects handler with the router
func (r rest) path() {
	r.router.POST("/local-shipping", r.local)
	r.router.POST("/remote-shipping", r.remote)
	r.router.POST("/manager", r.managerCreate)
	r.router.PATCH("/manager", r.managerUpdate)
	r.router.GET("/manager", r.managerGet)
	r.router.DELETE("/manager",r.managerDelete)
}

//Create a new Rest api service
func NewRestApi(r *router.Router, i interactors.InteractorInt) {
	re := &rest{router: r, inter: i}
	re.path()
}
