/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/report/api-gateway/internal"
	"github.com/ukama/ukama/systems/report/api-gateway/internal/client"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	crest "github.com/ukama/ukama/systems/common/rest"
	pmocks "github.com/ukama/ukama/systems/report/api-gateway/mocks"
)

const (
	pdfEndpoint = "/v1/pdf"
)

var (
	defaultCors = cors.Config{
		AllowAllOrigins: true,
	}

	routerConfig = &RouterConfig{
		serverConf: &crest.HttpConfig{
			Cors: defaultCors,
		},
		auth: &config.Auth{
			AuthAppUrl:    "http://localhost:4455",
			AuthServerUrl: "http://localhost:4434",
			AuthAPIGW:     "http://localhost:8080",
		},
	}

	testClientSet *Clients
)

func init() {
	gin.SetMode(gin.TestMode)

	testClientSet = NewClientsSet(
		&internal.GrpcEndpoints{
			Timeout:   1 * time.Second,
			Generator: "report:9090",
		},

		&internal.HttpEndpoints{
			Timeout: 1 * time.Second,
			Files:   `http://report:3000`,
		}, true)
}

func TestRouter_PingRoute(t *testing.T) {
	var pm = &pmocks.Pdf{}
	var arc = &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		p: pm,
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_Pdf(t *testing.T) {
	t.Run("ReportFound", func(t *testing.T) {
		invoiceId := uuid.NewV4().String()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pdfEndpoint, invoiceId), nil)

		var arc = &cmocks.AuthClient{}
		pm := &pmocks.Pdf{}

		var content = []byte("some fake pdf data")

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
		pm.On("GetPdf", invoiceId).Return(content, nil)

		r := NewRouter(&Clients{
			p: pm,
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ReportNotFound", func(t *testing.T) {
		invoiceId := uuid.NewV4().String()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pdfEndpoint, invoiceId), nil)

		var arc = &cmocks.AuthClient{}
		pm := &pmocks.Pdf{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
		pm.On("GetPdf", invoiceId).Return(nil, client.ErrInvoicePDFNotFound)

		r := NewRouter(&Clients{
			p: pm,
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
