package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
	metric "github.com/ukama/ukama/systems/common/metrics"
)
const (
	NumberOfUsers = "number_of_users"
	GaugeType           = "gauge"
)
type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Grpc              *config.Grpc     `default:"{}"`
	Service           *config.Service  `default:"{}"`
	MsgClient        *config.MsgClient `default:"{}"`
	OrgHost               string           `default:"org:9090"`
	OrgOWnerName      string           `default:"Ukama Root"`
	OrgOWnerEmail     string           `default:"hello@ukama.com"`
	OrgOWnerPhone     string           `default:"0000000000"`
	OrgOWnerUUID      string
	Queue            *config.Queue     `default:"{}"`
	PushMetricHost    string `default:"http://localhost:9091"`
	Org               string `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
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

var UserMetric = []metric.MetricConfig{{
	Name:   NumberOfUsers,
	Type:   GaugeType,
	Labels: map[string]string{"user":"","org": ""},
	Value:  0,
},
}