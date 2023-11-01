//go:build integration
// +build integration

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package integration

import (
	"github.com/ukama/ukama/systems/common/config"

	log "github.com/sirupsen/logrus"
)

// Before running test for the first time you have to create a test account in Identity manager and provide email and password for it

type TestConfig struct {
	ServiceHost string `default:"localhost:8080"`
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	log.Infof("Config: %+v", testConf)
}
