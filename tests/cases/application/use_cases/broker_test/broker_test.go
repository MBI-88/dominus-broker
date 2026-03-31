package broker_test

import (
	"testing"

	"dominus-broker/internal/application/use_cases/broker"
	"dominus-broker/mocks"

	"go.uber.org/mock/gomock"
)

func TestNewBroker(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	c := mocks.NewMockBrokerClient(ctrl)
	b := broker.NewBroker(c)
	if b == nil {
		t.Fatal("NewBroker returned nil")
	}
}
