package grpcconn_test

import (
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/mocks"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestStreamClientConn(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockStreamClient, *mocks.MockGrpcClient, *mocks.MockLogs, *mocks.MockGrpcDto)
		output    error
	}{
		{
			name: "StreamClientConn Ok",
			setupMock: func(msc *mocks.MockStreamClient, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {

				gomock.InOrder(
					msc.EXPECT().Recv().Return(mgd, nil),
					msc.EXPECT().Recv().Return(mgd, nil),
					msc.EXPECT().Recv().Return(mgd, nil),
					msc.EXPECT().Recv().Return(nil, fmt.Errorf("end")),
				)

				mgd.EXPECT().
					GetSubscribers().
					Return([]string{"server1.api.com", "server2.api.com", "server3.api.com", "127.0.0.1:8080"}).AnyTimes()

				mgc.EXPECT().
					ClientStream(gomock.All(), gomock.All(), gomock.All()).
					Do(func(urls, m, c any) {
						var ok = true
						for ok {
							_, ok = <-m.(<-chan []byte)
						}
					}).Times(1)

				mgd.EXPECT().
					GetPayload().
					Return([]byte("test-1")).AnyTimes()

				ml.EXPECT().
					WriteLog(gomock.All(), gomock.All()).Times(1)
			},
			output: fmt.Errorf("end"),
		},
		{
			name: "StreamClientConn Recv error",
			setupMock: func(msc *mocks.MockStreamClient, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {
				msc.EXPECT().Recv().Return(nil, fmt.Errorf("recv error"))
			},
			output: fmt.Errorf("recv error"),
		},
		{
			name: "StreamClientConn subscribers empty",
			setupMock: func(msc *mocks.MockStreamClient, mgc *mocks.MockGrpcClient, ml *mocks.MockLogs, mgd *mocks.MockGrpcDto) {
				msc.EXPECT().Recv().Return(mgd, nil)

				mgd.EXPECT().
					GetSubscribers().
					Return([]string{}).Times(1)
			},
			output: fmt.Errorf("subscribers empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLogs := mocks.NewMockLogs(ctrl)
			mockStreaClient := mocks.NewMockStreamClient(ctrl)
			mockGrpcDto := mocks.NewMockGrpcDto(ctrl)
			mockGrpcClient := mocks.NewMockGrpcClient(ctrl)
			mockTopics := mocks.NewMockTopics(ctrl)

			tt.setupMock(mockStreaClient, mockGrpcClient, mockLogs, mockGrpcDto)

			stream := grpcconn.NewGrpcService(mockLogs, mockGrpcClient, mockTopics)

			err := stream.StreamClientConn(mockStreaClient)

			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}

		})

	}
}
