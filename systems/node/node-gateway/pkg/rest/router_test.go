/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"

	"github.com/tj/assert"
	crest "github.com/ukama/ukama/systems/common/rest"
	nmocks "github.com/ukama/ukama/systems/node/health/pb/gen/mocks"

	"github.com/ukama/ukama/systems/node/node-gateway/pkg"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/client"
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &crest.HttpConfig{
		Cors: defaultCors,
	},
	
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout:    1 * time.Second,
		Health: "0.0.0.0:9092",
	})
}
func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}




func Test_GetRunningsApps(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/health/nodes/60285a2a-fe1d-4261-a868-5be480075b8f/performance", nil)
	c := &nmocks.HealhtServiceClient{}
	getRunningAppsReq := &hpb.GetRunningAppsRequest{
		NodeId: "60285a2a-fe1d-4261-a868-5be480075b8f",
	}

	// Set up the mock expectations for GetRunningApps.
	c.On("GetRunningApps", mock.Anything, getRunningAppsReq).Return(
		&hpb.GetRunningAppsResponse{
			RunningApps: &hpb.App{
				Id:        "60285a2a-fe1d-4261-a868-5be480075b8f",
				NodeId:    getRunningAppsReq.NodeId,
				Timestamp: "12-12-2024",
			},
		},
		nil,
	).Once() // Use Once() to indicate that this expectation should be called once.

	// Create a new router with the mock client.
	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(c),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}


func Test_StoreRunningApps(t *testing.T) {
	chealth := &nmocks.HealhtServiceClient{}

	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(chealth),
	}, routerConfig).f.Engine()

	
	t.Run("storeRunningApps", func(t *testing.T) {
		

		data := &hpb.StoreRunningAppsInfoRequest{
			NodeId:    "60285a2a-fe1d-4261-a868-5be480075b8f",
			Timestamp: "12-12-2023",
			System: []*hpb.System{
				{
					Name:  "SystemName1",
					Value: "SystemValue1",
				},
				{
					Name:  "SystemName2",
					Value: "SystemValue2",
				},
			},
			Capps:  []*hpb.Capps{
				{
					Name:   "CappsName1",
					Tag:    "CappsTag1",
					Resources: []*hpb.Resource{
						{
							Name:  "ResourceName1",
							Value: "ResourceValue1",
						},
						{
							Name:  "ResourceName2",
							Value: "ResourceValue2",
						},
					},
				},
				{
					Name:   "CappsName2",
					Tag:    "CappsTag2",
					Resources: []*hpb.Resource{
						{
							Name:  "ResourceName3",
							Value: "ResourceValue3",
						},
						{
							Name:  "ResourceName4",
							Value: "ResourceValue4",
						},
					},
				},
			},
		}
		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/health/nodes/60285a2a-fe1d-4261-a868-5be480075b8f/performance", bytes.NewReader(jdata))
		assert.NoError(t, err)
		

		chealth.On("StoreRunningAppsInfo", mock.Anything, data).Return(&hpb.StoreRunningAppsInfoResponse{
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		chealth.AssertExpectations(t)
	})
}