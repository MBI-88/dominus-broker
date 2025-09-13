package entities_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/tests/env"
	"encoding/json"
	"testing"
)

func TestQueue(t *testing.T) {

	t.Run("Enqueue_OK", func(t *testing.T) {
		queue := entities.NewQueue()
		message, _ := json.Marshal(env.Data)
		queue.Enqueue(message)

		if queue.Length() == 0 {
			t.Fatalf("Queue is empty")
		}
	})

	t.Run("Dequeue_OK", func(t *testing.T) {
		queue := entities.NewQueue()
		message, _ := json.Marshal(env.Data)
		queue.Enqueue(message)
		queue.Enqueue(message)

		_, ok := queue.Dequeue()
		if !ok {
			t.Fatal("Dequeue error")
		}

	})

	t.Run("Dequeue_Error", func(t *testing.T) {
		queue := entities.NewQueue()
		_, ok := queue.Dequeue()
		if ok {
			t.Fatal("Dequeue error")
		}
	})
}
