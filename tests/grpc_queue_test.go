package tests

import (
	"dominus-project/app/interactors"
	"encoding/json"
	"testing"
)

func TestSimpleConn(t *testing.T) {
	log := NewEventMock(true)
	inter := interactors.NewInteractor(log, 100)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("Connection-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		simple := inter.NewGrpcService()
		result := simple.SimpleConn(NewSimpleContextMock(payload, subsOk))
		if result != nil {
			t.Fatal(result)
		}
	})

	t.Run("Connection-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		simple := inter.NewGrpcService()
		result := simple.SimpleConn(NewSimpleContextMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}


