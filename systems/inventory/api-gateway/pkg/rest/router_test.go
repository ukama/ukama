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

// Test data constants
const (
	testInventoryID           = "1"
	testCategory              = "ACCESS"
	testComponentType         = "tower node"
	testDescription           = "best tower node"
	testDatasheetURL          = "https://datasheet.com"
	testImagesURL             = "https://images.com"
	testPartNumber            = "1234"
	testManufacturer          = "ukama"
	testManaged               = "ukama"
	testWarranty              = 1
	testSpecification         = "spec"
	testVat                   = "10"
	testItem                  = "Product-1"
	testAccountingInventoryID = "1"
	testOpexFee               = "100"
	testEffectiveDate         = "2023-01-01"
	testAccountingDescription = "Product-1 description"
	testPingResponse          = "pong"
	testHTTPStatusOK          = 200
)

// Test endpoints
const (
	testComponentEndpoint      = "/v1/components"
	testAccountingEndpoint     = "/v1/accounting"
	testPingEndpoint           = "/ping"
	testSyncComponentsEndpoint = "/v1/components/sync"
	testSyncAccountingEndpoint = "/v1/accounting/sync"
)

// Test server configurations
const (
	testNodeMetricsEndpoint = "localhost:8080"
	testAuthAppURL          = "http://localhost:4455"
	testAuthServerURL       = "http://localhost:4434"
	testAuthAPIGW           = "http://localhost:8080"
	testComponentGRPC       = "component:9090"
	testAccountingGRPC      = "accounting:9090"
	testTimeout             = 1 * time.Second
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &pkg.HttpEndpoints{
		NodeMetrics: testNodeMetricsEndpoint,
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    testAuthAppURL,
		AuthServerUrl: testAuthServerURL,
		AuthAPIGW:     testAuthAPIGW,
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout:    testTimeout,
		Component:  testComponentGRPC,
		Accounting: testAccountingGRPC,
	})
}

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	arc := &providers.AuthRestClient{}
	req, _ := http.NewRequest("GET", testPingEndpoint, nil)
	r := NewRouter(testClientSet, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, testHTTPStatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testPingResponse)
}

func TestGetComponent(t *testing.T) {
	// arrange
	var uId = uuid.NewV4()
	var cId = uuid.NewV4()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", testComponentEndpoint+"/"+cId.String(), nil)
	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	comp.On("Get", mock.Anything, mock.Anything).Return(&cpb.GetResponse{
		Component: &cpb.Component{
			Id:            cId.String(),
			UserId:        uId.String(),
			Inventory:     testInventoryID,
			Category:      testCategory,
			Type:          testComponentType,
			Description:   testDescription,
			DatasheetURL:  testDatasheetURL,
			ImagesURL:     testImagesURL,
			PartNumber:    testPartNumber,
			Manufacturer:  testManufacturer,
			Managed:       testManaged,
			Warranty:      testWarranty,
			Specification: testSpecification,
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
	req, _ := http.NewRequest("GET", testComponentEndpoint+"/user/"+uId.String(), nil)
	q := req.URL.Query()
	q.Add("category", testCategory)
	req.URL.RawQuery = q.Encode()

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	comp.On("GetByUser", mock.Anything, mock.Anything, mock.Anything).Return(&cpb.GetByUserResponse{
		Components: []*cpb.Component{
			{
				Id:            cId.String(),
				UserId:        uId.String(),
				Inventory:     testInventoryID,
				Category:      testCategory,
				Type:          testComponentType,
				Description:   testDescription,
				DatasheetURL:  testDatasheetURL,
				ImagesURL:     testImagesURL,
				PartNumber:    testPartNumber,
				Manufacturer:  testManufacturer,
				Managed:       testManaged,
				Warranty:      testWarranty,
				Specification: testSpecification,
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
	req, _ := http.NewRequest("GET", testAccountingEndpoint+"/"+aId.String(), nil)
	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	acc.On("Get", mock.Anything, mock.Anything).Return(&apb.GetResponse{
		Accounting: &apb.Accounting{
			Id:            aId.String(),
			Vat:           testVat,
			Item:          testItem,
			UserId:        uId.String(),
			Inventory:     testAccountingInventoryID,
			OpexFee:       testOpexFee,
			EffectiveDate: testEffectiveDate,
			Description:   testAccountingDescription,
		},
	}, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewAccountingInventoryFromClient(acc),
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
	req, _ := http.NewRequest("GET", testAccountingEndpoint+"/user/"+uId.String(), nil)

	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	acc.On("GetByUser", mock.Anything, mock.Anything).Return(&apb.GetByUserResponse{
		Accounting: []*apb.Accounting{
			{
				Id:            aId.String(),
				Vat:           testVat,
				Item:          testItem,
				UserId:        uId.String(),
				Inventory:     testAccountingInventoryID,
				OpexFee:       testOpexFee,
				EffectiveDate: testEffectiveDate,
				Description:   testAccountingDescription,
			},
		},
	}, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}

func TestSyncComponents(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", testSyncComponentsEndpoint, nil)
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
	req, _ := http.NewRequest("PUT", testSyncAccountingEndpoint, nil)
	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	acc.On("SyncAccounting", mock.Anything, mock.Anything).Return(&apb.SyncAcountingResponse{}, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}
