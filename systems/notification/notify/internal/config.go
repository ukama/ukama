package internal

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Metrics          *uconf.Metrics   `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	Queue            *uconf.Queue     `default:"{}"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	System           string `default:"notification"`
	OrgName          string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},

		Service: uconf.LoadServiceHostConfig(name),

		MsgClient: &uconf.MsgClient{
			Host:    "msg-client-notification:9095",
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.global.{{ .Org}}.nucleus.org.notification.sent",
				"event.cloud.global.{{ .Org}}.nucleus.user.notification.sent",
				"event.cloud.local.{{ .Org}}.registry.network.notification.sent",
				"event.cloud.local.{{ .Org}}.registry.node.notification.sent",
			},
		},
	}
}
