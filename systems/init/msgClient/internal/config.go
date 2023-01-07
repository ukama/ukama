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
}

func NewConfig() *Config {
	return &Config{

		Grpc: &uconf.Grpc{
			Port: 9095,
		},
	}
}
