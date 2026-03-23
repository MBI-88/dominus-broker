package cmemory

import (
	"context"
	"crypto/tls"
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repositories"
	"dominus-project/internal/infrastructure/enum"
	"dominus-project/internal/infrastructure/event"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type memory struct {
	rdb      *redis.Client
	lg       event.Event
	streamID string
}

func NewMemoryClient(
	PoolSize int64,
	IdleConn int64,
	MaxRetries int64,
	DialTimeOut int64,
	ReadTimeOut int64,
	WriteTimeOut int64,
	Port int64,
	Db int,
	Host string,
	Password string,
	Tls bool,
	Username string,
	StreamID string,
	lg event.Event,
) repositories.MemoryClient {

	var cfTls *tls.Config

	if Tls {
		cfTls = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(&redis.Options{
		PoolSize:     int(PoolSize),
		IdleTimeout:  time.Duration(IdleConn),
		MaxRetries:   int(MaxRetries),
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
		lg:       lg,
		streamID: StreamID,
	}
}

func (m *memory) SendMessage(ctx context.Context, q *entities.Message) error {
	go m.lg.WriteLog(ctx, enum.DEBUG, "SendMessage", enum.DEBUG_DESCRIPTION)

	data, err := jsoniter.Marshal(q)
	if err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "SendMessage", err.Error())
		return err
	}

	if err := m.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: m.streamID,
		Values: map[string]any{
			enum.PAYLOAD: data,
		},
		ID: q.GetID(),
	}).Err(); err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "SendMessage", err.Error())
		return err
	}
	return nil
}

func (m *memory) AckMessage(ctx context.Context, messageId, groupId string) error {
	go m.lg.WriteLog(ctx, enum.DEBUG, "DeleteMessage", enum.DEBUG_DESCRIPTION)
	if err := m.rdb.XAck(ctx, m.streamID, groupId, messageId).Err(); err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "AckMessage", err.Error())
		return err
	}
	return nil
}

func (m *memory) GetMessage(ctx context.Context, workerId, groupId string) (*entities.Message, error) {
	go m.lg.WriteLog(ctx, enum.DEBUG, "GetMessage", enum.DEBUG_DESCRIPTION)

	response, err := m.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupId,
		Consumer: workerId,
		Streams:  []string{m.streamID, ">"},
		Block:    100 * time.Millisecond,
		Count:    1,
	}).Result()

	if err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "GetMessage", err.Error())
		return nil, err
	}

	message := response[0].Messages[0].Values[enum.PAYLOAD].(string)

	var q *entities.Message
	if err := jsoniter.Unmarshal([]byte(message), q); err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "GetMessage", err.Error())
		return nil, err
	}

	return q, nil
}

func (m *memory) Group(groupId string) error {
	return m.rdb.XGroupCreateMkStream(context.Background(), m.streamID, groupId, enum.START_FROM_NEW_MESSAGE).Err()
}
