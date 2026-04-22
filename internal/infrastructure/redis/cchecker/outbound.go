package cchecker

import (
	"context"
	"crypto/tls"
	"dominus-broker/internal/domain/repositories"
	"dominus-broker/internal/infrastructure/enum"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type checker struct {
	rdb *redis.Client
	exp int
}

func NewCheckerClient(
	Port int64,
	Db int,
	Host string,
	Password string,
	Tls bool,
	Username string,
	ExpirationTime int,
) repositories.CheckerClient {

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
	return &checker{
		rdb: client,
		exp: ExpirationTime,
	}
}

func (c *checker) SaveConsumer(ctx context.Context, key string) error {
	if _, err := c.rdb.SetArgs(ctx, fmt.Sprintf("%s:%s", enum.IDEM_POTENCY_TAG, key), "", redis.SetArgs{
		Mode: "NX",
		TTL:  time.Duration(c.exp) * time.Second,
	}).Result(); err != nil {
		return fmt.Errorf("checker.SaveConsumer %s", err)
	}
	return nil
}

func (c *checker) CheckConsumer(ctx context.Context, key string) bool {
	result, err := c.rdb.Exists(ctx, fmt.Sprintf("%s:%s", enum.IDEM_POTENCY_TAG, key)).Result()
	if err != nil {
		return false
	}
	if result == 1 {
		return true
	}
	return false
}
