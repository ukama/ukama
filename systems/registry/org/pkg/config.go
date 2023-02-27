// package pkg

// import (
// 	"github.com/ukama/ukama/systems/common/config"
// )

//	type Config struct {
//		config.BaseConfig `mapstructure:",squash"`
//		DB                *config.Database `default:"{}"`
//		Grpc              *config.Grpc     `default:"{}"`
//		UsersHost         string           `default:"users:9090"`
//		Metrics           *config.Metrics  `default:"{}"`
//	}
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
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	UsersHost         string           `default:"users:9090"`

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