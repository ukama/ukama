package pkg

import (
<<<<<<< HEAD
	"time"

=======
>>>>>>> subscriber-sys_sim-manager
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
<<<<<<< HEAD
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"5s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.simManager.sim.allocation"},
		},
	}
=======
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
>>>>>>> subscriber-sys_sim-manager
}
