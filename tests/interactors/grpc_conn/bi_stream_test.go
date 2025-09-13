package grpcconn_test

import (
	"context"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/mocks"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestStreamBiConn(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockStreamBi, *mocks.MockGrpcClient, *mocks.MockLogs, *mocks.MockGrpcDto)
		output    error
	}{
		{
			name: "StreamBiConn Ok",
			setupMock: func(mss *mocks.MockStreamBi, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {
				ml.EXPECT().
					WriteLog(gomock.All(), gomock.All()).AnyTimes()

				gomock.InOrder(
					mss.EXPECT().Recv().Return(mgd, nil),
					mss.EXPECT().Recv().Return(mgd, nil),
					mss.EXPECT().Recv().Return(mgd, nil),
					mss.EXPECT().Recv().Return(nil, fmt.Errorf("connection closed")),
				)

				mgd.EXPECT().
					GetSubscribers().
					Return([]string{"server1.api.com", "server2.api.com", "server3.api.com", "127.0.0.1:8080"}).AnyTimes()

				mgd.EXPECT().
					GetPayload().
					Return([]byte("test")).AnyTimes()

				mgc.EXPECT().
					BidirectionalStream(gomock.All(), gomock.All(), gomock.All(), gomock.All(), gomock.All(), gomock.All(), gomock.All()).
					Do(func(subscribers, streamProv, streamSub, errMsg, closed, ctx, done any) {
						// received
						go func() {
							for msg := range streamProv.(<-chan []byte) {
								println(msg)
								errMsg.(chan<- error) <- fmt.Errorf("test error")
							}
							closed.(chan<- struct{}) <- struct{}{}
						}()

						// subscribers sending messages
						for {
							select {
							case <-time.Tick(1 * time.Second):
								streamSub.(chan<- []byte) <- []byte("test")
							case <-time.Tick(5 * time.Second):
							case <-ctx.(context.Context).Done():
								for range subscribers.([]string) {
									done.(chan<- struct{}) <- struct{}{}
								}
								return
							}
						}
					}).Times(1)

				gomock.InOrder(
					mss.EXPECT().Send(gomock.All()).Return(nil),
					mss.EXPECT().Send(gomock.All()).Return(nil),
					mss.EXPECT().Send(gomock.All()).Return(nil),
					mss.EXPECT().Send(gomock.All()).Return(fmt.Errorf("connection closed")),
				)

			},
			output: fmt.Errorf("connection closed"),
		},
		{
			name: "StreamBiConn Recv error",
			setupMock: func(msb *mocks.MockStreamBi, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {
				msb.EXPECT().Recv().Return(nil, fmt.Errorf("recv error"))
			},
			output: fmt.Errorf("recv error"),
		},
		{
			name: "StreamBiConn subscribers empty",
			setupMock: func(msb *mocks.MockStreamBi, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {
				msb.EXPECT().Recv().Return(mgd, nil).Times(1)

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
			mockGrpcDto := mocks.NewMockGrpcDto(ctrl)
			mockGrpcClient := mocks.NewMockGrpcClient(ctrl)
			mockTopics := mocks.NewMockTopics(ctrl)
			mockStreamBi := mocks.NewMockStreamBi(ctrl)

			tt.setupMock(mockStreamBi, mockGrpcClient, mockLogs, mockGrpcDto)

			stream := grpcconn.NewGrpcService(mockLogs, mockGrpcClient, mockTopics)

			err := stream.StreamBiConn(mockStreamBi)

			if tt.output.Error() != err.Error() {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
		})

	}
}
