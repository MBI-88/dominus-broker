package input

import (
	"dominus/app/interactors"

	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	//"github.com/valyala/fasthttp"
)

type rest struct {
	router *router.Router
	inter  interactors.InteractorInt
	js     jsoniter.API
}



/*
// Crate topic in memory and database
//
// # Parameters
//
// -> ctx: context fasthttp
func (r *rest) managerCreate(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	context := NewRestContext(ctx)
	
	if err := inter.CreateTopic(context); err != nil {
		message := make(map[string]any)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		message["message"] = err.Error()
		b, _ := r.js.Marshal(message)
		ctx.Response.SetBody(b)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
}

// Update topic in memory and database
//
// # Parameters
//
// -> ctx: context fasthttp
func (r *rest) managerUpdate(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	context := NewRestContext(ctx)

	if err := inter.UpdateTopic(context); err != nil {
		message := make(map[string]any)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		message["message"] = err.Error()
		b, _ := r.js.Marshal(message)
		ctx.Response.SetBody(b)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
}

// Returns topics in the database
//
// # Parameters
//
// -> ctx: context fasthttp
func (r *rest) managerGet(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	context := NewRestContext(ctx)
	
	message, err := inter.GetTopic(context)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		message["message"] = err.Error()
		b, _ := r.js.Marshal(message)
		ctx.Response.SetBody(b)
		return
	}

	body, err := r.js.Marshal(message)
	if err != nil {
		delete(message, "topic")
		message["message"] = err.Error()
	}

	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(body)
}

// Deletes a topic in memory and database
//
// # Parameters
//
// -> ctx: context fasthttp
func (r *rest) managerDelete(ctx *fasthttp.RequestCtx) {
	inter := r.inter.NewManager()
	context := NewRestContext(ctx)

	if err := inter.DeleteTopic(context); err != nil {
		message := make(map[string]any)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusExpectationFailed)
		message["message"] = err.Error()
		b, _ := r.js.Marshal(message)
		ctx.Response.SetBody(b)
	}

	ctx.Response.Header.Set("Content-Type", "application/text")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
}

// Path connects handler with the router
func (r *rest) path() {
	r.router.POST("/manager", r.managerCreate)
	r.router.PATCH("/manager", r.managerUpdate)
	r.router.GET("/manager", r.managerGet)
	r.router.DELETE("/manager", r.managerDelete)
}

// Create a new Rest api service
func NewRestApi(r *router.Router, i interactors.InteractorInt) {
	re := &rest{router: r, inter: i, js: jsoniter.ConfigCompatibleWithStandardLibrary}
	re.path()
}

**/