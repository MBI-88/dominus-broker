package tests

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/interactors"
	"encoding/json"
	"testing"
)

func TestClientStream(t *testing.T) {
	log := NewEventMock(true)
	topics := entities.NewTopics(100)
	inter := interactors.NewInteractor(log, topics, 100)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("ClientStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamClientConn(NewClienContextMock(payload, subsOk))
		if result.Error() != "End" {
			t.Fatal(result)
		}
	})

	t.Run("ClientStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamClientConn(NewClienContextMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

func TestServerStream(t *testing.T) {
	log := NewEventMock(true)
	topics := entities.NewTopics(100)
	inter := interactors.NewInteractor(log, topics, 100)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("ServerStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		msg, ctx := NewServerContextMock(payload, subsOk)
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}

	})

	t.Run("ServerStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		msg, ctx := NewServerContextMock(payload, make([]string, 0))
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})

}

func TestBidirectionalStream(t *testing.T) {
	log := NewEventMock(true)
	topics := entities.NewTopics(100)
	inter := interactors.NewInteractor(log, topics, 100)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("BidirectionalStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamBiConn(NewBiContextMock(payload, subsOk))
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}
	})

	t.Run("BidirectionalStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamBiConn(NewBiContextMock(payload, make([]string, 0)))
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})
}