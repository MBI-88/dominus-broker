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
		setupMock func(*mocks.MockBrokerClient, *mocks.MockBrokerServerDto, *mocks.MockBrokerRequestDto)
		output    error
	}{
		{
			name: "StreamServerConn ok",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerServerDto, mrq *mocks.MockBrokerRequestDto) {

				mrq.EXPECT().
					GetSubscribers().
					Return([]string{"server1.api.com", "server2.api.com", "server3.api.com", "127.0.0.1:8080"}).AnyTimes()

				mrq.EXPECT().
					GetPayload().
					Return([]byte("test")).Times(1)

				mc.EXPECT().
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
					mdto.EXPECT().Context().Return(context.Background()),
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(fmt.Errorf("connection closed")),
				)
			},
			output: fmt.Errorf("connection closed"),
		},
		{
			name: "StreamServerConn subscribers empty",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerServerDto, mrq *mocks.MockBrokerRequestDto) {
				mdto.EXPECT().Context().Return(context.Background())
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

			mockDto := mocks.NewMockBrokerServerDto(ctrl)
			mockClient := mocks.NewMockBrokerClient(ctrl)
			mockRequestDto := mocks.NewMockBrokerRequestDto(ctrl)

			tt.setupMock(mockClient, mockDto, mockRequestDto)

			service := broker.NewBroker(mockClient)

			err := service.StreamServerConn(mockRequestDto, mockDto)

			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}

		})

	}
}
