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
// 	"github.com/ukama/ukama/systems/node/api-gateway/pkg/rest"

// 	"github.com/ukama/ukama/systems/common/config"

// 	log "github.com/sirupsen/logrus"
// )

// const controllerApiEndpoint = "/v1/controllers/"

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

// 	var nt = rest.RestartNodeRequest{
// 	NodeId:      nodeId.String(),
// }
// 	t.Run("RestartNode", func(tt *testing.T) {
// 		body, err := json.Marshal(nt)
// 		if err != nil {
// 			t.Errorf("fail to marshal request data: %v. Error: %v", nt, err)
// 		}
// 		resp, err := client.R().
// 			EnableTrace().
// 			SetBody(bytes.NewReader(body)).
// 			Post(getApiUrl() + controllerApiEndpoint + "nodes/" + nt.NodeId + "/restart")
// 		assert.NoError(t, err)
// 		assert.Equal(tt, http.StatusOK, resp.StatusCode())
// 	})

// 	t.Run("RestartSite", func(tt *testing.T) {
// 		netId:= uuid.NewV4()

// 		restartSiteReq:= rest.RestartSiteRequest{
// 			SiteName: "site1",
// 			NetworkId: netId.String(),
// 		}
// 		body,err := json.Marshal(restartSiteReq)
// 		if err != nil {
// 			t.Errorf("fail to marshal request data: %v. Error: %v", restartSiteReq, err)
// 		}
// 		resp, err := client.R().

// 			SetBody(bytes.NewReader(body)).
// 			Post(getApiUrl() + controllerApiEndpoint + "networks/"+restartSiteReq.NetworkId+"/sites/"+ restartSiteReq.SiteName+"/restart")
// 		assert.NoError(t, err)
// 		assert.Equal(tt, http.StatusOK, resp.StatusCode())
// 	})
// 	t.Run("RestartNodes", func(tt *testing.T) {
// 		netId:= uuid.NewV4()

// 		restartNodesReq:= rest.RestartNodesRequest{
// 			NetworkId: netId.String(),
// 			NodeIds: []string{nodeId.String()},
// 		}
// 		body,err := json.Marshal(restartNodesReq)
// 		if err != nil {
// 			t.Errorf("fail to marshal request data: %v. Error: %v", restartNodesReq, err)
// 		}
// 		resp, err := client.R().

// 			SetBody(bytes.NewReader(body)).
// 			Post(getApiUrl() + controllerApiEndpoint + "networks/"+restartNodesReq.NetworkId+"/restart-nodes")
// 		assert.NoError(t, err)
// 		assert.Equal(tt, http.StatusOK, resp.StatusCode())

// 	})

// }
// func getApiUrl() string {
// 	return "http://localhost:8080"
// }
