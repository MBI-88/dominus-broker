package sqs_test

import (
	"context"
	"dominus-broker/internal/application/use_cases/sqs"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/mocks"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestProducer(t *testing.T) {
	payloadOK := []byte("\n{payload:{data:data}}\n")
	tests := []struct {
		name      string
		setupMock func(*testing.T, *mocks.MockProducerDto, *mocks.MockMemoryClient)
		output    error
	}{
		{
			name: "Producer OK",
			setupMock: func(t *testing.T, dto *mocks.MockProducerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().
					GetPayload().
					Return(payloadOK)

				mc.EXPECT().
					SendMessage(gomock.Any(), gomock.AssignableToTypeOf(&entities.Message{})).
					Do(func(_ context.Context, q *entities.Message) {
						if string(q.GetMessage()) != string(payloadOK) {
							t.Errorf("SendMessage message body: got %q want %q", q.GetMessage(), payloadOK)
						}
						if q.GetMessageId() == "" {
							t.Error("SendMessage: expected non-empty stream-style message id")
						}
					}).
					Return(nil)
			},
			output: nil,
		},
		{
			name: "Producer connection error",
			setupMock: func(t *testing.T, dto *mocks.MockProducerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().
					GetPayload().
					Return(payloadOK)

				mc.EXPECT().
					SendMessage(gomock.Any(), gomock.AssignableToTypeOf(&entities.Message{})).
					Do(func(_ context.Context, q *entities.Message) {
						if string(q.GetMessage()) != string(payloadOK) {
							t.Errorf("SendMessage message body: got %q want %q", q.GetMessage(), payloadOK)
						}
					}).
					Return(fmt.Errorf("connection error"))
			},
			output: fmt.Errorf("connection error"),
		},
		{
			name: "Producer empty payload",
			setupMock: func(_ *testing.T, dto *mocks.MockProducerDto, _ *mocks.MockMemoryClient) {
				dto.EXPECT().
					GetPayload().
					Return([]byte{})
			},
			output: fmt.Errorf("sqs.Producer empty payload"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockDto := mocks.NewMockProducerDto(ctrl)
			mockClient := mocks.NewMockMemoryClient(ctrl)

			tt.setupMock(t, mockDto, mockClient)
			service := sqs.NewSQS(mockClient)
			err := service.Producer(ctx, mockDto)

			if tt.output == nil {
				if err != nil {
					t.Fatalf("expected nil error got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %v got nil", tt.output)
			}

			if tt.output.Error() != err.Error() {
				t.Fatalf("expected %s got %s", tt.output, err)
			}

		})
	}
}
