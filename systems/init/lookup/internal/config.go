package internal

import (
	uconf "github.com/ukama/ukama/services/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Queue            *uconf.Queue    `default:"{}"`
	Metrics          *uconf.Metrics  `default:"{}"`
}
