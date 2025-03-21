package tests

import (
	"dominus-project/app/interactors"
	"encoding/json"
	"testing"
)

func TestSimpleConn(t *testing.T) {
	log := NewEventMock(true)
	inter := interactors.NewInteractor(log)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("Connection-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		simple := inter.NewGrpcService()
		result := simple.SimpleConn(NewSimpleConnMock(payload, subsOk))
		if result != nil {
			t.Fatal(result)
		}
	})

	t.Run("Connection-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		simple := inter.NewGrpcService()
		result := simple.SimpleConn(NewSimpleConnMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

func TestClientStream(t *testing.T) {
	log := NewEventMock(true)
	inter := interactors.NewInteractor(log)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("ClientStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamClientConn(NewClienConnMock(payload, subsOk))
		if result.Error() != "End" {
			t.Fatal(result)
		}
	})

	t.Run("ClientStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamClientConn(NewClienConnMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

func TestServerStream(t *testing.T) {
	log := NewEventMock(true)
	inter := interactors.NewInteractor(log)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("ServerStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		msg, ctx := NewServerConnMock(payload, subsOk)
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}

	})

	t.Run("ServerStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		msg, ctx := NewServerConnMock(payload, make([]string, 0))
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})

}

func TestBidirectionalStream(t *testing.T) {
	log := NewEventMock(true)
	inter := interactors.NewInteractor(log)
	inter = inter.Set(NewGrpcClientMock())
	t.Run("BidirectionalStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamBiConn(NewBiConnMock(payload, subsOk))
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}
	})

	t.Run("BidirectionalStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(data)
		stream := inter.NewGrpcService()
		result := stream.StreamBiConn(NewBiConnMock(payload, make([]string, 0)))
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})
}
