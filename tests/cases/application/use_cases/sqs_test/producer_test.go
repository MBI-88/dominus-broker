package sqs_test

import (
	"context"
	"dominus-project/internal/application/use_cases/sqs"
	"dominus-project/mocks"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestProducer(t *testing.T) {
	tests := []struct{
		name string
		setupMock func ( *mocks.MockProducerDto, *mocks.MockMemoryClient)
		output error
	}{
		{
			name: "Producer OK",
			setupMock: func(dto *mocks.MockProducerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().
				GetPayload().
				Return([]byte("\n{payload:{data:data}}\n"))

				mc.EXPECT(). 
				SendMessage(gomock.All(), gomock.All()). 
				Return(nil)
			},
			output: nil,
		},
		{
			name: "Producer connection error",
			setupMock: func(dto *mocks.MockProducerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().
				GetPayload().
				Return([]byte("\n{payload:{data:data}}\n"))

				mc.EXPECT(). 
				SendMessage(gomock.All(), gomock.All()). 
				Return(fmt.Errorf("connection error"))
			},
			output: fmt.Errorf("connection error"),

		},
		{
			name: "Producer empty payload",
			setupMock: func(dto *mocks.MockProducerDto, mc *mocks.MockMemoryClient) {
				dto.EXPECT().
				GetPayload().
				Return([]byte{})
			},
			output: fmt.Errorf("empty payload") ,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockDto := mocks.NewMockProducerDto(ctrl)
			mockClient := mocks.NewMockMemoryClient(ctrl)

			tt.setupMock(mockDto, mockClient)
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