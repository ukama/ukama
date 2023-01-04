package internal

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)


type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Queue            *uconf.Queue    `default:"{}"`
	Metrics          *uconf.Metrics  `default:"{}"`
	Timeout          time.Duration
	MsgClient        *uconf.MsgClient `default:"{}"`
}

func NewConfig() *Config {
	return &Config{
		MsgClient: &uconf.MsgClient{
			ListnerRoutes: []string{""},
		},
	}
}
