package config

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type GrpcConfig struct {
	GRPCPort      int64
	ConnectionKey string
}

type RestConfig struct {
	RestPort     int64
	ApiToken     string
	AllowOrigins string
}

type CertConfig struct {
	KeyFile   string
	SslCaCert string
	SslCert   string
}

type RedisConfig struct {
	Port         int64
	PoolSize     int
	MaxRetries   int
	DialTimeOut  int
	ReadTimeOut  int
	WriteTimeOut int
	MemoryDB     int
	CheckerDB    int
	IdPotencyEX  int
	Host         string
	Password     string
	Username     string
	StreamID     string
	GroupID      string
	Tls          bool
}

type InfraConfig struct {
	LogMode string
	LogURL  string
}

type Config struct {
	GrpcConfig  GrpcConfig
	RestConfig  RestConfig
	CertConfig  CertConfig
	RedisConfig RedisConfig
	InfraConfig InfraConfig
}

func NewConfig(prod bool) *Config {
	var cf Config

	if prod {
		viper.AutomaticEnv()
		raw := viper.GetString("APP_CONFIG")
		if raw == "" {
			panic("APP_CONFIG is empty or not set")
		}
		if err := json.Unmarshal([]byte(raw), &cf); err != nil {
			panic(err)
		}
	} else {
		viper.SetConfigType("json")
		viper.SetConfigFile("./../env/env.dev.json")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				panic("[-] File not found!")
			}
		}
		if err := viper.Unmarshal(&cf); err != nil {
			panic(err)
		}
	}

	return &cf
}
