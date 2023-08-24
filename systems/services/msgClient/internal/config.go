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
	OrgName          string
	MasterOrgName    string
	Shovel           Shovel
}

type Shovel struct {
	SrcProtocol     string `json:"src_protocol" default:"amqp091"`
	DestProtocol    string `default:"amqp091" json:"src-protocol"`
	SrcExchange     string `default:"amqp.Topic" json:"src-exchange"`
	SrcExchangeKey  string `json:"src-exchange-key,omitempty"`
	DestExchange    string `default:"amqp.Topic" json:"dest-exchange,omitempty"`
	DestExchangeKey string `json:"dest-exchange-key,omitempty"`
	DestQueue       string `json:"dest-queue,omitempty"`
	SrcQueue        string `json:"src-queue,omitempty"`
	SrcUri          string `json:"src-uri"`
	DestUri         string `json:"dest-uri"`
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
