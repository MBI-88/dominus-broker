package manager_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/rules"
	"dominus-project/internal/interactors/manager"
	"dominus-project/internal/interfaces/fasthttp/input"
	"dominus-project/tests/env"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)



func TestManagerController(t *testing.T) {
	topics := entities.NewTopics(100)
	rls := rules.NewValidator()
	manager := manager.NewManagerService(topics, rls)

	t.Run("AddTopic_Ok", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topic")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_create_body.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}

	})
	t.Run("AddTopic_Error", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topic")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_create_body_error.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("UpdateSubscribers_Ok", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/subscribers/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodPut)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_update_body.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}

	})
	t.Run("UpdateSubscribers_Error", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/subscribers/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodPut)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_update_error.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("DeleteTopic_Ok", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topic/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodDelete)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}
	})

	t.Run("DeleteTopic_Error", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topic/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodDelete)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})
}
