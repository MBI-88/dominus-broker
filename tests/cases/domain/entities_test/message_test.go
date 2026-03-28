package entities_test

import (
	"dominus-project/internal/domain/entities"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMessage(t *testing.T) {

	message := []byte("Hola mundo test")
	t.Run("Message field", func(t *testing.T) {

		msg := entities.NewMessageWithID(message)
		msg.SetMessage(message)
		rmessage := msg.GetMessage()

		if string(message) != string(rmessage) {
			t.Fatalf("%s != %s\n", message, rmessage)
		}

	})

	t.Run("CreatedAt field", func(t *testing.T) {
		msg := entities.NewMessageWithID(message)
		now := time.Now()
		msg.SetCreateAt(now)
		rnow := msg.GetCreatedAt()

		if !now.Equal(rnow) {
			t.Fatalf("%s != %s\n", now, rnow)
		}

	})

	t.Run("MessageId field valid Redis stream id", func(t *testing.T) {
		msg := entities.NewMessageWithID(message)
		const want = "1700000000099-42"
		if !msg.SetMessageId(want) {
			t.Fatal("SetMessageId should accept <ms>-<seq> format")
		}
		if msg.GetMessageId() != want {
			t.Fatalf("got %q want %q", msg.GetMessageId(), want)
		}
	})

	t.Run("MessageId rejects UUID and keeps previous id", func(t *testing.T) {
		msg := entities.NewMessageWithID(message)
		before := msg.GetMessageId()
		id := uuid.New()
		if msg.SetMessageId(id.String()) {
			t.Fatal("SetMessageId should reject non stream-id strings")
		}
		if msg.GetMessageId() != before {
			t.Fatalf("id should be unchanged on reject: got %q want %q", msg.GetMessageId(), before)
		}
	})

}
