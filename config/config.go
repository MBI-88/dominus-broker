package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type GrpcConfig struct {
	Port     int64  `json:"port" mapstructure:"port"`
	ApiToken string `json:"api_token" mapstructure:"api_token"`
}

type RestConfig struct {
	Port         int64  `json:"port" mapstructure:"port"`
	ApiToken     string `json:"api_token" mapstructure:"api_token"`
	AllowOrigins string `json:"allow_origins" mapstructure:"allow_origins"`
}

type CertConfig struct {
	KeyFile   string `json:"key_file" mapstructure:"key_file"`
	SslCaCert string `json:"ssl_ca_cert" mapstructure:"ssl_ca_cert"`
	SslCert   string `json:"ssl_cert" mapstructure:"ssl_cert"`
}

type RedisConfig struct {
	Port        int64  `json:"port" mapstructure:"port"`
	MemoryDB    int    `json:"memory_db" mapstructure:"memory_db"`
	CheckerDB   int    `json:"checker_db" mapstructure:"checker_db"`
	IdPotencyEx int    `json:"id_potency_ex" mapstructure:"id_potency_ex"`
	Host        string `json:"host" mapstructure:"host"`
	Password    string `json:"password" mapstructure:"password"`
	Username    string `json:"username" mapstructure:"username"`
	StreamID    string `json:"stream_id" mapstructure:"stream_id"`
	GroupID     string `json:"group_id" mapstructure:"group_id"`
	Tls         bool   `json:"tls" mapstructure:"tls"`
}

type LogConfig struct {
	LogMode string `json:"log_mode" mapstructure:"log_mode"`
	LogUrl  string `json:"log_url" mapstructure:"log_url"`
}

type Config struct {
	GrpcConfig  GrpcConfig  `json:"grpc_config" mapstructure:"grpc_config"`
	RestConfig  RestConfig  `json:"rest_config" mapstructure:"rest_config"`
	CertConfig  CertConfig  `json:"cert_config" mapstructure:"cert_config"`
	RedisConfig RedisConfig `json:"redis_config" mapstructure:"redis_config"`
	LogConfig   LogConfig   `json:"log_config" mapstructure:"infra_config"`
}

func NewConfig() *Config {
	var cf Config

	viper.SetConfigType("json")
	viper.AutomaticEnv()
	raw := viper.GetString("APP_CONFIG")
	if raw == "" {
		panic("APP_CONFIG is empty or not set")
	}
	if err := viper.ReadConfig(strings.NewReader(raw)); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&cf); err != nil {
		panic(err)
	}
	if err := cf.validateGrpcConfig(); err != nil {
		panic(err)
	}
	if err := cf.validateRedis(); err != nil {
		panic(err)
	}
	if err := cf.validateRestConfig(); err != nil {
		panic(err)
	}
	return &cf
}

func (c *Config) validateRedis() error {
	if c.RedisConfig.Host == "" {
		return fmt.Errorf("RedisConfig.Host empty")
	}
	if c.RedisConfig.Port == 0 {
		return fmt.Errorf("RedisConfig.Port empty")
	}
	return nil
}

func (c *Config) validateGrpcConfig() error {
	if c.GrpcConfig.Port == 0 {
		return fmt.Errorf("GrpcConfig.GrpcPort empty")
	}
	if c.GrpcConfig.ApiToken == "" {
		return fmt.Errorf("GrpcConfig.ApiToken empty")
	}
	return nil
}

func (c *Config) validateRestConfig() error {
	if c.RestConfig.Port == 0 {
		return fmt.Errorf("RestConfig.Port empty")
	}
	if c.RestConfig.ApiToken == "" {
		return fmt.Errorf("RestConfig.ApiToken empty")
	}
	if c.RestConfig.AllowOrigins == "" {
		return fmt.Errorf("RestConfig.AllowOrigins empty")
	}
	return nil
}
