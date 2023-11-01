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
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/rest"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
)

const apiEndpoint = "/v1/api"

var testConf *TestConfig

type TestConfig struct {
	ServiceHost string `default:"localhost:8080"`
}

func init() {
	testConf = &TestConfig{}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	log.Infof("Config: %+v", testConf)
}
