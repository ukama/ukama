package internal

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database   `default:"{}"`
	Grpc             *uconf.Grpc       `default:"{}"`
	Queue            *uconf.Queue      `default:"{}"`
	Metrics          *uconf.Metrics    `default:"{}"`
	Timeout          time.Duration     `default:"3s"`
	HeathCheck       HeathCheckRoutine `default:"{}"`
	System           string            `default:"init"`
}

type HeathCheckRoutine struct {
	Period      time.Duration `default:"60s"`
	AllowedMiss uint32        `default:"3"`
}

func NewConfig() *Config {
	return &Config{
		Grpc: &uconf.Grpc{
			Port: 9095,
		},
	}
}
