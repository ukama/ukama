package pkg

import (
	"github.com/ukama/ukamaX/common/config"
	"time"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	EtcdHost          string
	Grpc              config.Grpc
	Http              config.Http
	DialTimeoutSecond time.Duration
	NodeMetricsPort   int
	Dns               *DnsConfig
}

type DnsConfig struct {
	NodeDomain string // nodes domain like : ukama.node or mesh.node
}

func NewConfig() *Config {
	return &Config{
		EtcdHost:          "localhost:2379",
		DialTimeoutSecond: 3 * time.Second,
		Grpc: config.Grpc{
			Port: 9090,
		},
		Http: config.Http{
			Port: 8080,
		},
		NodeMetricsPort: 10250,
		Dns: &DnsConfig{
			NodeDomain: "node.mesh",
		},
	}
}
