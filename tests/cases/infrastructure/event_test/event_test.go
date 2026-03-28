package event_test

import (
	"bytes"
	"context"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/event"
	"io"
	"os"
	"testing"
)

func TestWriteLog(t *testing.T) {

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
				t.Fatal(err)
			}
		}()

		lgs := event.NewEvent("", "", false)
		ctx := context.WithValue(context.Background(), enum.ID, "test-1123")
		var buffer bytes.Buffer
		lgs.WriteLog(ctx, enum.INFO, "WriteLog", enum.DEBUG_DESCRIPTION)

		if err := w.Close(); err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(&buffer, r); err != nil {
			t.Fatal(err)
		}

		if buffer.String() == "" {
			t.Fatal("Expected write the console")
		}

	})

	t.Run("WriteLog Warn", func(t *testing.T) {
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = w

		defer func() {
			os.Stdout = old
			if err := r.Close(); err != nil {
				t.Fatal(err)
			}
		}()

		lgs := event.NewEvent("", "", false)
		ctx := context.WithValue(context.Background(), enum.ID, "test-1123")
		var buffer bytes.Buffer
		lgs.WriteLog(ctx, enum.WARN, "WriteLog", enum.DEBUG_DESCRIPTION)

		if err := w.Close(); err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(&buffer, r); err != nil {
			t.Fatal(err)
		}

		if buffer.String() == "" {
			t.Fatal("Expected write the console")
		}
	})

	t.Run("WriteLog Error", func(t *testing.T) {
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = w

		defer func() {
			os.Stdout = old
			if err := r.Close(); err != nil {
				t.Fatal(err)
			}
		}()

		lgs := event.NewEvent("", "", false)
		ctx := context.WithValue(context.Background(), enum.ID, "test-1123")
		var buffer bytes.Buffer
		lgs.WriteLog(ctx, enum.ERROR, "WriteLog", enum.DEBUG_DESCRIPTION)

		if err := w.Close(); err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(&buffer, r); err != nil {
			t.Fatal(err)
		}

		if buffer.String() == "" {
			t.Fatal("Expected write the console")
		}

	})

}

func TestCheckID(t *testing.T) {

	t.Run("CheckID with id", func(t *testing.T) {
		lgs := event.NewEvent("", "", false)
		ctx := context.WithValue(context.Background(), enum.ID, "test-1123")
		ctx = lgs.CheckID(ctx)

		if ctx.Value(enum.ID) != "test-1123" {
			t.Fatalf("Expected %s, got %s", "test-1123", ctx.Value(enum.ID))
		}
	})

	t.Run("CheckID without id", func(t *testing.T) {
		lgs := event.NewEvent("", "", false)
		ctx := context.Background()
		ctx = lgs.CheckID(ctx)
		if ctx.Value(enum.ID) == nil {
			t.Fatal("Expected check the id")
		}

	})
}
