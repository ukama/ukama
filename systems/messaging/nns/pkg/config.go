package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	uconf.BaseConfig  `mapstructure:",squash"`
	EtcdHost          string
	Http              rest.HttpConfig
	DialTimeoutSecond time.Duration
	NodeMetricsPort   int
	Dns               *DnsConfig
	Grpc              *uconf.Grpc      `default:"{}"`
	Queue             *uconf.Queue     `default:"{}"`
	Metrics           *uconf.Metrics   `default:"{}"`
	Timeout           time.Duration    `default:"3s"`
	MsgClient         *uconf.MsgClient `default:"{}"`
	Service           *uconf.Service
	Registry          string
}

type DnsConfig struct {
	NodeDomain string // nodes domain like : ukama.node or mesh.node
}

func NewConfig(name string) *Config {
	return &Config{
		EtcdHost:          "localhost:2379",
		DialTimeoutSecond: 3 * time.Second,
		NodeMetricsPort:   10250,
		Dns: &DnsConfig{
			NodeDomain: "node.mesh",
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.lookup.organization.create"},
		},
		Registry: "gateway.registry:8080",
	}
}
