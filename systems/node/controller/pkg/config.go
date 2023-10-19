package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	DB               *uconf.Database `default:"{}"`
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc   `default:"{}"`
	Queue            *uconf.Queue  `default:"{}"`
	Timeout          time.Duration `default:"20s"`
	Service          *uconf.Service
	MsgClient        *config.MsgClient `default:"{}"`
	OrgName          string            `default:"ukama"`
	RegistryHost     string            `default:"http://org:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.registry.node.node.create",
			},
		},
	}
}
