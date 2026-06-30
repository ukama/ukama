/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tj/assert"

	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

var serveRouterConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: cors.Config{AllowAllOrigins: true},
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

func noAuth(c *gin.Context, s string) error { return nil }

func newServeRouter(role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	clients := &Clients{
		Manager: &fakeManager{},
		Member:  &fakeMember{role: role},
	}
	return NewRouter(clients, serveRouterConfig, noAuth).f.Engine()
}

const testOpID = "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"

func TestRoutes_StartOperation(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"type":"RestartNode","system":"node","resource_key":"node:abc"}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/operations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	newServeRouter("ROLE_OWNER").ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRoutes_GetByResource(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/operations?resource_key=node:abc", nil)

	newServeRouter("ROLE_OWNER").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoutes_GetOperation(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/operations/"+testOpID, nil)

	newServeRouter("ROLE_OWNER").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoutes_MarkRunning(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"fencing_token":3}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/operations/"+testOpID+"/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	newServeRouter("ROLE_OWNER").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoutes_ForceUnlock(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"user_id":"` + testOpID + `","reason":"stuck"}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/operations/"+testOpID+"/force-unlock", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	newServeRouter("ROLE_OWNER").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoutes_StartOperation_MissingFields(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/operations", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	newServeRouter("ROLE_OWNER").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
