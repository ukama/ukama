package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"20s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	OrgName          string
	Http             HttpServices
}
type HttpServices struct {
	InitClient    string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Timeout: 3 * time.Second,
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 7 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.registry.site.site.create",
			},
		},
	}
}
