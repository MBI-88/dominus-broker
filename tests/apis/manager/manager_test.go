package manager_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/rules"
	"dominus-project/internal/interactors/manager"
	"dominus-project/internal/interfaces/fasthttp/input"
	"os"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func readJson(path string) []byte {
	var file []byte
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}
	file, _ = os.ReadFile(path)
	return file
}

func TestManagerController(t *testing.T) {
	topics := entities.NewTopics(100)
	rls := rules.NewValidator()
	manager := manager.NewManagerService(topics, rls, 100)

	t.Run("AddTopic_Ok", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topic")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody(readJson("./../../mocks/rest_create_body.json"))
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
		ctx.Request.SetBody(readJson("./../../mocks/rest_create_body_error.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("UpdatePartition_Ok", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/partition/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodPut)
		ctx.Request.SetBody(readJson("./../../mocks/rest_update_body.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}

	})
	t.Run("UpdatePatition_Error", func(t *testing.T) {
		router := router.New()
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/partition/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodPut)
		ctx.Request.SetBody(readJson("./../../mocks/rest_update_error.json"))
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
