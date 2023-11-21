package pkg

import (
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Service          *uconf.Service
	System           SystemConfig
	LogLevel         int    `default:"4`
	EncryptKey       string `default:"KEY"`
	OrgId            string `default:"ORG_ID"`
	OrgName          string `default:"ORG_NAME"`
	OrgOwnerId       string `default:"ORG_OWNER_ID"`
}

type SystemConfig struct {
	Dataplan        string `default:"http://localhost:8074"`
	Init            string `default:"http://localhost:8071"`
	Registry        string `default:"http://localhost:8075"`
	Metrics         string `default:"http://localhost:8072"`
	Subscriber      string `default:"http://localhost:8078"`
	Billing         string `default:"http://localhost:8079"`
	BillingProvider string `default:"http://localhost:3000"`
	Nucleus         string `default:"http://localhost:8060"`
	MessageBus      string `default:"amqp://guest:guest@localhost:5672/"`
}

func NewConfig() *Config {
	return &Config{
		OrgName:    "ukama-test-org",
		OrgId:      "8c6c2bec-5f90-4fee-8ffd-ee6456abf4fc",
		OrgOwnerId: "018688fa-d861-4e7b-b119-ffc5e1637ba8",
		EncryptKey: "the-key-has-to-be-32-bytes-long!",
		// EncryptKey: "ad5ca9af-9d60-47e3-9b8c-43d354d5f6c7",
		System: SystemConfig{
			MessageBus:      "amqp://guest:guest@localhost:5672/",
			Dataplan:        "http://localhost:8074",
			Init:            "http://localhost:8071",
			Registry:        "http://localhost:8075",
			Metrics:         "http://localhost:8072",
			Subscriber:      "http://localhost:8078",
			Billing:         "http://localhost:8680",
			BillingProvider: "http://localhost:3000",
			Nucleus:         "http://localhost:8060",
		},
	}
}
