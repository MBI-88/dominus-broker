package config

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type grpcConfig struct {
	GRPCPort      int64  `json:"grpc_port"`
	ConnectionKey string `json:"connection_key"`
}

type restConfig struct {
	RestPort     int64  `json:"rest_port"`
	ApiToken     string `json:"api_token"`
	AllowOrigins string `json:"allow_origins"`
}

type domainConfig struct {
	TopicLimit int `json:"topic_limit"`
	QueueLimit int `json:"queue_limit"`
}

type certConfig struct {
	KeyFile   string `json:"key_file"`
	SslCaCert string `json:"ssl_ca_cert"`
	SslCert   string `json:"ssl_cert"`
}

type redisConfig struct {
	PoolSize     int64  `json:"pool_size"`
	IdleConn     int64  `json:"idle_conn"`
	MaxRetries   int64  `json:"max_retries"`
	DialTimeOut  int64  `json:"dial_time_out"`
	ReadTimeOut  int64  `json:"read_time_out"`
	WriteTimeOut int64  `json:"write_time_out"`
	Port         int64  `json:"port"`
	Db           int64  `json:"db"`
	Host         string `json:"host"`
	Password     string `json:"password"`
	Tls          bool   `json:"tls"`
	Ssl          bool   `json:"ssl"`
}

type sqsConfig struct {
	MaxAttempts     int64  `json:"max_attempts"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SqsQueueURL     string `json:"sqs_queue_url"`
	Endpoint        string `json:"endpoint"`
	RetrieMode      string `json:"retrie_mode"`
	RoleARN         string `json:"role_arn"`
	ExternaID       string `json:"external_id"`
}

type providerConfig struct {
	RedisConfig redisConfig `json:"redis_config"`
	SqsConfig   sqsConfig   `json:"sqs_config"`
}

type infraConfig struct {
	LogMode string `json:"log_mode"`
	Host    string `json:"host"`
	LogURL  string `json:"log_url"`
}

type Config struct {
	GrpcConfig     grpcConfig     `json:"grpc_config"`
	RestConfig     restConfig     `json:"rest_config"`
	CertConfig     certConfig     `json:"cert_config"`
	ProviderConfig providerConfig `json:"provider_config"`
	InfraConfig    infraConfig    `json:"infra_config"`
}

func NewConfig(prod bool) Config {
	var (
		cf Config
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
		if err := viper.Unmarshal(&cf); err != nil {
			panic(err)
		}
	}

	return cf
}
