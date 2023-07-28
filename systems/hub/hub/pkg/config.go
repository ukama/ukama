package pkg

import (
	"time"

	cors "github.com/gin-contrib/cors"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics
	Server            rest.HttpConfig
	Storage           MinioConfig
	Chunker           ChunkerConfig
	Service           *config.Service
	Queue             *config.Queue     `default:"{}"`
	MsgClient         *config.MsgClient `default:"{}"`
}

type MinioConfig struct {
	TimeoutSecond      int
	Endpoint           string
	AccessKey          string
	SecretKey          string
	BucketSuffix       string
	Region             string
	SkipBucketCreation bool
}

type ChunkerConfig struct {
	Host          string
	TimeoutSecond int
}

func NewConfig(name string) *Config {
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

		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}
