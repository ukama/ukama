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

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/rest"

	log "github.com/sirupsen/logrus"
)

const notifyApiEndpoint = "/v1/notifications"

var testConf *TestConfig
var nodeId = ukama.NewVirtualHomeNodeId().String()
var nt = rest.AddNotificationReq{
	NodeId:      nodeId,
	Severity:    "high",
	Type:        "event",
	ServiceName: "noded",
	Status:      8300,
	Time:        uint32(time.Now().Unix()),
	Description: "Some random alert",
	Details:     `{"reason": "testing", "component":"router_test"}`,
}

type TestConfig struct {
	ServiceHost string `default:"localhost:8080"`
}

func init() {
	testConf = &TestConfig{}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	log.Infof("Config: %+v", testConf)
}

func TestNodeGateway_Endpoints(t *testing.T) {
	client := resty.New()

	t.Run("NotificationNotFound", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + notifyApiEndpoint + "/someNotificationWhichDoesnotExist")

		assert.Error(t, err)
		assert.Equal(tt, http.StatusNotFound, resp.StatusCode())
	})

	t.Run("AddNotification", func(tt *testing.T) {
		body, err := json.Marshal(nt)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nt, err)
		}

		resp, err := client.R().
			EnableTrace().
			SetBody(bytes.NewReader(body)).
			Post(getApiUrl() + notifyApiEndpoint)

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusCreated, resp.StatusCode())
	})

	t.Run("ListAllNotifications", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + notifyApiEndpoint)

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
	})

	t.Run("ListServiceNotifications", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + notifyApiEndpoint + "?service_name=noded")

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
	})

	t.Run("ListEventNotifications", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + notifyApiEndpoint + "?notification_type=event")

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
	})

	t.Run("PurgeNotifications", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Delete(getApiUrl() + notifyApiEndpoint)

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
	})
}

func getApiUrl() string {
	return testConf.ServiceHost
}
