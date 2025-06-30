package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/tests/mocks"
	"encoding/json"
	"testing"
)

func TestBidirectionalStream(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	stream := grpcconn.NewGrpcService(log, mocks.NewGrpcClientMock(), topics)
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
