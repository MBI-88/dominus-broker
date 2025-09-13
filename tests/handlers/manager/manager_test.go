package manager_test

import (

//"dominus-project/internal/interfaces/fasthttp/input"
//"encoding/json"
//"testing"

// "github.com/fasthttp/router"
// "github.com/valyala/fasthttp"
)

/*
func TestManagerController(t *testing.T) {

	t.Run("AddTopic_Ok", func(t *testing.T) {
		topics := entities.NewTopics(env.TopicLimit)
		router := router.New()
		manager := manager.NewManagerService(topics, env.QueueLimit)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody(env.ReadJson("./../../samples/rest_create_body.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, ctx.Response.StatusCode())
		}

	})
	t.Run("AddTopic_Error", func(t *testing.T) {
		topics := entities.NewTopics(env.TopicLimit)
		router := router.New()
		manager := manager.NewManagerService(topics, env.QueueLimit)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody(env.ReadJson("./../../samples/rest_create_body_error.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("UpdateSubscribers_Ok", func(t *testing.T) {
		topics := entities.NewTopics(env.TopicLimit)
		router := router.New()
		c := new(fasthttp.RequestCtx)
		c.Request.SetRequestURI("/subscribers/test")
		c.Request.Header.SetMethod(fasthttp.MethodPatch)
		c.Request.SetBody(env.ReadJson("./../../samples/rest_update_body.json"))

		ctx := dto.NewRestContext(c)
		topic := entities.NewTopic(ctx, env.QueueLimit)
		topic.SetName("test")
		topics.Append(topic)

		manager := manager.NewManagerService(topics, env.QueueLimit)
		input.NewManagerAPI(router, manager)
		router.Handler(c)

		if c.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, c.Response.StatusCode())
		}

	})
	t.Run("UpdateSubscribers_Error", func(t *testing.T) {
		topics := entities.NewTopics(env.TopicLimit)
		router := router.New()
		manager := manager.NewManagerService(topics, env.QueueLimit)
		input.NewManagerAPI(router, manager)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/subscribers/test")
		ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
		ctx.Request.SetBody(env.ReadJson("./../../samples/rest_update_error.json"))
		router.Handler(ctx)

		if ctx.Response.StatusCode() != fasthttp.StatusNotAcceptable {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNotAcceptable, ctx.Response.StatusCode())
		}
	})

	t.Run("DeleteTopic_Ok", func(t *testing.T) {
		topics := entities.NewTopics(10)
		router := router.New()
		c := new(fasthttp.RequestCtx)
		c.Request.SetRequestURI("/topics/test")
		c.Request.Header.SetMethod(fasthttp.MethodDelete)

		ctx := dto.NewRestContext(c)
		topic := entities.NewTopic(ctx, env.QueueLimit)
		topic.SetName("test")
		topics.Append(topic)

		manager := manager.NewManagerService(topics, env.QueueLimit)
		input.NewManagerAPI(router, manager)
		router.Handler(c)

		if c.Response.StatusCode() != fasthttp.StatusNoContent {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusNoContent, c.Response.StatusCode())
		}
	})

	t.Run("DeleteTopic_Error", func(t *testing.T) {
		topics := entities.NewTopics(env.TopicLimit)
		router := router.New()
		manager := manager.NewManagerService(topics, env.QueueLimit)
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
		topics := entities.NewTopics(env.TopicLimit)
		router := router.New()
		c := new(fasthttp.RequestCtx)
		c.Request.SetRequestURI("/topics")
		c.Request.Header.SetMethod(fasthttp.MethodGet)

		ctx := dto.NewRestContext(c)
		topic := entities.NewTopic(ctx, env.QueueLimit)
		topic.SetName("test")
		topic.SetSubscribers([]string{"server1.api.com", "server2.api.com", "server3.api.com"})
		topic.SetMessage([]byte("test for testing"))

		topics.Append(topic)

		manager := manager.NewManagerService(topics, env.QueueLimit)
		input.NewManagerAPI(router, manager)

		router.Handler(c)

		if c.Response.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("[-] Expected %d received %d", fasthttp.StatusOK, c.Response.StatusCode())
		}
		if body := c.Response.Body(); len(body) == 0 {
			t.Fatal("Body must be fullfit")
		}
	})

	t.Run("GetTopicInfo_ERROR", func(t *testing.T) {
		topics := entities.NewTopics(env.TopicLimit)
		router := router.New()

		manager := manager.NewManagerService(topics, env.QueueLimit)
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
*/
