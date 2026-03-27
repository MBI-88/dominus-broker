package broker_test

import (
	"context"
	"dominus-project/internal/application/use_cases/broker"
	"dominus-project/mocks"
	"fmt"
	"sync"
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
					mdto.EXPECT().Context().Return(context.Background()),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(nil, fmt.Errorf("connection closed")),
				)

				mrq.EXPECT().
					GetSubscribers().
					Return([]string{"server1.api.com", "server2.api.com", "server3.api.com", "127.0.0.1:8080"}).
					AnyTimes()

				mrq.EXPECT().
					GetPayload().
					Return([]byte("test")).
					AnyTimes()

				mc.EXPECT().
					BidirectionalStream(gomock.All(), gomock.All(), gomock.All(), gomock.All(), gomock.All(), gomock.All(), gomock.All()).
					Do(func(subscribers, streamProv, streamSub, errMsg, closed, ctx, done any) {
						var recvWG sync.WaitGroup
						recvWG.Add(1)
						go func() {
							defer recvWG.Done()
							for msg := range streamProv.(<-chan []byte) {
								_ = msg
								errMsg.(chan<- error) <- fmt.Errorf("test error")
							}
							closed.(chan<- struct{}) <- struct{}{}
						}()

						subTicker := time.NewTicker(time.Second)
						defer subTicker.Stop()
						bctx := ctx.(context.Context)
						for {
							select {
							case <-subTicker.C:
								streamSub.(chan<- []byte) <- []byte("test")
							case <-bctx.Done():
								subTicker.Stop()
								recvWG.Wait()
								for range subscribers.([]string) {
									done.(chan<- struct{}) <- struct{}{}
								}
								return
							}
						}
					}).
					Times(1)

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
				gomock.InOrder(
					mdto.EXPECT().Context().Return(context.Background()),
					mdto.EXPECT().Recv().Return(nil, fmt.Errorf("recv error")),
				)
			},
			output: fmt.Errorf("recv error"),
		},
		{
			name: "StreamBiConn subscribers empty",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerBidirectionalDto, mrq *mocks.MockBrokerRequestDto) {
				gomock.InOrder(
					mdto.EXPECT().Context().Return(context.Background()),
					mdto.EXPECT().Recv().Return(mrq, nil),
				)

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
