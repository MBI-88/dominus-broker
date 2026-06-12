package streambidirectionaltest

import (
	"context"
	streambidirectional "dominus-broker/internal/application/usecases/stream_bidirectional"
	"dominus-broker/mocks"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	c := mocks.NewMockBrokerClient(ctrl)
	b := streambidirectional.New(c)
	if b == nil {
		t.Fatal("StreamBidirectional returned nil")
	}
}

func TestStreamBidirectioal(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockBrokerClient, *mocks.MockBrokerBidirectionalDto, *mocks.MockBrokerRequestDto)
		output    error
	}{
		{
			name: "StreamBidirectional Ok",
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
					BidirectionalStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(urls, streamProv, streamSub, txArg, ctx any) {
						prov := streamProv.(<-chan []byte)
						sub := streamSub.(chan<- []byte)
						txCh := txArg.(chan<- struct{})
						cctx := ctx.(context.Context)

						go func() {
							for range prov {
							}
							// Let the handler drain streamSub for Send() expectations before tx closes the session.
							time.Sleep(100 * time.Millisecond)
							// Always signal tx after prov closes so StreamBidirectional can leave <-closed even if ctx was canceled first.
							txCh <- struct{}{}
						}()

						go func() {
							tick := time.NewTicker(5 * time.Millisecond)
							defer tick.Stop()
							for {
								select {
								case <-cctx.Done():
									return
								case <-tick.C:
									select {
									case sub <- []byte("test"):
									default:
									}
								}
							}
						}()
					}).
					Times(1)

				gomock.InOrder(
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(nil),
					mdto.EXPECT().Send(gomock.All()).Return(fmt.Errorf("connection closed")),
				)

			},
			output: fmt.Errorf("streambidirectional.StreamBidirectional connection closed"),
		},
		{
			name: "StreamBidirectional Recv error",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerBidirectionalDto, mrq *mocks.MockBrokerRequestDto) {
				gomock.InOrder(
					mdto.EXPECT().Context().Return(context.Background()),
					mdto.EXPECT().Recv().Return(nil, fmt.Errorf("recv error")),
				)
			},
			output: fmt.Errorf("streambidirectional.StreamBidirectional recv error"),
		},
		{
			name: "StreamBidirectional subscribers empty",
			setupMock: func(mc *mocks.MockBrokerClient, mdto *mocks.MockBrokerBidirectionalDto, mrq *mocks.MockBrokerRequestDto) {
				gomock.InOrder(
					mdto.EXPECT().Context().Return(context.Background()),
					mdto.EXPECT().Recv().Return(mrq, nil),
				)

				mrq.EXPECT().
					GetSubscribers().
					Return([]string{}).Times(1)
			},
			output: fmt.Errorf("streambidirectional.StreamBidirectional subscribers not found"),
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

			bidirectionalUseCase := streambidirectional.New(mockClient)

			err := bidirectionalUseCase.StreamBidirectional(mockDto)

			if tt.output.Error() != err.Error() {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
		})

	}
}
