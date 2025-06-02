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
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/common/config"

	log "github.com/sirupsen/logrus"
)

type TestConfig struct {
	config.BaseConfig
	ServiceHost string
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{
		ServiceHost: "http://localhost:8080",
	}

	config.LoadConfig("integration", tConfig)
	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

// Call webhost endpoint and check response
func Test_Metrics(t *testing.T) {
	rest := resty.New().EnableTrace().SetDebug(tConfig.DebugMode)

	t.Run("Ping", func(t *testing.T) {
		r, err := rest.R().Get(tConfig.ServiceHost + "/ping")
		assert.NoError(t, err)
		assert.Equal(t, r.StatusCode(), http.StatusOK)
	})

}
