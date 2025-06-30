package system_test

import (
	"dominus-project/internal/interactors/system"
	"dominus-project/tests/mocks"
	"testing"
)

func TestLogs(t *testing.T) {

	t.Run("Logs_ok", func(t *testing.T) {
		log := mocks.NewEventMock(true)
		system := system.NewSystemService(log)

		result, err := system.GetLogs()
		if err != nil {
			t.Fatal(err)
		}

		if len(result) == 0 {
			t.Fatalf("[-] Result must be different of empty")
		}

	})

	t.Run("Logs_error", func(t *testing.T) {
		log := mocks.NewEventMock(false)
		system := system.NewSystemService(log)

		_, err := system.GetLogs()
		if err.Error() != "[-] Error response" {
			t.Fatalf("[-] Error is different")
		}

	})
}
