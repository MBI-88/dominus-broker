package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/tests/env"
	"encoding/json"
	"testing"
)

func TestBidirectionalStream(t *testing.T) {
	log := env.NewEventMock(true)
	topics := entities.NewTopics(100)
	stream := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, true), topics)
	t.Run("BidirectionalStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(env.Data)
		result := stream.StreamBiConn(env.NewBiContextMock(payload, env.SubsOk))
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}
	})

	t.Run("BidirectionalStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(env.Data)
		result := stream.StreamBiConn(env.NewBiContextMock(payload, make([]string, 0)))
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})
}
