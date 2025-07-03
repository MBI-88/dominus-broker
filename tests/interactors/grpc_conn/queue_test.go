package grpcconn_test

import (
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/tests/env"
	"encoding/json"
	"testing"
	"time"
)

func TestSimpleConn(t *testing.T) {
	log := env.NewEventMock(true)
	topics := entities.NewTopics(100)
	topic := entities.NewTopic(nil)
	topic.SetName(env.Tps)
	topic.SetSubscribers(env.SubsOk)
	topics.Append(topic)
	simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, true), topics)

	t.Run("Connection-OK", func(t *testing.T) {
		payload, _ := json.Marshal(env.Data)
		simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, true), topics)
		if err := simple.SimpleConn(env.NewSimpleContextMock(payload, []string{env.Tps})); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Connection-Error", func(t *testing.T) {

		if err := simple.SimpleConn(env.NewSimpleContextMock(nil, make([]string, 0))); err == nil {
			t.Fatalf("Expected error but received nil")
		}
	})

	t.Run("Connection-Error_2", func(t *testing.T) {
		topic.SetSubscribers(make([]string, 0))
		if err := simple.SimpleConn(env.NewSimpleContextMock(nil, []string{env.Tps})); err == nil {
			t.Fatalf("Expected error but received nil")
		}
	})
}

func TestRunQueue(t *testing.T) {
	log := env.NewEventMock(true)
	topics := entities.NewTopics(100)
	topic := entities.NewTopic(nil)
	topic.SetName(env.Tps)
	topic.SetSubscribers(env.SubsOk)
	body, _ := json.Marshal(env.Data)
	topic.SetMessage(body)
	topics.Append(topic)

	t.Run("Find-Topic-OK", func(t *testing.T) {
		simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, false), topics)
		queue := make(chan struct{})
		go func() {
			time.Sleep(3 * time.Second)
			queue <- struct{}{}
		}()

		if err := simple.RunQueue(queue); err.Error() != "Queue closed" {
			t.Fatal("Queue error")
		}

		close(queue)
	})

	t.Run("Find-Topic-Connection-Error", func(t *testing.T) {
		simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(false, false), topics)
		queue := make(chan struct{})
		go func() {
			time.Sleep(3 * time.Second)
			queue <- struct{}{}
		}()

		if err := simple.RunQueue(queue); err.Error() != "Queue closed" {
			t.Fatal("Queue error")
		}

		close(queue)
	})
	
	t.Run("Find-Topic-Response-Error", func(t *testing.T) {
		simple := grpcconn.NewGrpcService(log, env.NewGrpcClientMock(true, false), topics)
		queue := make(chan struct{})
		go func() {
			time.Sleep(3 * time.Second)
			queue <- struct{}{}
		}()

		if err := simple.RunQueue(queue); err.Error() != "Queue closed" {
			t.Fatal("Queue error")
		}

		close(queue)
	})

}
