package config

import (
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ENV         string           `yaml:"env"`
	PostgresDsn string           `yaml:"storage_dsn"`
	Http        HttpConfig       `yaml:"http"`
	Https       HttpsConfig      `yaml:"https"`
	Kafka       KafkaConfig      `yaml:"kafka"`
	Redis       RedisConfig      `yaml:"redis"`
	Grpc        GrpcConfig       `yaml:"grpc"`
	Monitoring  MonitoringConfig `yaml:"monitoring"`
}
type HttpConfig struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}
type HttpsConfig struct {
	PrivateKey     string        `yaml:"private_key"`
	Certificate    string        `yaml:"certificate"`
	Address        string        `yaml:"address"`
	Timeout        time.Duration `yaml:"timeout"`
	Authentication Auth          `yaml:"auth"`
}

type GrpcConfig struct {
	Address string `yaml:"address"`
}

type KafkaConfig struct {
	Hosts []string `yaml:"Hosts"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

type Auth struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type MonitoringConfig struct {
	Prometheus struct {
		Address string `yaml:"address"`
	} `yaml:"prometheus-exporter"`
	Jaeger struct {
		Address string `yaml:"address"`
	} `yaml:"jaeger"`
}

func New(reader io.Reader) *Config {
	var config Config
	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		return nil
	}
	return &config
}
