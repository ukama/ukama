package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig  `mapstructure:",squash"`
	EtcdHost          string
	DialTimeoutSecond time.Duration
	NodeMetricsPort   int
	Dns               *DnsConfig
	Grpc              *uconf.Grpc      `default:"{}"`
	Queue             *uconf.Queue     `default:"{}"`
	Metrics           *uconf.Metrics   `default:"{}"`
	Timeout           time.Duration    `default:"3s"`
	MsgClient         *uconf.MsgClient `default:"{}"`
	Service           *uconf.Service
	Registry          string `default:"http://registry:8080"`
	Org               string `default:""`
}

type DnsConfig struct {
	NodeDomain string // nodes domain like : ukama.node or mesh.node
}

func NewConfig(name string) *Config {
	return &Config{
		EtcdHost:          "localhost:2379",
		DialTimeoutSecond: 5 * time.Second,
		NodeMetricsPort:   10250,
		Dns: &DnsConfig{
			NodeDomain: "node.mesh",
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.mesh.node.online", "event.cloud.mesh.node.offline", "event.cloud.registry.node.assigned", "event.cloud.registry.node.release"},
		},
		Registry: "http://192.168.0.23:8075",
	}
}
