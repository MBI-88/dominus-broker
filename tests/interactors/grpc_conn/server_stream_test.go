package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/tests/env"
	"encoding/json"
	"testing"
)

func TestServerStream(t *testing.T) {
	log := env.NewEventMock(true)
	topics := entities.NewTopics(100)
	stream := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true,true), topics)
	t.Run("ServerStream-OK", func(t *testing.T) {
		payload, _ := json.Marshal(env.Data)
		msg, ctx := env.NewServerContextMock(payload, env.SubsOk)
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() != "Connection closed" {
			t.Fatal(result)
		}

	})

	t.Run("ServerStream-Error", func(t *testing.T) {
		payload, _ := json.Marshal(env.Data)
		msg, ctx := env.NewServerContextMock(payload, make([]string, 0))
		result := stream.StreamServerConn(msg, ctx)
		if result.Error() == "End" {
			t.Fatalf("Expected error but received nil")
		}
	})

}
