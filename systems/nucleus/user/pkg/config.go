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
	OwnerName         string            `default:"Ukama Root"`
	OwnerEmail        string            `default:"hello@ukama.com"`
	OwnerPhone        string            `default:"0000000000"`
	OwnerId           string
	AuthId            string
	Queue             *config.Queue `default:"{}"`
	PushGatewayHost   string        `default:"http://localhost:9091"`
	OrgName           string
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
			DbName: name + "s",
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 7 * time.Second,
		},
	}
}
