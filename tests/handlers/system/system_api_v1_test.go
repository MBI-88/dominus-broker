package system_test

import (
	"dominus-project/internal/interfaces/fasthttp/input"
	"dominus-project/mocks"
	"fmt"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

func TestGetBackupController(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*mocks.MockSystemService, *mocks.MockLogs)
		statusCode int
	}{
		{
			name: "GetBackup Ok",
			setupMock: func(mss *mocks.MockSystemService, ml *mocks.MockLogs) {
				mss.EXPECT().
					GetLogs().
					Return([]string{"test1", "test2"}, nil).Times(1)
			},
			statusCode: fasthttp.StatusOK,
		},
		{
			name: "GetBackup error",
			setupMock: func(mss *mocks.MockSystemService, ml *mocks.MockLogs) {
				mss.EXPECT().
					GetLogs().
					Return([]string{}, fmt.Errorf("error")).Times(1)
			},
			statusCode: fasthttp.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockSystemService(ctrl)
			mockLogs := mocks.NewMockLogs(ctrl)

			tt.setupMock(mockService, mockLogs)

			router := router.New()
			input.NewSystemAPI(router, mockService)
			ctx := new(fasthttp.RequestCtx)
			ctx.Request.SetRequestURI("/v1/backup")
			ctx.Request.Header.SetMethod(fasthttp.MethodGet)
			router.Handler(ctx)

			if tt.statusCode != ctx.Response.StatusCode() {
				t.Fatalf("[-] Expected %d received %d", tt.statusCode, ctx.Response.StatusCode())
			}

		})
	}
}

func TestGetLogsController(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*mocks.MockSystemService, *mocks.MockLogs)
		statusCode int
	}{
		{
			name: "GetBackup Ok",
			setupMock: func(mss *mocks.MockSystemService, ml *mocks.MockLogs) {
				mss.EXPECT().
					GetLogs().
					Return([]string{"test1", "test2"}, nil).Times(1)
			},
			statusCode: fasthttp.StatusOK,
		},
		{
			name: "GetBackup error",
			setupMock: func(mss *mocks.MockSystemService, ml *mocks.MockLogs) {
				mss.EXPECT().
					GetLogs().
					Return([]string{}, fmt.Errorf("error")).Times(1)
			},
			statusCode: fasthttp.StatusNotAcceptable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockSystemService(ctrl)
			mockLogs := mocks.NewMockLogs(ctrl)

			tt.setupMock(mockService, mockLogs)

			router := router.New()
			input.NewSystemAPI(router, mockService)
			ctx := new(fasthttp.RequestCtx)
			ctx.Request.SetRequestURI("/v1/logs")
			ctx.Request.Header.SetMethod(fasthttp.MethodGet)
			router.Handler(ctx)

			if tt.statusCode != ctx.Response.StatusCode() {
				t.Fatalf("[-] Expected %d received %d", tt.statusCode, ctx.Response.StatusCode())
			}

		})
	}
}
