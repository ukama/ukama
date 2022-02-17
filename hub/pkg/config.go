package pkg

import (
	cors "github.com/gin-contrib/cors"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics
	Server            rest.HttpConfig
	Storage           MinioConfig
	Chunker           ChunkerConfig
}

type MinioConfig struct {
	TimeoutSecond int
	Endpoint      string
	AccessKey     string
	SecretKey     string
	BucketSuffix  string
	Region        string
}

type ChunkerConfig struct {
	Host          string
	TimeoutSecond int
}

func NewConfig() *Config {
	return &Config{
		Server: rest.HttpConfig{
			Port: 8080,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost"},
			},
		},
		Metrics: config.DefaultMetrics(),
		Storage: MinioConfig{
			Endpoint:      "localhost:9000",
			AccessKey:     "minio",
			SecretKey:     "minio123",
			BucketSuffix:  "local-test",
			TimeoutSecond: 3,
		},
		Chunker: ChunkerConfig{
			Host:          "http://localhost:8080",
			TimeoutSecond: 3,
		},
	}
}
