package tests

import (
	"dominus-project/app/domain/rules"
	"dominus-project/app/interactors"
	"dominus-project/app/interfaces/fasthttp/input"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func TestGetLogsController(t *testing.T) {
	rls := rules.NewRules()
	t.Run("Response_StatusOK", func(t *testing.T) {
		router := router.New()
		log := NewEventMock(true)
		inter := interactors.NewInteractor(log, rls)
		inter = inter.Set(NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)
		
		if fasthttp.StatusOK != ctx.Response.StatusCode() {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
	})

	t.Run("Response_Error", func(t *testing.T) {
		router := router.New()
		log := NewEventMock(false)
		inter := interactors.NewInteractor(log, rls)
		inter = inter.Set(NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if fasthttp.StatusNotFound != ctx.Response.StatusCode() {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotFound, ctx.Response.StatusCode())
		}
	})
}

func TestGetBackupController(t *testing.T) {
	rls := rules.NewRules()
	t.Run("Response_StatusOK", func(t *testing.T) {
		router := router.New()
		log := NewEventMock(true)
		inter := interactors.NewInteractor(log, rls)
		inter = inter.Set(NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs-backup")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if fasthttp.StatusOK != ctx.Response.StatusCode() {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
	})

	t.Run("Response_Error", func(t *testing.T) {
		router := router.New()
		log := NewEventMock(false)
		inter := interactors.NewInteractor(log, rls)
		inter = inter.Set(NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs-backup")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if fasthttp.StatusNotFound != ctx.Response.StatusCode() {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotFound, ctx.Response.StatusCode())
		}
	})
}