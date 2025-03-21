package config

import (
	"github.com/spf13/viper"
)

var (
	grpcConnKey = "DominusKeyConnectionMBI@88"
)

type config struct {
	RestPort                         int64
	GrpcPort                         int64
	KeyFile                          string
	SslCaCert                        string
	SslCert                          string
	ApiToken                         string
	AllowOrigins                     string
	ConnectionKey                    string
	Logs                             string
	Host                             string
}

func (s *config) setEnv() {
	s.ApiToken = viper.GetString("API_TOKEN")
	s.RestPort = viper.GetInt64("REST_PORT")
	s.GrpcPort = viper.GetInt64("GRPC_PORT")
	s.SslCert = viper.GetString("SSL_CERT")
	s.SslCaCert = viper.GetString("SSL_CA")
	s.KeyFile = viper.GetString("KEY_FILE")
	s.AllowOrigins = viper.GetString("ALLOW_ORIGINS")
	s.ConnectionKey = grpcConnKey
	s.Logs = viper.GetString("LOGS")
	s.Host = viper.GetString("HOST")
}

func (s config) GetEnvVar(prod bool) config {
	if prod {
		viper.AutomaticEnv()
	} else {
		viper.SetConfigFile("./.env")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				panic("[-] File not found!")
			}
		}
	}
	s.setEnv()
	s.checkVars()
	return s
}

func (s *config) GetEnvVarTest() config {
	viper.SetConfigFile("./../.env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("[-] File not found!")
		}
	}
	s.setEnv()
	s.checkVars()
	return *s
}

func (s *config) checkVars() {
	if s.GrpcPort == 0 {
		panic("[-] GrcpPort must be different from 0")
	}else if s.RestPort == 0 {
		panic("[-] RestPort must be different from 0")
	}else if s.ApiToken == "" {
		panic("[-] ApiToken must be different from empty")
	}else if s.AllowOrigins == "" {
		panic("[-] AllowOrigins must be different from empty")
	}else if s.Logs == "" {
		panic("[-] Logs must be different from empty")
	}
}

type restConfigInt interface {
	GetEnvVarTest() config
	GetEnvVar(prod bool) config
}

func NewConfig() restConfigInt {
	return new(config)
}
