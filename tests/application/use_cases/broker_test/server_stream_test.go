package broker_test

import (
	"context"
	"dominus-project/internal/application/use_cases/broker"
	"dominus-project/mocks"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestStreamServerConn(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockStreamServer, *mocks.MockGrpcClient, *mocks.MockLogs, *mocks.MockGrpcDto)
		output    error
	}{
		{
			name: "StreamServerConn ok",
			setupMock: func(mss *mocks.MockStreamServer, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {
				mgd.EXPECT().
					GetSubscribers().
					Return([]string{"server1.api.com", "server2.api.com", "server3.api.com", "127.0.0.1:8080"}).AnyTimes()

				mgd.EXPECT().
					GetPayload().
					Return([]byte("test")).Times(1)

				mgc.EXPECT().
					ServerStream(gomock.All(), gomock.All(), gomock.All(), gomock.All(), gomock.All()).
					Do(func(subscribers, initialRequest, stream, ctx, done any) {
						for {
							select {
							case <-time.Tick(1 * time.Second):
								stream.(chan<- []byte) <- []byte("test")
							case <-ctx.(context.Context).Done():
								for range subscribers.([]string) {
									done.(chan<- struct{}) <- struct{}{}
								}
								return
							case <-time.Tick(5 * time.Second):
							}
						}
					}).Times(1)

				gomock.InOrder(
					mss.EXPECT().Send(gomock.All()).Return(nil),
					mss.EXPECT().Send(gomock.All()).Return(nil),
					mss.EXPECT().Send(gomock.All()).Return(nil),
					mss.EXPECT().Send(gomock.All()).Return(fmt.Errorf("connection closed")),
				)

				ml.EXPECT().
					WriteLog(gomock.All(), gomock.All()).AnyTimes()
			},
			output: fmt.Errorf("connection closed"),
		},
		{
			name: "StreamServerConn subscribers empty",
			setupMock: func(mss *mocks.MockStreamServer, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {
				mgd.EXPECT().
					GetSubscribers().
					Return([]string{}).Times(1)

			},
			output: fmt.Errorf("subscribers not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLogs := mocks.NewMockLogs(ctrl)
			mockStreaServer := mocks.NewMockStreamServer(ctrl)
			mockGrpcDto := mocks.NewMockGrpcDto(ctrl)
			mockGrpcClient := mocks.NewMockGrpcClient(ctrl)

			tt.setupMock(mockStreaServer, mockGrpcClient, mockLogs, mockGrpcDto)

			stream := broker.NewBroker(mockLogs, mockGrpcClient)

			err := stream.StreamServerConn(mockGrpcDto, mockStreaServer)

			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}

		})

	}
}
