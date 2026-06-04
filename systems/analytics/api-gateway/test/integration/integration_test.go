//go:build integration
// +build integration

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package integration

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
)

type TestConfig struct {
	ServiceHost string `default:"localhost:8080"`
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	config.LoadConfig("integration", testConf)
	log.Infof("Config: %+v", testConf)
}

func getApiUrl() string {
	return "http://" + testConf.ServiceHost
}

// Integration tests run against a live analytics system, e.g.:
//
// func Test_AnalyticsApi(t *testing.T) {
//	client := resty.New()
//
//	t.Run("BusinessHome", func(tt *testing.T) {
//		resp, err := client.R().
//			EnableTrace().
//			Get(getApiUrl() + "/v1/analytics/business/home?period=week")
//
//		if assert.NoError(tt, err) {
//			assert.Equal(tt, http.StatusOK, resp.StatusCode())
//		}
//	})
//
//	t.Run("CustomerOverview", func(tt *testing.T) {
//		resp, err := client.R().
//			EnableTrace().
//			Get(getApiUrl() + "/v1/analytics/customers/overview")
//
//		if assert.NoError(tt, err) {
//			assert.Equal(tt, http.StatusOK, resp.StatusCode())
//		}
//	})
//
//	t.Run("NetworkOverview", func(tt *testing.T) {
//		resp, err := client.R().
//			EnableTrace().
//			Get(getApiUrl() + "/v1/analytics/network/overview")
//
//		if assert.NoError(tt, err) {
//			assert.Equal(tt, http.StatusOK, resp.StatusCode())
//		}
//	})
// }
