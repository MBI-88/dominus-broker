package redisclient

import (
	"context"
	"crypto/tls"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type redisClient struct {
	rdb        *redis.Client
	expiration int
}

func NewRedisClient(
	PoolSize int64,
	IdleConn int64,
	MaxRetries int64,
	DialTimeOut int64,
	ReadTimeOut int64,
	WriteTimeOut int64,
	Port int64,
	Db int64,
	Host string,
	Password string,
	Tls bool,
	Username string,
) adapters.MemoryClient {

	cfTls := &tls.Config{}

	op := &redis.Options{
		PoolSize:     int(PoolSize),
		IdleTimeout:  time.Duration(IdleConn),
		MaxRetries:   int(MaxRetries),
		DialTimeout:  time.Duration(DialTimeOut),
		ReadTimeout:  time.Duration(ReadTimeOut),
		WriteTimeout: time.Duration(WriteTimeOut),
		DB:           int(Db),
		Addr:         fmt.Sprintf("%s:%d", Host, Port),
		Password:     Password,
		TLSConfig:    cfTls,
		Username:     Username,
	}

	client := redis.NewClient(op)

	return &redisClient{
		rdb: client,
	}
}

func (r *redisClient) SendMessage(ctx context.Context, q *entities.Queue) error {
	status := r.rdb.Set(ctx, q.ID.String(), q, time.Duration(r.expiration) * time.Hour)
	if status.Err() != nil {

	}
	return nil
}

func (r *redisClient) DeleteMessage(ctx context.Context, ID string) error {
	status := r.rdb.Del(ctx, ID)
	if status.Err() != nil {

	}
	return nil
}

func (r *redisClient) GetMessage(ctx context.Context, key string) (*entities.Queue, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return  nil, redis.ErrClosed
	}

	var q *entities.Queue
	if err := jsoniter.Unmarshal([]byte(result), q); err != nil {
		return  nil, err
	}	

	return q, nil
}
