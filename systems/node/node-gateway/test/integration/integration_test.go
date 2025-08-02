/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package integration

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"testing"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/ukama/ukama/systems/common/uuid"
// 	"github.com/ukama/ukama/systems/node/node-gateway/pkg/rest"

// 	"github.com/ukama/ukama/systems/common/config"

// 	log "github.com/sirupsen/logrus"
// )

// const healthApiEndpoint = "/v1/health/"

// var testConf *TestConfig

// type TestConfig struct {
// 	ServiceHost string `default:"localhost:8080"`
// }

// func init() {
// 	testConf = &TestConfig{}

// 	log.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
// 	config.LoadConfig("integration", testConf)
// 	log.Infof("Config: %+v", testConf)
// }

// func TestApiGateway_Endpoints(t *testing.T) {
// 	client := resty.New()
// 	nodeId :=  uuid.NewV4()
// 	var nt = rest.GetRunningAppsRequest{
// 	NodeId:      nodeId.String(),
// }
// 	t.Run("StoreRunningAppsInfo", func(tt *testing.T) {
// 		body, err := json.Marshal(nt)
// 		if err != nil {
// 			t.Errorf("fail to marshal request data: %v. Error: %v", nt, err)
// 		}
// 		resp, err := client.R().
// 			EnableTrace().
// 			SetBody(bytes.NewReader(body)).
// 			Post(getApiUrl() + healthApiEndpoint + "nodes/" + nt.NodeId + "/performance")
// 		assert.NoError(t, err)
// 		assert.Equal(tt, http.StatusOK, resp.StatusCode())
// 	})

// 	t.Run("GetRunningAppsInfo", func(tt *testing.T) {

// 		resp, err := client.R().
// 			EnableTrace().
// 			Get(getApiUrl() + healthApiEndpoint + "nodes/" + nt.NodeId + "/performance")
// 		assert.NoError(t, err)
// 		assert.Equal(tt, http.StatusOK, resp.StatusCode())
// 		})
// }
// func getApiUrl() string {
// 	return "http://" + testConf.ServiceHost
// }
