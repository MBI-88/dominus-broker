package tests

import (
	"dominus-project/app/domain/entities"
	"encoding/json"
	"sync"
	"testing"
)

func TestTopic(t *testing.T) {
	body, _ := json.Marshal(data)

	t.Run("Push_Ok", func(t *testing.T) {
		topic := &entities.Topic{
			Name:       tps,
			Partitions: subsOk,
			Queue:      entities.NewQueue(),
			Lck:        new(sync.Mutex),
			Limit:      3,
		}

		if err := topic.Push(body); err != nil {
			t.Fatal(err)
		}

		if err := topic.Push(body); err != nil {
			t.Fatal(err)
		}

		if err := topic.Push(body); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Push_error", func(t *testing.T) {
		topic := &entities.Topic{
			Name:       tps,
			Partitions: subsOk,
			Queue:      entities.NewQueue(),
			Lck:        new(sync.Mutex),
			Limit:      1,
		}
		topic.Push(body)

		if err := topic.Push(body); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})

	t.Run("Pop_Ok", func(t *testing.T) {
		topic := &entities.Topic{
			Name:       tps,
			Partitions: subsOk,
			Queue:      entities.NewQueue(),
			Lck:        new(sync.Mutex),
			Limit:      1,
		}
		topic.Push(body)
		topic.Push(body)
		topic.Push(body)

		if body := topic.Pop(); len(body) == 0 {
			t.Fatal("Body must be different from 0")
		}
		if body := topic.Pop(); len(body) == 0 {
			t.Fatal("Body must be different from 0")
		}
		if body := topic.Pop(); len(body) == 0 {
			t.Fatal("Body must be different from 0")
		}
	})

	t.Run("Pop_error", func(t *testing.T) {
		topic := &entities.Topic{
			Name:       tps,
			Partitions: subsOk,
			Queue:      entities.NewQueue(),
			Lck:        new(sync.Mutex),
			Limit:      1,
		}
		body := topic.Pop()

		if len(body) != 0 {
			t.Fatal("Body must be equal to 0")
		}
	})

}


func TestTopiParallelRW(t *testing.T) {
	
}