/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
	apb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
	amocks "github.com/ukama/ukama/systems/inventory/accounting/pb/gen/mocks"
	"github.com/ukama/ukama/systems/inventory/api-gateway/pkg"
	"github.com/ukama/ukama/systems/inventory/api-gateway/pkg/client"
	cpb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
	cmocks "github.com/ukama/ukama/systems/inventory/component/pb/gen/mocks"
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &pkg.HttpEndpoints{
		NodeMetrics: "localhost:8080",
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout:    1 * time.Second,
		Component:  "component:9090",
		Accounting: "accounting:9090",
	})
}

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	arc := &providers.AuthRestClient{}
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestGetComponent(t *testing.T) {
	// arrange
	var uId = uuid.NewV4()
	var cId = uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/components/"+cId.String(), nil)
	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	comp.On("Get", mock.Anything, mock.Anything).Return(&cpb.GetResponse{
		Component: &cpb.Component{
			Id:            cId.String(),
			UserId:        uId.String(),
			Inventory:     "2",
			Category:      1,
			Type:          "tower node",
			Description:   "best tower node",
			DatasheetURL:  "https://datasheet.com",
			ImagesURL:     "https://images.com",
			PartNumber:    "1234",
			Manufacturer:  "ukama",
			Managed:       "ukama",
			Warranty:      1,
			Specification: "spec",
		},
	}, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestGetComponents(t *testing.T) {
	// arrange
	var uId = uuid.NewV4()
	var cId = uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/components/user/"+uId.String(), nil)
	q := req.URL.Query()
	q.Add("category", cpb.ComponentCategory_name[1])
	req.URL.RawQuery = q.Encode()

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	comp.On("GetByUser", mock.Anything, mock.Anything, mock.Anything).Return(&cpb.GetByUserResponse{
		Components: []*cpb.Component{
			{
				Id:            cId.String(),
				UserId:        uId.String(),
				Inventory:     "2",
				Category:      1,
				Type:          "tower node",
				Description:   "best tower node",
				DatasheetURL:  "https://datasheet.com",
				ImagesURL:     "https://images.com",
				PartNumber:    "1234",
				Manufacturer:  "ukama",
				Managed:       "ukama",
				Warranty:      1,
				Specification: "spec",
			},
		},
	}, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestGetAccounting(t *testing.T) {
	var uId = uuid.NewV4()
	var aId = uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/accounting/"+aId.String(), nil)
	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	acc.On("Get", mock.Anything, mock.Anything).Return(&apb.GetResponse{
		Accounting: &apb.Accounting{
			Id:            aId.String(),
			Vat:           "10",
			Item:          "Product-1",
			UserId:        uId.String(),
			Inventory:     "1",
			OpexFee:       "100",
			EffectiveDate: "2023-01-01",
			Description:   "Product-1 description",
		},
	}, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewNewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}

func TestGetAccountings(t *testing.T) {
	// arrange
	var uId = uuid.NewV4()
	var aId = uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/accounting/user/"+uId.String(), nil)

	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	acc.On("GetByUser", mock.Anything, mock.Anything).Return(&apb.GetByUserResponse{
		Accounting: []*apb.Accounting{
			{
				Id:            aId.String(),
				Vat:           "10",
				Item:          "Product-1",
				UserId:        uId.String(),
				Inventory:     "1",
				OpexFee:       "100",
				EffectiveDate: "2023-01-01",
				Description:   "Product-1 description",
			},
		},
	}, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewNewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}

func TestSyncComponents(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/components/sync", nil)
	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	comp.On("SyncComponents", mock.Anything, mock.Anything).Return(&cpb.SyncComponentsResponse{}, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestSyncAccounting(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/accounting/sync", nil)
	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	acc.On("SyncAccounting", mock.Anything, mock.Anything).Return(&apb.SyncAcountingResponse{}, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewNewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}
