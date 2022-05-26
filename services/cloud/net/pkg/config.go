package pkg

import (
	"time"

	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	EtcdHost          string
	Grpc              config.Grpc
	Http              rest.HttpConfig
	DialTimeoutSecond time.Duration
	NodeMetricsPort   int
	Dns               *DnsConfig
	Metrics           config.Metrics
	OrgMetricsTarget  OrgMetricsConf
}

type OrgMetricsConf struct {
	Url            string
	ScrapeInterval time.Duration
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
		Http: rest.HttpConfig{
			Port: 8080,
		},
		NodeMetricsPort: 10250,
		Dns: &DnsConfig{
			NodeDomain: "node.mesh",
		},
		Metrics: config.Metrics{
			Port:    10250,
			Enabled: true,
		},
		OrgMetricsTarget: OrgMetricsConf{
			Url:            "https://localhost:10251", // full url with port path and http
			ScrapeInterval: 1 * time.Minute,
		},
	}
}
