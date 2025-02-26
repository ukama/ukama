package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
	evt "github.com/ukama/ukama/systems/common/events"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	Service          *uconf.Service
	OrgName          string
	MsgClient        *uconf.MsgClient `default:"{}"`
}



func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				evt.EventRoutingKey[evt.EventSiteCreate],
			},
		},
	}
}