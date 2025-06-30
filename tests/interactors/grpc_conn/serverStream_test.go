package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/tests/mocks"
	"encoding/json"
	"testing"
)

func TestServerStream(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	stream := grpcconn.NewGrpcService(log, mocks.NewGrpcClientMock(), topics)
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
