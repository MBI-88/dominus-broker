package tests

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/interactors"
	"encoding/json"
	"sync"
	"testing"
)

func TestSimpleConn(t *testing.T) {
	log := NewEventMock(true)
	topics := entities.NewTopics(100)
	topic := &entities.Topic{
		Name: tps,
		Partitions: subsOk,
		Lck: new(sync.Mutex),
		Queue: entities.NewQueue(),
		Limit: 100,
	}
	topics.Append(topic)
	inter := interactors.NewInteractor(log, topics, 100)

	t.Run("Connection-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		simple := inter.NewGrpcService()
		
		if err := simple.SimpleConn(NewSimpleContextMock(payload, []string{tps})); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Connection-Error", func(t *testing.T) {
		simple := inter.NewGrpcService()
		
		if err := simple.SimpleConn(NewSimpleContextMock(nil, make([]string, 0))); err == nil {
			t.Fatalf("Expected error but received nil")
		}
	})

	t.Run("Connection-Error_2", func(t *testing.T) {
		simple := inter.NewGrpcService()
		topic.Partitions = make([]string, 0)

		if err := simple.SimpleConn(NewSimpleContextMock(nil, []string{tps})); err == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}


