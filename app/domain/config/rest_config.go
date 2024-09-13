package config

import "github.com/spf13/viper"

type restConfig struct {
	ApiToken                          string
	Port                              int
	SslCert                           string
	WriteTimeout                      int
	ReadTimeout                       int
	IdleTimeout                       int
	MaxConnsPerIp                     int
	MaxRequestPerConn                 int
	MaxRequestBodySize                int
	ReduceMemoryUsage                 bool
	DisablePreparseMultipartForm      bool
	DisableHeaderNamesNormalizing     bool
	SleepWhenConcurrencyLimitExcedeed int
	NoDefaultDate                     bool
	KeepHijackedConns                 bool
	CloseOnShutdown                   bool
	StreamRequestBody                 bool
	Dsn                               string
	Cidr                              string
}

func (s *restConfig) setEnv() {
	s.ApiToken = viper.GetString("API_TOKEN")
	s.Port = viper.GetInt("PORT")
	s.SslCert = viper.GetString("SSL_CERT")
	s.WriteTimeout = viper.GetInt("WRITE_TIMEOUT")
	s.ReadTimeout = viper.GetInt("READ_TIMEOUT")
	s.IdleTimeout = viper.GetInt("IDLE_TIMOUT")
	s.MaxConnsPerIp = viper.GetInt("MAX_CONNS_PER_IP")
	s.MaxRequestPerConn = viper.GetInt("MAX_REQUEST_PER_CONN")
	s.MaxRequestBodySize = viper.GetInt("MAX_REQUEST_BODY_SIZE")
	s.ReduceMemoryUsage = viper.GetBool("REDUCE_MEMORY_USAGE")
	s.DisablePreparseMultipartForm = viper.GetBool("DISABLE_PREPARSE_MULTIPART_FORM")
	s.DisableHeaderNamesNormalizing = viper.GetBool("DISABLE_HEADER_NAMES_NORMALIZING")
	s.SleepWhenConcurrencyLimitExcedeed = viper.GetInt("SLEEP_WHEN_CONCURRENCY_LIMIT_EXCEDEED")
	s.NoDefaultDate = viper.GetBool("NO_DEFAULT_DATE")
	s.KeepHijackedConns = viper.GetBool("KEEP_HIJACKED_CONNS")
	s.CloseOnShutdown = viper.GetBool("CLOSE_ON_SHUTDOWN")
	s.StreamRequestBody = viper.GetBool("STREAM_REQUEST_BODY")
	s.Dsn = viper.GetString("DSN")
	s.Cidr = viper.GetString("CIDR")

}

func (s restConfig) GetEnvVar(prod bool) restConfig {
	if prod {
		viper.SetEnvPrefix("")
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

func (s restConfig) GetEnvVarTest() restConfig {
	viper.SetConfigFile("./../.env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("[-] File not found!")
		}
	}
	s.setEnv()
	return s
}

type restConfigInt interface {
	GetEnvVarTest() restConfig
	GetEnvVar(prod bool) restConfig
}

func NewRestConfig() restConfigInt {
	return new(restConfig)
}
