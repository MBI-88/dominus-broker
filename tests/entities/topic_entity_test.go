package entities_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/tests/env"
	"encoding/json"
	"errors"
	"sync"
	"testing"
)

func TestTopic(t *testing.T) {
	body, _ := json.Marshal(env.Data)

	t.Run("Push_Ok", func(t *testing.T) {
		topic := &entities.Topic{
			Name:        env.Tps,
			Subscribers: env.SubsOk,
			Queue:       entities.NewQueue(),
			Lck:         new(sync.Mutex),
			Limit:       3,
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
			Name:        env.Tps,
			Subscribers: env.SubsOk,
			Queue:       entities.NewQueue(),
			Lck:         new(sync.Mutex),
			Limit:       1,
		}
		topic.Push(body)

		if err := topic.Push(body); err == nil {
			t.Fatal("Error must be different from nil")
		}

	})

	t.Run("Pop_Ok", func(t *testing.T) {
		topic := &entities.Topic{
			Name:        env.Tps,
			Subscribers: env.SubsOk,
			Queue:       entities.NewQueue(),
			Lck:         new(sync.Mutex),
			Limit:       1,
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
			Name:        env.Tps,
			Subscribers: env.SubsOk,
			Queue:       entities.NewQueue(),
			Lck:         new(sync.Mutex),
			Limit:       1,
		}
		body := topic.Pop()

		if len(body) != 0 {
			t.Fatal("Body must be equal to 0")
		}
	})

}

func TestTopiParallelRW(t *testing.T) {
	topic := &entities.Topic{
		Name:        env.Tps,
		Subscribers: env.SubsOk,
		Queue:       entities.NewQueue(),
		Lck:         new(sync.Mutex),
		Limit:       3,
	}
	body, _ := json.Marshal(env.Data)

	t.Run("Push_Ok", func(t *testing.T) {
		var (
			wg   = new(sync.WaitGroup)
			err1 error
			err2 error
		)
		wg.Add(2)

		go func() {
			defer wg.Done()
			if err := topic.Push(body); err != nil {
				err1 = err
			}
		}()
		go func() {
			defer wg.Done()
			if err := topic.Push(body); err != nil {
				err2 = err
			}
		}()

		wg.Wait()

		if err1 != nil {
			t.Fatal(err1)
		}
		if err2 != nil {
			t.Fatal(err2)
		}
	})

	t.Run("Pop_Ok", func(t *testing.T) {
		var (
			wg   = new(sync.WaitGroup)
			err1 error
			err2 error
		)
		wg.Add(2)

		go func() {
			defer wg.Done()
			if body := topic.Pop(); len(body) == 0 {
				err1 = errors.New("Body is empty")
			}
		}()

		go func() {
			defer wg.Done()
			if body := topic.Pop(); len(body) == 0 {
				err2 = errors.New("Body is empty")
			}
		}()

		wg.Wait()

		if err1 != nil {
			t.Fatal(err1)
		}

		if err2 != nil {
			t.Fatal(err2)
		}
	})

}
