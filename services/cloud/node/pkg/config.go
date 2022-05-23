package pkg

import (
	uconf "github.com/ukama/ukama/services/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               uconf.Database
	Grpc             uconf.Grpc
	Queue            uconf.Queue
}
