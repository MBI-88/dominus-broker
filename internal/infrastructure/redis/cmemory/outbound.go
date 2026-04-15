package cmemory

import (
	"context"
	"crypto/tls"
	"dominus-broker/internal/domain/entities"
	"dominus-broker/internal/domain/repositories"
	"dominus-broker/internal/infrastructure/enum"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

type memory struct {
	rdb      *redis.Client
	streamID string
}

func NewMemoryClient(
	Port int64,
	Db int,
	Host string,
	Password string,
	Tls bool,
	Username string,
	StreamID string,
) repositories.MemoryClient {

	var cfTls *tls.Config

	if Tls {
		cfTls = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(&redis.Options{
		DB:        Db,
		Addr:      fmt.Sprintf("%s:%d", Host, Port),
		Password:  Password,
		TLSConfig: cfTls,
		Username:  Username,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return &memory{
		rdb:      client,
		streamID: StreamID,
	}
}

func (m *memory) SendMessage(ctx context.Context, q *entities.Message) error {
	data, err := jsoniter.Marshal(&MessageDto{
		Message:   q.GetMessage(),
		MessageId: q.GetMessageId(),
		CreatedAt: q.GetCreatedAt(),
	})
	if err != nil {
		return err
	}
	if _, err := m.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: m.streamID,
		Values: map[string]any{
			enum.PAYLOAD: data,
		},
		ID: q.GetMessageId(),
	}).Result(); err != nil {
		return fmt.Errorf("memory.SendMessage %s", err)
	}
	return nil
}

func (m *memory) AckMessage(ctx context.Context, messageId, groupId string) error {
	if err := m.rdb.XAck(ctx, m.streamID, groupId, messageId).Err(); err != nil {
		return fmt.Errorf("memory.AckMessage %s", err)
	}
	return nil
}

func (m *memory) GetMessage(ctx context.Context, workerId, groupId string) (*entities.Message, error) {
	response, err := m.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupId,
		Consumer: workerId,
		Streams:  []string{m.streamID, ">"},
		Block:    100 * time.Millisecond,
		Count:    1,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("memory.GetMessage %s", err)
	}

	streamMsg := response[0].Messages[0]
	message := streamMsg.Values[enum.PAYLOAD].(string)

	var data MessageDto
	if err := jsoniter.Unmarshal([]byte(message), &data); err != nil {
		return nil, fmt.Errorf("memory.GetMessage %s", err)
	}

	entity := entities.NewMessage()
	if !entity.SetMessageId(streamMsg.ID) {
		return nil, fmt.Errorf("memory.invalid messageID format")
	}

	entity.SetCreateAt(data.CreatedAt)
	entity.SetMessage(data.Message)
	return entity, nil
}

func (m *memory) Group(groupId string) error {
	return m.rdb.XGroupCreateMkStream(context.Background(), m.streamID, groupId, enum.START_FROM_NEW_MESSAGE).Err()
}
