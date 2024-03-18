/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc    `default:"{}"`
	Metrics          *uconf.Metrics `default:"{}"`
	Service          *uconf.Service
	RepoUrl          string `default:"https://github.com/ukama/networks.git"`
	GitUsername      string `default:"ukama"`
	Token            string `default:"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"`
	RootConfigURL    string `default:"https://raw.githubusercontent.com/ukama/networks/main/root.json"`
}

func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
	}
}
