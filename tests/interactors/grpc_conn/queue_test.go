package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/tests/mocks"
	"encoding/json"
	"sync"
	"testing"
)

func TestSimpleConn(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	topic := &entities.Topic{
		Name:        mocks.Tps,
		Subscribers: mocks.SubsOk,
		Lck:         new(sync.Mutex),
		Queue:       entities.NewQueue(),
		Limit:       100,
	}
	topics.Append(topic)
	simple := grpcconn.NewGrpcService(log, mocks.NewGrpcClientMock(), topics)

	t.Run("Connection-OK", func(t *testing.T) {
		payload, _ := json.Marshal(mocks.Data)

		if err := simple.SimpleConn(mocks.NewSimpleContextMock(payload, []string{mocks.Tps})); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Connection-Error", func(t *testing.T) {

		if err := simple.SimpleConn(mocks.NewSimpleContextMock(nil, make([]string, 0))); err == nil {
			t.Fatalf("Expected error but received nil")
		}
	})

	t.Run("Connection-Error_2", func(t *testing.T) {
		topic.Subscribers = make([]string, 0)

		if err := simple.SimpleConn(mocks.NewSimpleContextMock(nil, []string{mocks.Tps})); err == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

// Probar RunQueue
