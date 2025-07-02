package manager_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/rules"
	"dominus-project/internal/interactors/manager"
	"dominus-project/internal/interfaces/fasthttp/input"
	"dominus-project/tests/env"
	"encoding/json"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func TestManagerController(t *testing.T) {
	rls := rules.NewValidator()

	t.Run("AddTopic_Ok", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()
		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_create_body.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}

	})
	t.Run("AddTopic_Error", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()
		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_create_body_error.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("UpdateSubscribers_Ok", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()
		topic := new(entities.Topic)

		data := env.ReadJson("./../../mocks/rest_create_body.json")
		if err := json.Unmarshal(data, topic); err != nil {
			t.Fatal(err)
		}

		topics.Append(topic)
		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/subscribers/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_update_body.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}

	})
	t.Run("UpdateSubscribers_Error", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()
		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/subscribers/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
		ctx.Request.SetBody(env.ReadJson("./../../mocks/rest_update_error.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("DeleteTopic_Ok", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()
		topic := new(entities.Topic)

		data := env.ReadJson("./../../mocks/rest_create_body.json")
		if err := json.Unmarshal(data, topic); err != nil {
			t.Fatal(err)
		}

		topics.Append(topic)
		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodDelete)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}
	})

	t.Run("DeleteTopic_Error", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()
		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodDelete)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("GetTopicInfo_OK", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()
		topic := new(entities.Topic)

		data := env.ReadJson("./../../mocks/rest_create_body.json")
		if err := json.Unmarshal(data, topic); err != nil {
			t.Fatal(err)
		}

		topics.Append(topic)
		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, ctx.Response.StatusCode())
		}
		if body := ctx.Response.Body(); len(body) == 0 {
			t.Fatal("Body must be fullfit")
		}
	})

	t.Run("GetTopicInfo_ERROR", func(t *testing.T) {
		topics := entities.NewTopics(100)
		router := router.New()

		manager := manager.NewManagerService(topics, rls)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		response := make(map[string]any, 3)

		if err := json.Unmarshal(ctx.Response.Body(), &response); err != nil {
			t.Fatal(err)
		}

		arrayTopics := response["topics"].([]any)
	
		if len(arrayTopics) > 0 {
			t.Fatal("Topics must be empty")
		}
	})
}
