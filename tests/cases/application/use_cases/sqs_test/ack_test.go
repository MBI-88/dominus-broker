package sqs_test

import (
	"context"
	"dominus-broker/internal/application/usecases/sqs"
	"dominus-broker/mocks"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAck(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockConsumerDto, *mocks.MockMemoryClient)
		output    error
	}{
		{
			name: "Ack ok",
			setupMock: func(dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient) {
				// Ack llama GetMessageId dos veces: validación y AckMessage.
				dto.EXPECT().GetMessageId().Return("1700000000001-0").Times(2)
				dto.EXPECT().GetGroupId().Return("group-1")

				mc.EXPECT().
					AckMessage(gomock.Any(), "1700000000001-0", "group-1").
					Return(nil)
			},
			output: nil,
		},
		{
			name: "Ack connection error",
			setupMock: func(dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().GetMessageId().Return("1700000000001-0").Times(2)
				dto.EXPECT().GetGroupId().Return("group-1")

				mc.EXPECT().
					AckMessage(gomock.Any(), "1700000000001-0", "group-1").
					Return(fmt.Errorf("connection error"))
			},
			output: fmt.Errorf("connection error"),
		},
		{
			name: "Ack invalid message id arbitrary string",
			setupMock: func(dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().GetMessageId().Return("message-1").Times(1)
			},
			output: fmt.Errorf("sqs.Ack invalid messageId"),
		},
		{
			name: "Ack invalid message id without sequence part",
			setupMock: func(dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().GetMessageId().Return("123456789").Times(1)
			},
			output: fmt.Errorf("sqs.Ack invalid messageId"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockDto := mocks.NewMockConsumerDto(ctrl)
			mockClient := mocks.NewMockMemoryClient(ctrl)

			tt.setupMock(mockDto, mockClient)
			service := sqs.NewSqs(mockClient)
			err := service.Ack(ctx, mockDto)

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
