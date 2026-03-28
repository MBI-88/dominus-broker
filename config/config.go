package config

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type GrpcConfig struct {
	GRPCPort      int64  `json:"grpc_port"`
	ConnectionKey string `json:"connection_key"`
}

type RestConfig struct {
	RestPort     int64  `json:"rest_port"`
	ApiToken     string `json:"api_token"`
	AllowOrigins string `json:"allow_origins"`
}

type domainConfig struct {
	TopicLimit int `json:"topic_limit"`
	QueueLimit int `json:"queue_limit"`
}

type CertConfig struct {
	KeyFile   string `json:"key_file"`
	SslCaCert string `json:"ssl_ca_cert"`
	SslCert   string `json:"ssl_cert"`
}

type RedisConfig struct {
	Port         int64  `json:"port"`
	PoolSize     int  `json:"pool_size"`
	MaxRetries   int  `json:"max_retries"`
	DialTimeOut  int  `json:"dial_time_out"`
	ReadTimeOut  int  `json:"read_time_out"`
	WriteTimeOut int  `json:"write_time_out"`
	MemoryDB     int    `json:"memory_db"`
	CheckerDB    int    `json:"checker_db"`
	IdPotency    int    `json:"id_potency"`
	Host         string `json:"host"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	StreamID     string `json:"stream_id"`
	GroupID      string `json:"group_id"`
	Tls          bool   `json:"tls"`
}

type InfraConfig struct {
	LogMode string `json:"log_mode"`
	LogURL  string `json:"log_url"`
}

type Config struct {
	GrpcConfig  *GrpcConfig  `json:"grpc_config"`
	RestConfig  *RestConfig  `json:"rest_config"`
	CertConfig  *CertConfig  `json:"cert_config"`
	RedisConfig *RedisConfig `json:"redis_config"`
	InfraConfig *InfraConfig `json:"infra_config"`
}

func NewConfig(prod bool) *Config {
	var (
		cf *Config
	)

	if prod {
		viper.AutomaticEnv()
		if err := json.Unmarshal([]byte(viper.GetString("APP_CONFIG")), &cf); err != nil {
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
		if err := viper.Unmarshal(cf); err != nil {
			panic(err)
		}
	}

	return cf
}
