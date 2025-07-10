package entities_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/tests/env"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestTopic(t *testing.T) {
	body, _ := json.Marshal(env.Data)

	t.Run("SetMessage_Ok", func(t *testing.T) {
		topic := entities.NewTopic(nil, env.QueueLimit)
		if err := topic.SetMessage(body); err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("SetMessage_OK", func(t *testing.T) {
		topic := entities.NewTopic(nil, 1)
		if err := topic.SetMessage(body); err != nil {
			t.Fatal(err.Error())
		}

		if err := topic.SetMessage(body); err == nil {
			t.Fatal("SetMessage must return error")
		}
	})

	t.Run("GetMessage_Ok", func(t *testing.T) {
		topic := entities.NewTopic(nil, env.QueueLimit)

		if err := topic.SetMessage(body); err != nil {
			t.Fatal(err.Error())
		}

		if body := topic.GetMessage(); len(body) == 0 {
			t.Fatal("Body must be different from 0")
		}
	})

	t.Run("GetMessage_error", func(t *testing.T) {
		topic := entities.NewTopic(nil, env.QueueLimit)
		body := topic.GetMessage()

		if len(body) != 0 {
			t.Fatal("Body must be equal to 0")
		}
	})

}

func TestTopicParallelRW(t *testing.T) {
	topic := entities.NewTopic(nil, env.QueueLimit)
	body, _ := json.Marshal(env.Data)

	t.Run("SetMessage_GetMessage_RW", func(t *testing.T) {
		var (
			wg   = new(sync.WaitGroup) 
			errGet error
			errSet error
		)
		wg.Add(2)

		go func() {
			defer wg.Done()
			errSet = topic.SetMessage(body)
			errSet = topic.SetMessage(body)
			errSet = topic.SetMessage(body)
		}()
		go func() {
			defer wg.Done()
			time.Sleep(1 * time.Second)
			if body := topic.GetMessage(); len(body) == 0 {
				errGet = errors.New("Body is empty")
			}
		}()

		wg.Wait()

		if errGet != nil || errSet != nil {
			t.Fatal(errGet.Error() + " " + errSet.Error())
		}
	})

}
