/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Grpc              *config.Grpc    `default:"{}"`
	Metrics           *config.Metrics `default:"{}"`
	Service           *config.Service
	Timeout           time.Duration `default:"3s"`
}

func NewConfig(name string) *Config {
	// Sanitize name
	name = strings.ReplaceAll(name, "-", "_")

	return &Config{
		Service: config.LoadServiceHostConfig(name),
	}
}
