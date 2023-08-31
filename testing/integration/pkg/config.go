package pkg

import (
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Service          *uconf.Service
	System           SystemConfig
	LogLevel         int `default:"4`
}

type SystemConfig struct {
	Dataplan   string `default:"http://localhost:8074"`
	Init       string `default:"http://localhost:8071"`
	Registry   string `default:"http://localhost:8075"`
	Metrics    string `default:"http://localhost:8072"`
	Subscriber string `default:"http://localhost:8078"`
	Billing    string `default:"http://localhost:8079"`
	MessageBus string "amqp://guest:guest@localhost:5672/"
}

func NewConfig() *Config {
	return &Config{
		System: SystemConfig{
			MessageBus: "amqp://guest:guest@localhost:5672/",
			Dataplan:   "http://localhost:8074",
			Init:       "http://localhost:8071",
			Registry:   "http://localhost:8075",
			Metrics:    "http://localhost:8072",
			Subscriber: "http://localhost:8078",
			Billing:    "http://localhost:8079",
		},
	}
}
