package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	DB               *uconf.Database `default:"{}"`
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc   `default:"{}"`
	Queue            *uconf.Queue  `default:"{}"`
	Timeout          time.Duration `default:"20s"`
	Service          *uconf.Service
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgName          string           `default:"ukama"`
	RegistryHost     string           `default:"registry"`
	LatestConfigHash string 
	StoreUser        string
	StoreUrl         string
	AccessToken      string
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
