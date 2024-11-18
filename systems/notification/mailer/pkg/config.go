/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type MailerConfig struct {
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
	Timeout          time.Duration   `default:"50s"`
	Service          *uconf.Service
	Mailer           *MailerConfig
	TemplatesPath    string `default:"templates"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		Mailer:  loadMailerConfig(name),
	}
}

func envKey(name, suffix string) string {
	return strings.ToUpper(name + suffix)
}

func getEnvOrLog(key, description string) (string, bool) {
	val, exists := os.LookupEnv(key)
	if !exists {
		logrus.Errorf("%s env not found: %s", description, key)
		return "", false
	}
	return val, true
}

func loadMailerConfig(name string) *MailerConfig {
	const (
		hostSuffix     = "_HOST"
		portSuffix     = "_PORT"
		usernameSuffix = "_USERNAME"
		passwordSuffix = "_PASSWORD"
		fromSuffix     = "_FROM"
	)

	config := &MailerConfig{}

	if val, ok := getEnvOrLog(envKey(name, fromSuffix), "mailer from"); ok {
		config.From = val
	}

	if val, ok := getEnvOrLog(envKey(name, hostSuffix), "mailer host"); ok {
		config.Host = val
	}

	if val, ok := getEnvOrLog(envKey(name, portSuffix), "mailer port"); ok {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Errorf("Invalid port value for %s: %v", envKey(name, portSuffix), err)
		} else {
			config.Port = port
		}
	}

	if val, ok := getEnvOrLog(envKey(name, usernameSuffix), "mailer username"); ok {
		config.Username = val
	}

	if val, ok := getEnvOrLog(envKey(name, passwordSuffix), "mailer password"); ok {
		config.Password = val
	}

	return config
}

func (m *MailerConfig) Validate() error {
	if m.Host == "" {
		return fmt.Errorf("mailer host is required")
	}
	if m.Port <= 0 || m.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", m.Port)
	}
	if m.From == "" {
		return fmt.Errorf("from address is required")
	}
	return nil
}