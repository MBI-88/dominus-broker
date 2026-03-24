package sqs_test

import (
	"context"
	"dominus-project/internal/application/use_cases/sqs"
	"dominus-project/mocks"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAck(t *testing.T) {
	tests := []struct{
		name string 
		setupMock func (context.Context, *mocks.MockConsumerDto, *mocks.MockMemoryClient)
		output error
	}{	
		{
			name: "Ack OK",
			setupMock: func(ctx context.Context, dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient){
				dto.EXPECT(). 
				GetMessageId().Return("213fd31er3e1r3")
				dto.EXPECT(). 
				GetGroupId().Return("consumer-1")

				mc.EXPECT(). 
				AckMessage(ctx, "213fd31er3e1r3", "consumer-1"). 
				Return(nil)
			},
			output: nil,
		},
		{
			name: "Ack connection error",
			setupMock: func(ctx context.Context, dto *mocks.MockConsumerDto, mc *mocks.MockMemoryClient){
				dto.EXPECT(). 
				GetMessageId().Return("213fd31er3e1r3")
				dto.EXPECT(). 
				GetGroupId().Return("consumer-1")

				mc.EXPECT(). 
				AckMessage(ctx, "213fd31er3e1r3", "consumer-1"). 
				Return(fmt.Errorf("connection error"))
			},
			output: fmt.Errorf("connection error"),
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
			err := service.Ack(ctx, mockDto)

			if tt.output.Error() != err.Error() {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
		})
	}
}