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
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	api "github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Before running test for the first time you have to create a test account in Identity manager and provide email and password for it

type TestConfig struct {
	ServiceHost string `default:"localhost:8080"`
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_LookupClientApi(t *testing.T) {

	org := "org-name"
	nodeId := ukama.NewVirtualHomeNodeId().String()
	system := "sys-name"

	client := resty.New()

	t.Run("GetOrgNotFound", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/orgs/someOrgWhichDoesnotExist")

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusNotFound, resp.StatusCode())
			assert.Contains(tt, resp.String(), "org record not found")
		}
	})

	t.Run("AddOrg", func(tt *testing.T) {
		id := uuid.NewV4().String()
		rd := api.AddOrgRequest{
			OrgName:     "org-name",
			Certificate: "helloOrg",
			Ip:          "0.0.0.0",
			OrgId:       id,
		}

		jd, err := json.Marshal(&rd)
		resp, err := client.R().
			EnableTrace().
			SetBody(jd).
			Put(getApiUrl() + "/v1/orgs/" + org)

		if err != nil {
			if assert.Error(t, err) {
				assert.Equal(tt, http.StatusCreated, resp.StatusCode())
			}
		} else {
			if assert.NoError(t, err) {
				assert.Equal(tt, http.StatusConflict, resp.StatusCode())
			}
		}

	})

	t.Run("GetOrg", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/orgs/" + org)

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
	})

	t.Run("AddNode", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Put(getApiUrl() + "/v1/orgs/" + org + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
		}
	})

	t.Run("GetNodes", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/orgs/" + org + "/nodes/" + nodeId)
		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})

	t.Run("DeleteNode", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Delete(getApiUrl() + "/v1/orgs/" + org + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}
	})

	t.Run("AddSystem", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetBody(strings.NewReader(`{ "org":"org-name", "system":"sys", "ip":"0.0.0.0", "certificate":"certs", "port":100}`)).
			Put(getApiUrl() + "/v1/orgs/" + org + "/systems/" + system)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
			fmt.Println("Response: ", resp.String())
		}
	})

	t.Run("GetSystems", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/orgs/" + org + "/systems/" + system)
		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})

	t.Run("DeleteSystems", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Delete(getApiUrl() + "/v1/orgs/" + org + "/systems/" + system)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}
	})
}

func getApiUrl() string {
	return testConf.ServiceHost
}
