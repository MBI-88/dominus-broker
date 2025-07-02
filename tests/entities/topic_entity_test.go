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

	t.Run("SetMessage_Ok", func(t *testing.T) {
		topic := &entities.Topic{
			Name:        env.Tps,
			Subscribers: env.SubsOk,
			Lck:         new(sync.Mutex),
		}

		topic.SetMessage(body)

	})

	t.Run("GetMessage_Ok", func(t *testing.T) {
		topic := &entities.Topic{
			Name:        env.Tps,
			Subscribers: env.SubsOk,
			
			Lck:         new(sync.Mutex),
			
		}
		topic.SetMessage(body)

		if body := topic.GetMessage(); len(body) == 0 {
			t.Fatal("Body must be different from 0")
		}
	})

	t.Run("GetMessage_error", func(t *testing.T) {
		topic := &entities.Topic{
			Name:        env.Tps,
			Subscribers: env.SubsOk,
			Lck:         new(sync.Mutex),
			
		}
		body := topic.GetMessage()

		if len(body) != 0 {
			t.Fatal("Body must be equal to 0")
		}
	})

}

func TestTopiParallelRW(t *testing.T) {
	topic := &entities.Topic{
		Name:        env.Tps,
		Subscribers: env.SubsOk,
		Lck:         new(sync.Mutex),
	}
	body, _ := json.Marshal(env.Data)

	t.Run("SetMessage_Ok", func(t *testing.T) {
		var (
			wg   = new(sync.WaitGroup)
		)
		wg.Add(2)

		go func() {
			defer wg.Done()
			topic.SetMessage(body)
		}()
		go func() {
			defer wg.Done()
			topic.SetMessage(body)
		}()

		wg.Wait()
	})

	t.Run("GetMessage_Ok", func(t *testing.T) {
		var (
			wg   = new(sync.WaitGroup)
			err1 error
			err2 error
		)
		wg.Add(2)

		go func() {
			defer wg.Done()
			if body := topic.GetMessage(); len(body) == 0 {
				err1 = errors.New("Body is empty")
			}
		}()

		go func() {
			defer wg.Done()
			if body := topic.GetMessage(); len(body) == 0 {
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
