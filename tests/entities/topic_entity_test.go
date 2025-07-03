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
		topic := entities.NewTopic(nil)
		topic.SetMessage(body)

	})

	t.Run("GetMessage_Ok", func(t *testing.T) {
		topic := entities.NewTopic(nil)
		topic.SetMessage(body)

		if body := topic.GetMessage(); len(body) == 0 {
			t.Fatal("Body must be different from 0")
		}
	})

	t.Run("GetMessage_error", func(t *testing.T) {
		topic := entities.NewTopic(nil)
		body := topic.GetMessage()

		if len(body) != 0 {
			t.Fatal("Body must be equal to 0")
		}
	})

}

func TestTopicParallelRW(t *testing.T) {
	topic := entities.NewTopic(nil)
	body, _ := json.Marshal(env.Data)

	t.Run("SetMessage_GetMessage_RW", func(t *testing.T) {
		var (
			wg   = new(sync.WaitGroup) 
			err error
		)
		wg.Add(2)

		go func() {
			defer wg.Done()
			topic.SetMessage(body)
		}()
		go func() {
			defer wg.Done()
			time.Sleep(2 * time.Second)
			if body := topic.GetMessage(); len(body) == 0 {
				err = errors.New("Body is empty")
			}
		}()

		wg.Wait()

		if err != nil {
			t.Fatal(err)
		}
	})

}
