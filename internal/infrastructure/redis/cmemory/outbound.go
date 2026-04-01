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
	PoolSize int,
	MaxRetries int,
	DialTimeOut int,
	ReadTimeOut int,
	WriteTimeOut int,
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
		PoolSize:     PoolSize,
		MaxRetries:   MaxRetries,
		DialTimeout:  time.Duration(DialTimeOut),
		ReadTimeout:  time.Duration(ReadTimeOut),
		WriteTimeout: time.Duration(WriteTimeOut),
		DB:           Db,
		Addr:         fmt.Sprintf("%s:%d", Host, Port),
		Password:     Password,
		TLSConfig:    cfTls,
		Username:     Username,
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
	data, err := jsoniter.Marshal(q)
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
		return fmt.Errorf("cmemory.SendMessage %s", err)
	}
	return nil
}

func (m *memory) AckMessage(ctx context.Context, messageId, groupId string) error {
	if err := m.rdb.XAck(ctx, m.streamID, groupId, messageId).Err(); err != nil {
		return fmt.Errorf("cmemory.AckMessage %s", err)
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
		return nil, fmt.Errorf("cmemory.GetMessage %s", err)
	}

	streamMsg := response[0].Messages[0]
	message := streamMsg.Values[enum.PAYLOAD].(string)

	var q entities.Message
	if err := jsoniter.Unmarshal([]byte(message), &q); err != nil {
		return nil, fmt.Errorf("cmemory.GetMessage %s", err)
	}

	if !q.SetMessageId(streamMsg.ID) {
		return nil, fmt.Errorf("invalid messageID format")
	}

	return &q, nil
}

func (m *memory) Group(groupId string) error {
	return m.rdb.XGroupCreateMkStream(context.Background(), m.streamID, groupId, enum.START_FROM_NEW_MESSAGE).Err()
}
