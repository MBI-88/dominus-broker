package entities_test

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/tests/env"
	"sync"
	"testing"
)

func TestTopics(t *testing.T) {
	topic := &entities.Topic{
		Name:        env.Tps,
		Subscribers: env.SubsOk,
		Lck:         new(sync.Mutex),
	}

	t.Run("Append_Ok", func(t *testing.T) {
		topics := entities.NewTopics(1)
		if err := topics.Append(topic); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Append_error", func(t *testing.T) {
		topics := entities.NewTopics(1)
		topics.Append(topic)

		if err := topics.Append(topic); err == nil {
			t.Fatal("Error must be different fron nil")
		}
	})

	t.Run("Update_Ok", func(t *testing.T) {
		topics := entities.NewTopics(1)
		topics.Append(topic)

		if err := topics.Update(&entities.Topic{
			Name:        env.Tps,
			Subscribers: []string{"http://localhost:80", "http://localhost:8081"},
			Lck:         new(sync.Mutex),
		}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Update_error", func(t *testing.T) {
		topics := entities.NewTopics(1)
		if err := topics.Update(topic); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	t.Run("Find_Ok", func(t *testing.T) {
		topics := entities.NewTopics(1)
		topics.Append(topic)

		if _, err := topics.Find(env.Tps); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Find_error", func(t *testing.T) {
		topics := entities.NewTopics(1)

		if _, err := topics.Find(env.Tps); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})

	t.Run("Next_Ok", func(t *testing.T) {
		topics := entities.NewTopics(2)
		topics.Append(topic)
		topics.Append(topic)

		tp, err := topics.Next()
		if err != nil {
			t.Fatal(err)
		}
		if tp.Name != env.Tps {
			t.Fatal("Topic name are differents")
		}
	})

	t.Run("Next_error", func(t *testing.T) {
		topics := entities.NewTopics(1)
		topics.Append(topic)

		topics.Next()

		if _, err := topics.Next(); err == nil {
			t.Fatal("Error must be defferent from nil")
		}

	})

	t.Run("Delete_Ok", func(t *testing.T) {
		topics := entities.NewTopics(1)
		topics.Append(topic)

		if err := topics.Delete(topic); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Delete_error", func(t *testing.T) {
		topics := entities.NewTopics(1)

		if err := topics.Delete(topic); err == nil {
			t.Fatal("Error must be different from nil")
		}
	})
}
