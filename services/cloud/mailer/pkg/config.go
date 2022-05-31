package pkg

import (
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics  `default:"{}"`
	Server            *rest.HttpConfig `default:"{}"`
	Queue             *QueueConfig     `default:"{}"`
	Smtp              *SmtpConfig      `default:"{}"`
	TemplatesPath     string           `default:"templates"`
}

type ServerConfig struct {
	Port int `default:"10251"`
}

type QueueConfig struct {
	config.Queue `mapstructure:",squash"`
	QueueName    string `default:"mailer"`
}

type SmtpConfig struct {
	From     string `default:"hello@dev.ukama.com" validation:"required"`
	Host     string `default:"localhost" validation:"required"`
	Port     int    `default:"25" validation:"required"`
	Password string
	Username string
}
