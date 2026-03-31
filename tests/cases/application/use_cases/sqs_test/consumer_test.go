package sqs_test

import (
	"context"
	"dominus-broker/internal/application/use_cases/sqs"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/mocks"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestConsumer(t *testing.T) {
	fixedAt := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	wantMsg := &entities.Message{
		Message:   []byte("consumer-ok-body"),
		MeesageId: "1700000000123-0",
		CreatedAt: fixedAt,
	}

	tests := []struct {
		name      string
		setupMock func(context.Context, *mocks.MockConsumerDto, *mocks.MockMemoryClient)
		output    error
		want      *entities.Message
	}{
		{
			name: "Consumer OK",
			setupMock: func(ctx context.Context, dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().
					GetWorkerId().
					Return("worker-1")

				dto.EXPECT().
					GetGroupId().
					Return("consumer-1")

				mc.EXPECT().
					GetMessage(ctx, "worker-1", "consumer-1").
					Return(wantMsg, nil)
			},
			output: nil,
			want:   wantMsg,
		},
		{
			name: "Consumer connection error",
			setupMock: func(ctx context.Context, dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().
					GetWorkerId().
					Return("worker-1")

				dto.EXPECT().
					GetGroupId().
					Return("consumer-1")

				mc.EXPECT().
					GetMessage(ctx, "worker-1", "consumer-1").
					Return(&entities.Message{}, fmt.Errorf("connection error"))
			},
			output: fmt.Errorf("connection error"),
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockDto := mocks.NewMockConsumerDto(ctrl)
			mockClient := mocks.NewMockMemoryClient(ctrl)

			tt.setupMock(ctx, mockDto, mockClient)
			service := sqs.NewSQS(mockClient)
			got, err := service.Consumer(ctx, mockDto)

			if tt.output == nil {
				if err != nil {
					t.Fatalf("expected nil error got %v", err)
				}
				if tt.want != nil {
					if string(got.GetMessage()) != string(tt.want.GetMessage()) {
						t.Fatalf("message body: got %q want %q", got.GetMessage(), tt.want.GetMessage())
					}
					if got.GetMessageId() != tt.want.GetMessageId() {
						t.Fatalf("message id: got %q want %q", got.GetMessageId(), tt.want.GetMessageId())
					}
					if !got.GetCreatedAt().Equal(tt.want.GetCreatedAt()) {
						t.Fatalf("created at: got %v want %v", got.GetCreatedAt(), tt.want.GetCreatedAt())
					}
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
