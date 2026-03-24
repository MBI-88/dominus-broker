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

		msg := entities.NewMessage(message)
		msg.SetMessage(message)
		rmessage := msg.GetMessage()

		if string(message) != string(rmessage) {
			t.Fatalf("%s != %s\n", message, rmessage)
		}

	}) 

	t.Run("CreatedAt field", func(t *testing.T){
		msg := entities.NewMessage(message)
		now := time.Now()
		msg.SetCreateAt(now)
		rnow := msg.GetCreatedAt() 

		if !now.Equal(rnow) {
			t.Fatalf("%s != %s\n", now, rnow)
		}

	})

	t.Run("MessageId field", func(t *testing.T) {
		msg := entities.NewMessage(message)
		id := uuid.New()
		msg.SetMessageId(id.String())
		rid := msg.GetMessageId()

		if id.String() != rid {
			t.Fatalf("%s != %s\n", id, rid)
		}

	})


}