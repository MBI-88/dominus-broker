package interactors_test

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/interactors"
	"dominus-project/tests/mocks"
	"encoding/json"
	"testing"
)

func TestClientStream(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	inter := interactors.NewInteractor(log, topics, 100)
	inter = inter.Set(mocks.NewGrpcClientMock())
	stream := inter.NewGrpcService()
	t.Run("ClientStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(mocks.Data)
		result := stream.StreamClientConn(mocks.NewClienContextMock(payload, mocks.SubsOk))
		if result.Error() != "End" {
			t.Fatal(result)
		}
	})

	t.Run("ClientStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(mocks.Data)
		result := stream.StreamClientConn(mocks.NewClienContextMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

func TestServerStream(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	inter := interactors.NewInteractor(log, topics, 100)
	inter = inter.Set(mocks.NewGrpcClientMock())
	stream := inter.NewGrpcService()
	t.Run("ServerStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(mocks.Data)
		msg, ctx := mocks.NewServerContextMock(payload, mocks.SubsOk)
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}

	})

	t.Run("ServerStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(mocks.Data)
		msg, ctx := mocks.NewServerContextMock(payload, make([]string, 0))
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})

}

func TestBidirectionalStream(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	inter := interactors.NewInteractor(log, topics, 100)
	inter = inter.Set(mocks.NewGrpcClientMock())
	stream := inter.NewGrpcService()
	t.Run("BidirectionalStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(mocks.Data)
		result := stream.StreamBiConn(mocks.NewBiContextMock(payload, mocks.SubsOk))
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}
	})

	t.Run("BidirectionalStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(mocks.Data)
		result := stream.StreamBiConn(mocks.NewBiContextMock(payload, make([]string, 0)))
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})
}