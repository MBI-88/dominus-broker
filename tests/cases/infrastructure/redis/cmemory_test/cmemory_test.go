package cmemory_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/domain/repositories"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"
	"dominus-broker/internal/infrastructure/redis/cmemory"
	"dominus-broker/mocks"

	"github.com/alicebob/miniredis/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

func newMemoryClient(t *testing.T, s *miniredis.Miniredis, streamID string, lg event.Event) repositories.MemoryClient {
	t.Helper()
	host, portStr, err := net.SplitHostPort(s.Addr())
	if err != nil {
		t.Fatalf("SplitHostPort: %v", err)
	}
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		t.Fatalf("ParseInt port: %v", err)
	}
	return cmemory.NewMemoryClient(
		10,
		0,
		int(time.Second),
		int(time.Second),
		int(time.Second),
		port,
		0,
		host,
		"",
		false,
		"",
		streamID,
		lg,
	)
}

func testRedisClient(t *testing.T, s *miniredis.Miniredis) *redis.Client {
	t.Helper()
	return redis.NewClient(&redis.Options{
		Addr:         s.Addr(),
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
}

func payloadFromStreamEntry(values []string) (string, bool) {
	for i := 0; i+1 < len(values); i += 2 {
		if values[i] == enum.PAYLOAD {
			return values[i+1], true
		}
	}
	return "", false
}

func TestSendMessage(t *testing.T) {
	const streamID = "stream-send-test"

	t.Run("SendMessage ok", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		client := newMemoryClient(t, s, streamID, lgsMock)
		ctx := context.Background()

		lgsMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		msg := entities.NewMessageWithID([]byte("hello-redis"))
		// Redis stream IDs deben ser <ms>-<seq>, no UUID.
		msg.SetMessageId("1-0")

		if err := client.SendMessage(ctx, msg); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}

		entries, err := s.Stream(streamID)
		if err != nil {
			t.Fatalf("Stream: %v", err)
		}
		if len(entries) != 1 {
			t.Fatalf("expected 1 stream entry, got %d", len(entries))
		}
		if entries[0].ID != "1-0" {
			t.Fatalf("expected id 1-0, got %q", entries[0].ID)
		}
		raw, ok := payloadFromStreamEntry(entries[0].Values)
		if !ok {
			t.Fatalf("missing %q field in stream values: %v", enum.PAYLOAD, entries[0].Values)
		}
		var got entities.Message
		if err := jsoniter.Unmarshal([]byte(raw), &got); err != nil {
			t.Fatalf("Unmarshal payload: %v", err)
		}
		if string(got.GetMessage()) != "hello-redis" {
			t.Fatalf("message body: got %q", got.GetMessage())
		}
		if got.GetMessageId() != "1-0" {
			t.Fatalf("message id: got %q", got.GetMessageId())
		}
	})

	t.Run("SendMessage error on redis failure", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		client := newMemoryClient(t, s, streamID, lgsMock)
		ctx := context.Background()

		lgsMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		msg := entities.NewMessageWithID([]byte("x"))
		msg.SetMessageId("2-0")

		s.SetError("simulated XADD failure")
		if err := client.SendMessage(ctx, msg); err == nil {
			t.Fatal("expected error when Redis returns failure")
		}
	})
}

func TestAckMessage(t *testing.T) {
	const streamID = "stream-ack-test"
	const groupID = "ack-group"

	t.Run("AckMessage ok", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		lgsMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		client := newMemoryClient(t, s, streamID, lgsMock)
		ctx := context.Background()

		if err := client.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
		}

		msg := entities.NewMessageWithID([]byte("ack-body"))
		msg.SetMessageId("10-0")
		if err := client.SendMessage(ctx, msg); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}

		rdb := testRedisClient(t, s)
		t.Cleanup(func() { _ = rdb.Close() })

		read, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupID,
			Consumer: "worker-1",
			Streams:  []string{streamID, ">"},
			Count:    1,
			Block:    time.Second,
		}).Result()
		if err != nil {
			t.Fatalf("XReadGroup: %v", err)
		}
		if len(read) != 1 || len(read[0].Messages) != 1 {
			t.Fatalf("expected 1 message in XReadGroup, got %+v", read)
		}
		mid := read[0].Messages[0].ID

		if err := client.AckMessage(ctx, mid, groupID); err != nil {
			t.Fatalf("AckMessage: %v", err)
		}

		pending, err := rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
			Stream: streamID,
			Group:  groupID,
			Start:  "-",
			End:    "+",
			Count:  10,
		}).Result()
		if err != nil {
			t.Fatalf("XPendingExt: %v", err)
		}
		for _, p := range pending {
			if p.ID == mid {
				t.Fatalf("message %q should be acked, still pending", mid)
			}
		}
	})

	t.Run("AckMessage error on redis failure", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		lgsMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		client := newMemoryClient(t, s, streamID, lgsMock)
		ctx := context.Background()

		if err := client.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
		}
		msg := entities.NewMessageWithID([]byte("x"))
		msg.SetMessageId("20-0")
		if err := client.SendMessage(ctx, msg); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}

		rdb := testRedisClient(t, s)
		t.Cleanup(func() { _ = rdb.Close() })
		if _, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupID,
			Consumer: "worker-2",
			Streams:  []string{streamID, ">"},
			Count:    1,
			Block:    time.Second,
		}).Result(); err != nil {
			t.Fatalf("XReadGroup: %v", err)
		}

		s.SetError("simulated XACK failure")
		if err := client.AckMessage(ctx, "20-0", groupID); err == nil {
			t.Fatal("expected error when Redis returns failure")
		}
	})
}

func TestGetMessage(t *testing.T) {
	const streamID = "stream-get-test"
	const groupID = "get-group"

	t.Run("GetMessage ok", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		lgsMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		client := newMemoryClient(t, s, streamID, lgsMock)
		ctx := context.Background()

		if err := client.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
		}

		sent := entities.NewMessageWithID([]byte("get-payload"))
		sent.SetMessageId("30-0")
		if err := client.SendMessage(ctx, sent); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}

		got, err := client.GetMessage(ctx, "worker-get-1", groupID)
		if err != nil {
			t.Fatalf("GetMessage: %v", err)
		}
		if string(got.GetMessage()) != "get-payload" {
			t.Fatalf("body: got %q", got.GetMessage())
		}
		if got.GetMessageId() != "30-0" {
			t.Fatalf("id: got %q", got.GetMessageId())
		}
	})

	t.Run("GetMessage error on redis failure", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		lgsMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		client := newMemoryClient(t, s, streamID, lgsMock)
		ctx := context.Background()

		if err := client.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
		}
		msg := entities.NewMessageWithID([]byte("y"))
		msg.SetMessageId("31-0")
		if err := client.SendMessage(ctx, msg); err != nil {
			t.Fatalf("SendMessage: %v", err)
		}

		s.SetError("simulated XREADGROUP failure")
		if _, err := client.GetMessage(ctx, "worker-get-2", groupID); err == nil {
			t.Fatal("expected error when Redis returns failure")
		}
	})

	t.Run("GetMessage error on invalid payload JSON", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		lgsMock.EXPECT().
			WriteLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		client := newMemoryClient(t, s, streamID, lgsMock)
		ctx := context.Background()

		if err := client.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
		}
		rdb := testRedisClient(t, s)
		t.Cleanup(func() { _ = rdb.Close() })
		if _, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamID,
			ID:     "99-0",
			Values: map[string]any{enum.PAYLOAD: "{not-valid-json"},
		}).Result(); err != nil {
			t.Fatalf("XAdd: %v", err)
		}

		if _, err := client.GetMessage(ctx, "worker-bad-json", groupID); err == nil {
			t.Fatal("expected unmarshal error")
		}
	})

}

func TestGroupMessage(t *testing.T) {
	const streamID = "stream-group-test"

	t.Run("Group ok", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		lgsMock := mocks.NewMockEvent(ctrl)
		client := newMemoryClient(t, s, streamID, lgsMock)
		// Group no usa el logger en outbound.

		const groupID = "new-group"
		if err := client.Group(groupID); err != nil {
			t.Fatalf("Group: %v", err)
		}

		rdb := testRedisClient(t, s)
		t.Cleanup(func() { _ = rdb.Close() })
		groups, err := rdb.XInfoGroups(context.Background(), streamID).Result()
		if err != nil {
			t.Fatalf("XInfoGroups: %v", err)
		}
		var found bool
		for _, g := range groups {
			if g.Name == groupID {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected group %q in stream %q, got %#v", groupID, streamID, groups)
		}
	})

	t.Run("Group error on duplicate", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		client := newMemoryClient(t, s, streamID, mocks.NewMockEvent(ctrl))

		const groupID = "dup-group"
		if err := client.Group(groupID); err != nil {
			t.Fatalf("first Group: %v", err)
		}
		if err := client.Group(groupID); err == nil {
			t.Fatal("expected error creating the same consumer group twice")
		}
	})

	t.Run("Group error on redis failure", func(t *testing.T) {
		s := miniredis.RunT(t)
		ctrl := gomock.NewController(t)
		client := newMemoryClient(t, s, streamID, mocks.NewMockEvent(ctrl))

		s.SetError("simulated XGROUP failure")
		if err := client.Group("fail-group"); err == nil {
			t.Fatal("expected error when Redis returns failure")
		}
	})
}
