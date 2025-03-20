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
	return *s
}

type restConfigInt interface {
	//Get variables for testing
	GetEnvVarTest() config
	//Get variables for development/production.
	//In production mode variables are setted up from system environment
	//In development mode variables are setted up form env file
	GetEnvVar(prod bool) config
}

func NewConfig() restConfigInt {
	return new(config)
}
