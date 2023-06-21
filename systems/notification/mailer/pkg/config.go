package pkg

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Mailer struct {
	Host     string `default:"localhost"`
	Port     int    `default:"587"`
	Username string `default:""`
	Password string `default:""`
	From     string `default:"hello@dev.ukama.com"`
}

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Queue            *uconf.Queue    `default:"{}"`
	Timeout          time.Duration   `default:"10s"`
	Service          *uconf.Service
	Mailer           *Mailer
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		Mailer:  LoadMailerHostConfig(name),
	}
}

func LoadMailerHostConfig(name string) *Mailer {
	s := &Mailer{}
	mailerHost := "_MAILER_HOST"
	mailerPort := "_MAILER_PORT"
	mailerUsername := "_MAILER_USERNAME"
	mailerPassword := "_MAILER_PASSWORD"
	mailerFrom := "_MAILER_FROM"

	val, present := os.LookupEnv(strings.ToUpper(name + mailerFrom))
	if present {
		s.From = val
	} else {
		logrus.Errorf("%s mailer from env not found", name)
	}

	val, present = os.LookupEnv(strings.ToUpper(name + mailerHost))
	if present {
		s.Host = val

	} else {
		logrus.Errorf("%s mailer host env not found", name)
	}

	val, present = os.LookupEnv(strings.ToUpper(name + mailerPort))

	if present {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Errorf("Failed to convert %s mailer port to int: %s", name, err)
		} else {
			s.Port = port
		}
	} else {
		logrus.Errorf("%s mailer port env not found", name)
	}

	val, present = os.LookupEnv(strings.ToUpper(name + mailerUsername))
	if present {
		s.Username = val
	} else {
		logrus.Errorf("%s mailer username env not found", name)
	}

	val, present = os.LookupEnv(strings.ToUpper(name + mailerPassword))
	if present {
		s.Password = val
	} else {
		logrus.Errorf("%s mailer password env not found", name)
	}
	return s
}
