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

func TestStreamBiConn(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockBrokerClient, *mocks.MockBrokerBidirectionalDto, *mocks.MockBrokerRequestDto)
		output    error
	}{
		{
			name: "StreamBiConn Ok",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerBidirectionalDto, mrq *mocks.MockBrokerRequestDto) {
				
				gomock.InOrder(
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(nil, fmt.Errorf("connection closed")),
				)

				mrq.EXPECT().
					GetSubscribers().
					Return([]string{"server1.api.com", "server2.api.com", "server3.api.com", "127.0.0.1:8080"}).AnyTimes()

				mrq.EXPECT().
					GetPayload().
					Return([]byte("test")).AnyTimes()

				mc.EXPECT().
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
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(fmt.Errorf("connection closed")),
				)

			},
			output: fmt.Errorf("connection closed"),
		},
		{
			name: "StreamBiConn Recv error",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerBidirectionalDto, mrq *mocks.MockBrokerRequestDto) {
				mdto.EXPECT().Recv().Return(nil, fmt.Errorf("recv error"))
			},
			output: fmt.Errorf("recv error"),
		},
		{
			name: "StreamBiConn subscribers empty",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerBidirectionalDto, mrq *mocks.MockBrokerRequestDto) {
				mdto.EXPECT().Recv().Return(mrq, nil).Times(1)

				mrq.EXPECT().
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

			mockDto := mocks.NewMockBrokerBidirectionalDto(ctrl)
			mockClient := mocks.NewMockBrokerClient(ctrl)
			mockRequestDto := mocks.NewMockBrokerRequestDto(ctrl)

			tt.setupMock(mockClient, mockDto, mockRequestDto)

			service := broker.NewBroker(mockClient)

			err := service.StreamBiConn(mockDto)

			if tt.output.Error() != err.Error() {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
		})

	}
}
