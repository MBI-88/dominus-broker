package grpcconn_test

import (
	"dominus-project/app/domain/entities"
	grpcconn "dominus-project/app/interactors/grpc_conn"
	"dominus-project/tests/mocks"
	"encoding/json"
	"testing"
)

func TestClientStream(t *testing.T) {
	log := mocks.NewEventMock(true)
	topics := entities.NewTopics(100)
	stream := grpcconn.NewGrpcService(log, mocks.NewGrpcClientMock(), topics)

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