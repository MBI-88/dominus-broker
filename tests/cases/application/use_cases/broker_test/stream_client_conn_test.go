package broker_test

import (
	"context"
	"dominus-project/internal/application/use_cases/broker"
	"dominus-project/mocks"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestStreamClientConn(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockBrokerClient, *mocks.MockBrokerClientDto, *mocks.MockBrokerRequestDto)
		output    error
	}{
		{
			name: "StreamClientConn Ok",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerClientDto, mrq *mocks.MockBrokerRequestDto) {

				gomock.InOrder(
					mdto.EXPECT().Context().Return(context.Background()),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(mrq, nil),
					mdto.EXPECT().Recv().Return(nil, fmt.Errorf("end")),
				)

				mrq.EXPECT().
					GetSubscribers().
					Return([]string{"server1.api.com", "server2.api.com", "server3.api.com", "127.0.0.1:8080"}).AnyTimes()

				mc.EXPECT().
					ClientStream(gomock.All(), gomock.All(), gomock.All()).
					Do(func(urls, m, c any) {
						var ok = true
						for ok {
							_, ok = <-m.(<-chan []byte)
						}
					}).Times(1)

				mrq.EXPECT().
					GetPayload().
					Return([]byte("test-1")).AnyTimes()
			},
			output: fmt.Errorf("end"),
		},
		{
			name: "StreamClientConn Recv error",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerClientDto, mrq *mocks.MockBrokerRequestDto) {
				mdto.EXPECT().Context().Return(context.Background())
				mdto.EXPECT().Recv().Return(nil, fmt.Errorf("recv error"))
			},
			output: fmt.Errorf("recv error"),
		},
		{
			name: "StreamClientConn subscribers empty",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerClientDto, mrq *mocks.MockBrokerRequestDto) {
				mdto.EXPECT().Context().Return(context.Background())
				mdto.EXPECT().Recv().Return(mrq, nil)

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

			mockClientDto := mocks.NewMockBrokerClientDto(ctrl)
			mockClient := mocks.NewMockBrokerClient(ctrl)
			mockRequestDto := mocks.NewMockBrokerRequestDto(ctrl)

			tt.setupMock(mockClient, mockClientDto, mockRequestDto)

			service := broker.NewBroker(mockClient)

			err := service.StreamClientConn(mockClientDto)

			if tt.output != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.output)
				}
				if tt.output.Error() != err.Error() {
					t.Fatalf("expected error %q, got %q", tt.output.Error(), err.Error())
				}
			} else if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

		})

	}
}
