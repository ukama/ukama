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
	"fmt"

	"net/http"

	"testing"

	"github.com/ukama/ukama/systems/common/config"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

const ChunksPath = "/v1/chunks"

type TestConfig struct {
	config.BaseConfig
	DistributionHost string
}

var (
	cappname    = "ukamaos"
	cappversion = "1.0.1"
)

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{
		DistributionHost: "http://localhost:8098",
	}

	config.LoadConfig("integration", tConfig)
	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

// Call webhost endpoint and check response
func Test_PutChunks(t *testing.T) {
	appUrl := fmt.Sprintf("%s%s", tConfig.DistributionHost, ChunksPath)

	rest := resty.New().EnableTrace().SetDebug(tConfig.DebugMode)

	t.Run("Ping", func(t *testing.T) {
		r, err := rest.R().Get(tConfig.DistributionHost + "/ping")

		assert.NoError(t, err)
		assert.Equal(t, r.StatusCode(), http.StatusOK)
	})

	t.Run("Put", func(tt *testing.T) {
		r, err := rest.R().SetBody(map[string]interface{}{
			"store": "testdata/art"}).Put(appUrl + "/" + cappname + "/" + cappversion)

		assert.NoError(tt, err)
		log.Infof("Response: '%d'", r.StatusCode())
		assert.Equal(tt, r.StatusCode(), http.StatusOK)

	})
}
