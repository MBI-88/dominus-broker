package system_test

import (
	"dominus-project/internal/interactors/system"
	"dominus-project/mocks"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestLogs(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockLogs)
		output    error
	}{
		{
			name: "GetLogs Ok",
			setupMock: func(ml *mocks.MockLogs) {
				ml.EXPECT().
					GetLogs().
					Return([]string{"log-1", "log-2"}, nil).Times(1)
			},
			output: nil,
		},
		{
			name: "GetLogs error",
			setupMock: func(ml *mocks.MockLogs) {
				ml.EXPECT().
					GetLogs().
					Return([]string{}, fmt.Errorf("logs error")).Times(1)
			},
			output: fmt.Errorf("logs error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLogs := mocks.NewMockLogs(ctrl)
			tt.setupMock(mockLogs)

			service := system.NewSystemService(mockLogs)

			_, err := service.GetLogs()
			if tt.output != nil && err == nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
			if tt.output == nil && err != nil {
				t.Fatalf("Expected output different from output %s != %s", tt.output, err)
			}
		})
	}
}
