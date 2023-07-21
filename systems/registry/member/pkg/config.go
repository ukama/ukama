package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
	metric "github.com/ukama/ukama/systems/common/metrics"
)

const (
	NumberOfMembers         = "number_of_members"
	NumberOfActiveMembers   = "active_members"
	NumberOfInactiveMembers = "inactive_members"
	GaugeType               = "gauge"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgRegistryHost  string           `default:"http://org:8080"`
	OrgOwnerUUID     string
	OrgId            string
	OrgName          string
	Service          *uconf.Service
	PushGateway      string `default:"http://localhost:9091"`
}

var MemberMetric = []metric.MetricConfig{
	{
		Name:   NumberOfMembers,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
		Value:  0,
	},
	{
		Name:   NumberOfActiveMembers,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
		Value:  0,
	},
	{
		Name:   NumberOfInactiveMembers,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
		Value:  0,
	},
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}
