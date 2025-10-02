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
		t.Run(tt.name, func(t *testing.T) {
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

		})

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
		t.Run(tt.name, func(t *testing.T) {
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

		})

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
		t.Run(tt.name, func(t *testing.T) {
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
		})

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
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}
