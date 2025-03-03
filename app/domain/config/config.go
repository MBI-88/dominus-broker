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
	WriteTimeout                     int
	ReadTimeout                      int
	IdleTimeout                      int
	MaxConnsPerIp                    int
	MaxRequestPerConn                int
	MaxRequestBodySize               int
	SleepWhenConcurrencyLimitExceded int
	Host                             string
	KeyFile                          string
	SslCaCert                        string
	SslCert                          string
	ApiToken                         string
	Dsn                              string
	Cidr                             string
	ConnectionKey                    string
	Database                         string
	Collections                      string
	Logs                             string
	StreamRequestBody                bool
	CloseOnShutdown                  bool
	KeepHijackedConns                bool
	NoDefaultDate                    bool
	DisableHeaderNamesNormalizing    bool
	DisablePreparseMultipartForm     bool
	ReduceMemoryUsage                bool
}

func (s *config) setEnv() {
	s.ApiToken = viper.GetString("API_TOKEN")
	s.RestPort = viper.GetInt64("REST_PORT")
	s.GrpcPort = viper.GetInt64("GRPC_PORT")
	s.SslCert = viper.GetString("SSL_CERT")
	s.SslCaCert = viper.GetString("SSL_CA")
	s.KeyFile = viper.GetString("KEY_FILE")
	s.WriteTimeout = viper.GetInt("WRITE_TIMEOUT")
	s.ReadTimeout = viper.GetInt("READ_TIMEOUT")
	s.IdleTimeout = viper.GetInt("IDLE_TIMOUT")
	s.MaxConnsPerIp = viper.GetInt("MAX_CONNS_PER_IP")
	s.MaxRequestPerConn = viper.GetInt("MAX_REQUEST_PER_CONN")
	s.MaxRequestBodySize = viper.GetInt("MAX_REQUEST_BODY_SIZE")
	s.ReduceMemoryUsage = viper.GetBool("REDUCE_MEMORY_USAGE")
	s.DisablePreparseMultipartForm = viper.GetBool("DISABLE_PREPARSE_MULTIPART_FORM")
	s.DisableHeaderNamesNormalizing = viper.GetBool("DISABLE_HEADER_NAMES_NORMALIZING")
	s.SleepWhenConcurrencyLimitExceded = viper.GetInt("SLEEP_WHEN_CONCURRENCY_LIMIT_EXCEDED")
	s.NoDefaultDate = viper.GetBool("NO_DEFAULT_DATE")
	s.KeepHijackedConns = viper.GetBool("KEEP_HIJACKED_CONNS")
	s.CloseOnShutdown = viper.GetBool("CLOSE_ON_SHUTDOWN")
	s.StreamRequestBody = viper.GetBool("STREAM_REQUEST_BODY")
	s.Dsn = viper.GetString("DSN")
	s.Cidr = viper.GetString("CIDR")
	s.ConnectionKey = grpcConnKey
	s.Database = viper.GetString("DATABASE")
	s.Collections = viper.GetString("COLLECTIONS")
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
	s.Database = "test"
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
