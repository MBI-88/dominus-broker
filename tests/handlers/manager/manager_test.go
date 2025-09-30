package manager_test

import (
	"dominus-project/mocks"
	"fmt"
	"testing"

	"dominus-project/internal/interfaces/fasthttp/input"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

//"dominus-project/internal/interfaces/fasthttp/input"
//"encoding/json"
//"testing"

// "github.com/valyala/fasthttp"

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

func TestAddTopic(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*mocks.MockManagerService)
		statusCode int
	}{
		{
			name: "AddTopic Ok",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					AddTopic(gomock.All()).
					Return(nil).Times(1)
			},
			statusCode: fasthttp.StatusNoContent,
		},
		{
			name: "AddTopic error",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					AddTopic(gomock.All()).
					Return(fmt.Errorf("error")).Times(1)
			},
			statusCode: fasthttp.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		router := router.New()
		managerMock := mocks.NewMockManagerService(ctrl)

		tt.setupMock(managerMock)

		input.NewManagerAPI(router, managerMock)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/v1/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.SetBody([]byte{})
		router.Handler(ctx)

		if ctx.Response.StatusCode() != tt.statusCode {
			t.Fatalf("[-] Expected %d received %d", tt.statusCode, ctx.Response.StatusCode())
		}

	}
}

func TestUpdateSubscribers(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*mocks.MockManagerService)
		statusCode int
	}{
		{
			name: "UpdateSubscribers Ok",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					UpdateSubscribers(gomock.All()).
					Return(nil).Times(1)
			},
			statusCode: fasthttp.StatusNoContent,
		},
		{
			name: "UpdateSubscribers error",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					UpdateSubscribers(gomock.All()).
					Return(fmt.Errorf("error")).Times(1)
			},
			statusCode: fasthttp.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		router := router.New()
		managerMock := mocks.NewMockManagerService(ctrl)

		tt.setupMock(managerMock)

		input.NewManagerAPI(router, managerMock)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/v1/subscribers/mysubscriber")
		ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
		ctx.Request.SetBody([]byte{})
		router.Handler(ctx)

		if ctx.Response.StatusCode() != tt.statusCode {
			t.Fatalf("[-] Expected %d received %d", tt.statusCode, ctx.Response.StatusCode())
		}

	}
}

func TestDeleteTopic(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*mocks.MockManagerService)
		statusCode int
	}{
		{
			name: "DeleteTopic Ok",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					DeleteTopic(gomock.All()).
					Return(nil).Times(1)
			},
			statusCode: fasthttp.StatusNoContent,
		},
		{
			name: "DeleteTopic error",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					DeleteTopic(gomock.All()).
					Return(fmt.Errorf("error")).Times(1)
			},
			statusCode: fasthttp.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		router := router.New()
		managerMock := mocks.NewMockManagerService(ctrl)

		tt.setupMock(managerMock)

		input.NewManagerAPI(router, managerMock)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/v1/topics/test_topic")
		ctx.Request.Header.SetMethod(fasthttp.MethodDelete)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != tt.statusCode {
			t.Fatalf("[-] Expected %d received %d", tt.statusCode, ctx.Response.StatusCode())
		}

	}

}

func TestGetTopics(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*mocks.MockManagerService)
		statusCode int
	}{
		{
			name: "GetTopics Ok",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					GetQueueInfo().
					Return(map[string]any{
						"data": map[string]any{
							"queue": 10,
						},
					}).Times(1)
			},
			statusCode: fasthttp.StatusOK,
		},
		{
			name: "GetTopics error",
			setupMock: func(mms *mocks.MockManagerService) {
				mms.EXPECT().
					GetQueueInfo().
					Return(map[string]any{
						"channel": make(chan int),
					}).Times(1)
			},
			statusCode: fasthttp.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		router := router.New()
		managerMock := mocks.NewMockManagerService(ctrl)

		tt.setupMock(managerMock)

		input.NewManagerAPI(router, managerMock)
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.SetRequestURI("/v1/topics")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		router.Handler(ctx)

		if ctx.Response.StatusCode() != tt.statusCode {
			t.Fatalf("[-] Expected %d received %d", tt.statusCode, ctx.Response.StatusCode())
		}

	}
}
