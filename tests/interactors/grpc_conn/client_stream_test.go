package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/tests/env"
	"encoding/json"
	"testing"
)

func TestClientStream(t *testing.T) {
	log := env.NewEventMock(true)
	topics := entities.NewTopics(100)
	stream := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, true), topics)

	t.Run("ClientStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(env.Data)
		result := stream.StreamClientConn(env.NewClienContextMock(payload, env.SubsOk))
		if result.Error() != "End" {
			t.Fatal(result)
		}
	})

	t.Run("ClientStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(env.Data)
		result := stream.StreamClientConn(env.NewClienContextMock(payload, make([]string, 0)))
		if result == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}
