package swagger_test

import (
	"dominus-project/internal/infrastructure/fasthttp/input"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func TestSwaggerController(t *testing.T) {
	t.Run("Swagger_Ok", func(t *testing.T) {
		router := router.New()
		input.NewSwaggerAPI(router)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/swagger/index.html")
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}

	})

	t.Run("Swagger_Error", func(t *testing.T) {
		router := router.New()
		input.NewSwaggerAPI(router)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/swagger/")
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotFound {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotFound, ctx.Response.StatusCode())
		}
	})

}
