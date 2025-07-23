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
	testCategory              = "access"
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
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, testHTTPStatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testPingResponse)
}

func TestRouter_GetComponent(t *testing.T) {
	// arrange
	componentId := uuid.NewV4().String()
	userId := uuid.NewV4().String()

	req := GetRequest{
		Uuid: componentId,
	}

	jReq, err := json.Marshal(req)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/components/"+componentId, bytes.NewReader(jReq))

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	compReq := &cpb.GetRequest{
		Id: componentId,
	}

	compResp := &cpb.GetResponse{
		Component: &cpb.Component{
			Id:            componentId,
			UserId:        userId,
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
	}

	comp.On("Get", mock.Anything, compReq).Return(compResp, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestRouter_GetComponents(t *testing.T) {
	// arrange
	userId := uuid.NewV4().String()
	componentId := uuid.NewV4().String()

	req := GetComponents{
		UserId:   userId,
		Category: testCategory,
	}

	jReq, err := json.Marshal(req)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/components/user/"+userId, bytes.NewReader(jReq))
	q := hreq.URL.Query()
	q.Add("category", req.Category)
	hreq.URL.RawQuery = q.Encode()

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	compReq := &cpb.GetByUserRequest{
		UserId:   userId,
		Category: testCategory,
	}

	compResp := &cpb.GetByUserResponse{
		Components: []*cpb.Component{
			{
				Id:            componentId,
				UserId:        userId,
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
	}

	comp.On("GetByUser", mock.Anything, compReq).Return(compResp, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestRouter_GetAccounting(t *testing.T) {
	// arrange
	accountingId := uuid.NewV4().String()
	userId := uuid.NewV4().String()

	req := GetRequest{
		Uuid: accountingId,
	}

	jReq, err := json.Marshal(req)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/accounting/"+accountingId, bytes.NewReader(jReq))

	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	accReq := &apb.GetRequest{
		Id: accountingId,
	}

	accResp := &apb.GetResponse{
		Accounting: &apb.Accounting{
			Id:            accountingId,
			Vat:           testVat,
			Item:          testItem,
			UserId:        userId,
			Inventory:     testAccountingInventoryID,
			OpexFee:       testOpexFee,
			EffectiveDate: testEffectiveDate,
			Description:   testAccountingDescription,
		},
	}

	acc.On("Get", mock.Anything, accReq).Return(accResp, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}

func TestRouter_GetAccountings(t *testing.T) {
	// arrange
	userId := uuid.NewV4().String()
	accountingId := uuid.NewV4().String()

	req := GetAccounts{
		UserId: userId,
	}

	jReq, err := json.Marshal(req)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/accounting/user/"+userId, bytes.NewReader(jReq))

	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	accReq := &apb.GetByUserRequest{
		UserId: userId,
	}

	accResp := &apb.GetByUserResponse{
		Accounting: []*apb.Accounting{
			{
				Id:            accountingId,
				Vat:           testVat,
				Item:          testItem,
				UserId:        userId,
				Inventory:     testAccountingInventoryID,
				OpexFee:       testOpexFee,
				EffectiveDate: testEffectiveDate,
				Description:   testAccountingDescription,
			},
		},
	}

	acc.On("GetByUser", mock.Anything, accReq).Return(accResp, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}

func TestRouter_SyncComponents(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("PUT", "/v1/components/sync", nil)
	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	compResp := &cpb.SyncComponentsResponse{}

	comp.On("SyncComponents", mock.Anything, mock.Anything).Return(compResp, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestRouter_SyncAccounting(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("PUT", "/v1/accounting/sync", nil)
	arc := &providers.AuthRestClient{}
	acc := &amocks.AccountingServiceClient{}

	accResp := &apb.SyncAcountingResponse{}

	acc.On("SyncAccounting", mock.Anything, mock.Anything).Return(accResp, nil)

	r := NewRouter(&Clients{
		Accounting: client.NewAccountingInventoryFromClient(acc),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	acc.AssertExpectations(t)
}

func TestRouter_ListComponents(t *testing.T) {
	// arrange
	componentId := uuid.NewV4().String()
	userId := uuid.NewV4().String()
	partNumber := "TEST-123"
	category := "access"

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/components", nil)

	// Add query parameters
	q := hreq.URL.Query()
	q.Add("id", componentId)
	q.Add("user_id", userId)
	q.Add("part_number", partNumber)
	q.Add("category", category)
	hreq.URL.RawQuery = q.Encode()

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	compReq := &cpb.ListRequest{
		Id:         componentId,
		UserId:     userId,
		PartNumber: partNumber,
		Category:   category,
	}

	compResp := &cpb.ListResponse{
		Components: []*cpb.Component{
			{
				Id:            componentId,
				UserId:        userId,
				Inventory:     testInventoryID,
				Category:      category,
				Type:          testComponentType,
				Description:   testDescription,
				DatasheetURL:  testDatasheetURL,
				ImagesURL:     testImagesURL,
				PartNumber:    partNumber,
				Manufacturer:  testManufacturer,
				Managed:       testManaged,
				Warranty:      testWarranty,
				Specification: testSpecification,
			},
		},
	}

	comp.On("List", mock.Anything, compReq).Return(compResp, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestRouter_ListComponents_WithDefaultCategory(t *testing.T) {
	// arrange
	componentId := uuid.NewV4().String()
	userId := uuid.NewV4().String()
	partNumber := "TEST-456"
	// category is not provided, should default to "all"

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/components", nil)

	// Add query parameters without category
	q := hreq.URL.Query()
	q.Add("id", componentId)
	q.Add("user_id", userId)
	q.Add("part_number", partNumber)
	hreq.URL.RawQuery = q.Encode()

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	compReq := &cpb.ListRequest{
		Id:         componentId,
		UserId:     userId,
		PartNumber: partNumber,
		Category:   "all", // default value
	}

	compResp := &cpb.ListResponse{
		Components: []*cpb.Component{
			{
				Id:            componentId,
				UserId:        userId,
				Inventory:     testInventoryID,
				Category:      "all",
				Type:          testComponentType,
				Description:   testDescription,
				DatasheetURL:  testDatasheetURL,
				ImagesURL:     testImagesURL,
				PartNumber:    partNumber,
				Manufacturer:  testManufacturer,
				Managed:       testManaged,
				Warranty:      testWarranty,
				Specification: testSpecification,
			},
		},
	}

	comp.On("List", mock.Anything, compReq).Return(compResp, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestRouter_ListComponents_EmptyResponse(t *testing.T) {
	// arrange
	componentId := uuid.NewV4().String()
	userId := uuid.NewV4().String()
	partNumber := "TEST-789"
	category := "power"

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/components", nil)

	// Add query parameters
	q := hreq.URL.Query()
	q.Add("id", componentId)
	q.Add("user_id", userId)
	q.Add("part_number", partNumber)
	q.Add("category", category)
	hreq.URL.RawQuery = q.Encode()

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	compReq := &cpb.ListRequest{
		Id:         componentId,
		UserId:     userId,
		PartNumber: partNumber,
		Category:   category,
	}

	compResp := &cpb.ListResponse{
		Components: []*cpb.Component{}, // empty response
	}

	comp.On("List", mock.Anything, compReq).Return(compResp, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestRouter_ListComponents_MultipleComponents(t *testing.T) {
	// arrange
	componentId1 := uuid.NewV4().String()
	componentId2 := uuid.NewV4().String()
	userId := uuid.NewV4().String()
	partNumber := "TEST-MULTI"
	category := "switch"

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/components", nil)

	// Add query parameters
	q := hreq.URL.Query()
	q.Add("id", componentId1)
	q.Add("user_id", userId)
	q.Add("part_number", partNumber)
	q.Add("category", category)
	hreq.URL.RawQuery = q.Encode()

	arc := &providers.AuthRestClient{}
	comp := &cmocks.ComponentServiceClient{}

	compReq := &cpb.ListRequest{
		Id:         componentId1,
		UserId:     userId,
		PartNumber: partNumber,
		Category:   category,
	}

	compResp := &cpb.ListResponse{
		Components: []*cpb.Component{
			{
				Id:            componentId1,
				UserId:        userId,
				Inventory:     testInventoryID,
				Category:      category,
				Type:          testComponentType,
				Description:   testDescription,
				DatasheetURL:  testDatasheetURL,
				ImagesURL:     testImagesURL,
				PartNumber:    partNumber,
				Manufacturer:  testManufacturer,
				Managed:       testManaged,
				Warranty:      testWarranty,
				Specification: testSpecification,
			},
			{
				Id:            componentId2,
				UserId:        userId,
				Inventory:     testInventoryID,
				Category:      category,
				Type:          "switch node",
				Description:   "Another switch component",
				DatasheetURL:  testDatasheetURL,
				ImagesURL:     testImagesURL,
				PartNumber:    partNumber,
				Manufacturer:  testManufacturer,
				Managed:       testManaged,
				Warranty:      testWarranty,
				Specification: testSpecification,
			},
		},
	}

	comp.On("List", mock.Anything, compReq).Return(compResp, nil)

	r := NewRouter(&Clients{
		Component: client.NewComponentInventoryFromClient(comp),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	comp.AssertExpectations(t)
}

func TestRouter_ListComponents_AllCategories(t *testing.T) {
	// Test all valid categories
	categories := []string{"all", "access", "backhaul", "power", "switch", "spectrum"}

	for _, category := range categories {
		t.Run("Category_"+category, func(t *testing.T) {
			// arrange
			componentId := uuid.NewV4().String()
			userId := uuid.NewV4().String()
			partNumber := "TEST-" + category

			w := httptest.NewRecorder()
			hreq, _ := http.NewRequest("GET", "/v1/components", nil)

			// Add query parameters
			q := hreq.URL.Query()
			q.Add("id", componentId)
			q.Add("user_id", userId)
			q.Add("part_number", partNumber)
			q.Add("category", category)
			hreq.URL.RawQuery = q.Encode()

			arc := &providers.AuthRestClient{}
			comp := &cmocks.ComponentServiceClient{}

			compReq := &cpb.ListRequest{
				Id:         componentId,
				UserId:     userId,
				PartNumber: partNumber,
				Category:   category,
			}

			compResp := &cpb.ListResponse{
				Components: []*cpb.Component{
					{
						Id:            componentId,
						UserId:        userId,
						Inventory:     testInventoryID,
						Category:      category,
						Type:          testComponentType,
						Description:   testDescription,
						DatasheetURL:  testDatasheetURL,
						ImagesURL:     testImagesURL,
						PartNumber:    partNumber,
						Manufacturer:  testManufacturer,
						Managed:       testManaged,
						Warranty:      testWarranty,
						Specification: testSpecification,
					},
				},
			}

			comp.On("List", mock.Anything, compReq).Return(compResp, nil)

			r := NewRouter(&Clients{
				Component: client.NewComponentInventoryFromClient(comp),
			}, routerConfig, arc.MockAuthenticateUser).f.Engine()

			// act
			r.ServeHTTP(w, hreq)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			comp.AssertExpectations(t)
		})
	}
}
