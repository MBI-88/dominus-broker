package sqs_test

import (
	"testing"

	"dominus-broker/internal/application/use_cases/sqs"
	"dominus-broker/mocks"

	"go.uber.org/mock/gomock"
)

func TestNewSQS(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	c := mocks.NewMockMemoryClient(ctrl)
	s := sqs.NewSQS(c)
	if s == nil {
		t.Fatal("NewSQS returned nil")
	}
}
