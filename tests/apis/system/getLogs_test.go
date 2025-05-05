package system_test

import (
	"dominus-project/app/interactors/system"
	"dominus-project/app/interfaces/fasthttp/input"
	"dominus-project/tests/mocks"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func TestGetLogsController(t *testing.T) {
	t.Run("Response_StatusOK", func(t *testing.T) {
		router := router.New()
		log := mocks.NewEventMock(true)
		system := system.NewSystemService(log)
		input.NewSystemAPI(router, system)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
	})

	t.Run("Response_Error", func(t *testing.T) {
		router := router.New()
		log := mocks.NewEventMock(false)
		system := system.NewSystemService(log)
		input.NewSystemAPI(router, system)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})
}