package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig     `mapstructure:",squash"`
	DB                   *uconf.Database   `default:"{}"`
	Grpc                 *uconf.Grpc       `default:"{}"`
	Queue                *uconf.Queue      `default:"{}"`
	Timeout              time.Duration     `default:"3s"`
	MsgClient            *config.MsgClient `default:"{}"`
	OrgRegistryHost      string            `default:"http://org:8080"`
	Service              *uconf.Service
	InvitationExpiryTime time.Time
	NotificationHost     string `default:"http://192.168.1.81:8089"`
	AuthLoginbaseURL     string `default:"http://localhost:4455/auth/login"`
	OrgName              string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 7 * time.Second,
		},
	}
}
