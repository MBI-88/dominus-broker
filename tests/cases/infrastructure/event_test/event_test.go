package event_test

import (
	"bytes"
	"context"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/event"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestWriteLog(t *testing.T) {
	lgs := event.NewEvent("", "", false)
	ctx := context.WithValue(context.Background(), enum.ID, "test-1123")
	var buffer bytes.Buffer

	t.Run("WriteLog info", func(t *testing.T) {
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = w

		defer func() {
			os.Stdout = old
			if err := r.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		lgs.WriteLog(ctx, enum.INFO, "WriteLog", enum.DEBUG_DESCRIPTION)

		if err := w.Close(); err != nil {
			fmt.Println(err)
		}

		if _, err := io.Copy(&buffer, r); err != nil {
			t.Fatal(err)
		}

		
		if buffer.String() == "" {
			t.Fatal("Expected write the console")
		}

	})

	t.Run("WriteLog Warn", func(t *testing.T) {
		lgs.WriteLog(ctx, enum.WARN, "WriteLog", enum.DEBUG_DESCRIPTION)

		
	})

	t.Run("WriteLog Error", func(t *testing.T) {
		lgs.WriteLog(ctx, enum.ERROR, "WriteLog", enum.DEBUG_DESCRIPTION)

		
	})

}

func TestCheckID(t *testing.T) {

}
