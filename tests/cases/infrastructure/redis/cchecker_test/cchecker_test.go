package cchecker_test

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"

	"dominus-broker/internal/domain/repositories"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/redis/cchecker"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newCheckerClient(t *testing.T, s *miniredis.Miniredis) repositories.CheckerClient {
	t.Helper()
	host, portStr, err := net.SplitHostPort(s.Addr())
	if err != nil {
		t.Fatalf("SplitHostPort: %v", err)
	}
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		t.Fatalf("ParseInt port: %v", err)
	}
	return cchecker.NewCheckerClient(
		port,
		0,
		host,
		"",
		false,
		"",
		3600,
	)
}

func TestSaveConsumer(t *testing.T) {
	s := miniredis.RunT(t)

	t.Run("SaveConsumer ok", func(t *testing.T) {
		checker := newCheckerClient(t, s)
		ctx := context.Background()
		const key = "consumer-ok"

		if err := checker.SaveConsumer(ctx, key); err != nil {
			t.Fatalf("SaveConsumer: %v", err)
		}

		redisKey := enum.ID_POTENCY_TAG + ":" + key
		if !s.DB(0).Exists(redisKey) {
			t.Fatalf("miniredis: key %q should exist", redisKey)
		}
	})

	t.Run("SaveConsumer error", func(t *testing.T) {
		checker := newCheckerClient(t, s)
		ctx := context.Background()
		const key = "consumer-dup"
		if err := checker.SaveConsumer(ctx, key); err != nil {
			t.Fatalf("first SaveConsumer: %v", err)
		}
		err := checker.SaveConsumer(ctx, key)
		if err == nil {
			t.Fatal("second SaveConsumer: expected error (SET NX)")
		}
		// SaveConsumer wraps with fmt.Errorf("... %s", err) (no %w), so errors.Is(redis.Nil) does not apply.
		if !strings.Contains(err.Error(), redis.Nil.Error()) {
			t.Fatalf("expected wrapped redis nil error, got: %v", err)
		}
	})
}

func TestConsumer(t *testing.T) {
	s := miniredis.RunT(t)

	t.Run("CheckConsumer ok when key exists", func(t *testing.T) {
		checker := newCheckerClient(t, s)
		ctx := context.Background()
		const key = "check-ok"
		if err := checker.SaveConsumer(ctx, key); err != nil {
			t.Fatalf("SaveConsumer: %v", err)
		}
		if !checker.CheckConsumer(ctx, key) {
			t.Fatal("CheckConsumer: expected true for existing key")
		}
	})

	t.Run("CheckConsumer false when key missing", func(t *testing.T) {
		checker := newCheckerClient(t, s)
		ctx := context.Background()
		if checker.CheckConsumer(ctx, "never-saved-key") {
			t.Fatal("CheckConsumer: expected false when key was never saved")
		}
	})

	t.Run("CheckConsumer false on redis error", func(t *testing.T) {
		checker := newCheckerClient(t, s)
		ctx := context.Background()
		s.SetError("simulated failure")
		if checker.CheckConsumer(ctx, "any-key") {
			t.Fatal("CheckConsumer: expected false when Redis returns error")
		}
	})
}
