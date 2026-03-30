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
					ServerStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(urls, initialRequest, stream, ctx, tx any) {
						ch := stream.(chan<- []byte)
						cctx := ctx.(context.Context)
						txCh := tx.(chan<- struct{})
						tick := time.NewTicker(5 * time.Millisecond)
						defer tick.Stop()
						for {
							select {
							case <-tick.C:
								select {
								case ch <- []byte("test"):
								default:
								}
							case <-cctx.Done():
								txCh <- struct{}{}
								return
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
		{
			name: "StreamServerConn send error then connection closed",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerServerDto, mrq *mocks.MockBrokerRequestDto) {
				mrq.EXPECT().GetSubscribers().Return([]string{"a"}).Times(1)
				mrq.EXPECT().GetPayload().Return([]byte("init")).Times(1)
				mdto.EXPECT().Context().Return(context.Background()).Times(1)
				mc.EXPECT().
					ServerStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(_, _, stream, ctx, tx any) {
						ch := stream.(chan<- []byte)
						cctx := ctx.(context.Context)
						txCh := tx.(chan<- struct{})
						go func() {
							select {
							case ch <- []byte("body"):
							case <-cctx.Done():
								return
							}
							<-cctx.Done()
							txCh <- struct{}{}
						}()
					}).Times(1)
				mdto.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("send failed")).Times(1)
			},
			output: fmt.Errorf("connection closed"),
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
