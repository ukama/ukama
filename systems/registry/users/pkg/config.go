package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
	metric "github.com/ukama/ukama/systems/common/metrics"
)

const (
	NumberOfActiveUsers   = "platform_active_users"
	NumberOfInactiveUsers = "platform_inactive_users"
	GaugeType             = "gauge"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Service           *config.Service   `default:"{}"`
	MsgClient         *config.MsgClient `default:"{}"`
	Org               string            `default:"org:9090"`
	OrgOWnerName      string            `default:"Ukama Root"`
	OrgOWnerEmail     string            `default:"hello@ukama.com"`
	OrgOWnerPhone     string            `default:"0000000000"`
	OrgOWnerUUID      string
	Queue             *config.Queue `default:"{}"`
	PushGatewayHost   string        `default:"http://localhost:9091"`
}

var UserMetric = []metric.MetricConfig{
	{
		Name:   NumberOfActiveUsers,
		Type:   GaugeType,
		Labels: map[string]string{"state": "active"},
		Value:  0,
	},
	{
		Name:   NumberOfInactiveUsers,
		Type:   GaugeType,
		Labels: map[string]string{"state": "inactive"},
		Value:  0,
	},
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 7 * time.Second,
		},
	}
}
