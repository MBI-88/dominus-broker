package apis_test

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/interactors"
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
		topics := entities.NewTopics(100)
		inter := interactors.NewInteractor(log, topics, 100)
		inter = inter.Set(mocks.NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() !=  fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
	})

	t.Run("Response_Error", func(t *testing.T) {
		router := router.New()
		log := mocks.NewEventMock(false)
		topics := entities.NewTopics(100)
		inter := interactors.NewInteractor(log, topics, 100)
		inter = inter.Set(mocks.NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusInternalServerError {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
		}
	})
}

func TestGetBackupController(t *testing.T) {
	t.Run("Response_StatusOK", func(t *testing.T) {
		router := router.New()
		log := mocks.NewEventMock(true)
		topics := entities.NewTopics(100)
		inter := interactors.NewInteractor(log, topics, 100)
		inter = inter.Set(mocks.NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs-backup")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
	})

	t.Run("Response_Error", func(t *testing.T) {
		router := router.New()
		log := mocks.NewEventMock(false)
		topics := entities.NewTopics(100)
		inter := interactors.NewInteractor(log, topics, 100)
		inter = inter.Set(mocks.NewGrpcClientMock())
		input.NewSystemAPI(router, inter.NewSystemService())
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/logs-backup")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusInternalServerError {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
		}
	})
}
