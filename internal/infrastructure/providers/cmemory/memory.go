package cmemory

import (
	"context"
	"crypto/tls"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/enum"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type memory struct {
	rdb       *redis.Client
	exp       int
	lg        adapters.Logs
	batchSize int64
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
	ExpirationTime int,
	BatchSize int64,
	lg adapters.Logs,
) adapters.MemoryClient {

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
		rdb:       client,
		exp:       ExpirationTime,
		lg:        lg,
		batchSize: BatchSize,
	}
}

func (m *memory) SendMessage(ctx context.Context, q *entities.Queue) error {
	go m.lg.WriteLog(ctx, enum.DEBUG, "SendMessage", enum.DEBUG_DESCRIPTION)
	if _, err := m.rdb.Set(ctx, q.ID.String(), q, time.Duration(m.exp)*time.Hour).Result(); err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "SendMessage", err.Error())
		return err
	}
	return nil
}

func (m *memory) DeleteMessage(ctx context.Context, key string) error {
	go m.lg.WriteLog(ctx, enum.DEBUG, "DeleteMessage", enum.DEBUG_DESCRIPTION)
	if _, err := m.rdb.Del(ctx, key).Result(); err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "DeleteMessage", err.Error())
		return err
	}
	return nil
}

func (m *memory) GetMessage(ctx context.Context, key string) (*entities.Queue, error) {
	go m.lg.WriteLog(ctx, enum.DEBUG, "GetMessage", enum.DEBUG_DESCRIPTION)
	result, err := m.rdb.Get(ctx, key).Result()
	if err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "GetMessage", err.Error())
		return nil, redis.ErrClosed
	}

	var q *entities.Queue
	if err := jsoniter.Unmarshal([]byte(result), q); err != nil {
		go m.lg.WriteLog(ctx, enum.ERROR, "GetMessage", err.Error())
		return nil, err
	}
	return q, nil
}

func (m *memory) GetKeys(ctx context.Context, mem entities.Memory) error {
	go m.lg.WriteLog(ctx, enum.DEBUG, "GetKeys", enum.DEBUG_DESCRIPTION)

	var cursor uint64
	for {
		keys, nextCursor, error := m.rdb.Scan(ctx, cursor, enum.ALL_KEYS, m.batchSize).Result()
		if error != nil {
			go m.lg.WriteLog(ctx, enum.ERROR, "GetKeys", error.Error())
			return error
		}

		for _, key := range keys {
			go m.lg.WriteLog(ctx, enum.INFO, "GetKeys", key)
			mem.Set(key)
		}

		cursor = nextCursor

		if cursor == 0 {
			go m.lg.WriteLog(ctx, enum.DEBUG, "GetKeys", enum.DEBUG_DESCRIPTION)
			break
		}
	}
	return nil
}
