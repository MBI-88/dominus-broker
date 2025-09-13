package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	grpcConnKey = "DominusKeyConnectionMBI@88"
)

type config struct {
	RestPort      int64
	GrpcPort      int64
	TopicLimit    int
	QueueLimit    int
	KeyFile       string
	SslCaCert     string
	SslCert       string
	ApiToken      string
	AllowOrigins  string
	ConnectionKey string
	Logs          string
	Host          string
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
	s.TopicLimit = viper.GetInt("TOPIC_LIMIT")
	s.QueueLimit = viper.GetInt("QUEUE_LIMIT")
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
		log.Println("[+] GrpcPort was setted to default option 5000")
		s.GrpcPort = 5000
	} else if s.RestPort == 0 {
		log.Println("[+] RestPort was setted to default option 8000")
		s.RestPort = 8000
	} else if s.ApiToken == "" {
		panic("[-] ApiToken must be different from empty")
	} else if s.AllowOrigins == "" {
		log.Println("[+] AllowOrigins was setted to default option 0.0.0.0/24")
		s.AllowOrigins = "0.0.0.0/24"
	} else if s.Logs == "" {
		log.Println("[+] Logs dir was setted to default option ./logs")
		s.Logs = "./logs"
	} else if s.TopicLimit == 0 {
		log.Println("[+] TopicLimit was setted to default option 100")
		s.TopicLimit = 100
	}
}

type restConfigInt interface {
	GetEnvVarTest() config
	GetEnvVar(prod bool) config
}

func NewConfig() restConfigInt {
	return new(config)
}
