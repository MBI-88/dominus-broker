package cchecker

import (
	"context"
	"crypto/tls"
	"dominus-project/internal/domain/repositories"
	"dominus-project/internal/domain/enum"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type checker struct {
	rdb *redis.Client
	exp int
}

func NewCheckerClient(
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
) repositories.ChckerClient {

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
	return &checker{
		rdb: client,
		exp: ExpirationTime,
	}
}

func (c *checker) SaveConsumer(ctx context.Context, key string) error {
	if _, err := c.rdb.SetNX(ctx, fmt.Sprintf("%s:%s", enum.ID_POTENCY_TAG, key), map[string]any{}, time.Duration(c.exp)*time.Second).Result(); err != nil {
		return err
	}
	return nil
}

func (c *checker) CheckConsumer(ctx context.Context, key string) bool {
	result, err := c.rdb.Exists(ctx, fmt.Sprintf("%s:%s", enum.ID_POTENCY_TAG, key)).Result()
	if err != nil {
		return false
	}
	if result == 1 {
		return true
	}
	return false
}
