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
	"strconv"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/client"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	bmocks "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen/mocks"
	ppb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	pmocks "github.com/ukama/ukama/systems/data-plan/package/pb/gen/mocks"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	rmocks "github.com/ukama/ukama/systems/data-plan/rate/pb/gen/mocks"
)

// Test data constants
const (
	testCountry     = "USA"
	testProvider    = "ABC"
	testSimType     = "ukama_data"
	testApn         = "ukama.tel"
	testVpmn        = "TTC"
	testPackageName = "Test Package"
	testUpdatedName = "Updated Package Name"
	testCurrency    = "USD"
	testDataUnit    = "MegaBytes"
	testVoiceUnit   = "seconds"
	testPackageType = "postpaid"
)

// Test data values
var (
	testMarkupValue     = 10.0
	testMarkupHistory1  = 5.5
	testDataRate        = 0.0014
	testSmsMoRate       = 0.0100
	testSmsMtRate       = 0.0001
	testImsiValue       = int64(1)
	testSmsVolume       = int64(100)
	testDataVolume      = int64(1024)
	testVoiceVolume     = int64(1000)
	testDuration        = uint64(30)
	testAmount          = 10.50
	testMarkup          = 5.0
	testOverdraft       = 0.0
	testTrafficPolicy   = uint32(1)
	testEffectiveAt     = "2023-10-10"
	testEffectiveAtTime = "2023-10-12T07:20:50.52Z"
	testEndAtTime       = "2043-10-12T07:20:50.52Z"
	testFromTime        = "2022-10-12T07:20:50.52Z"
	testToTime          = "2023-10-12T07:20:50.52Z"
	testFromDate        = "2023-04-01T00:00:00Z"
	testToDate          = "2023-05-01T00:00:00Z"
	testCreatedAt1      = "2021-11-12T11:45:26.371Z"
	testCreatedAt2      = "2022-11-12T11:45:26.371Z"
	testDeletedAt       = "2022-11-12T11:45:26.371Z"
	testNetworks        = []string{"network1", "network2"}
	testFileURL         = "https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/base-rate/template/template.csv"
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
	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	testClientSet = &Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}
}

func TestRouter_PingRoute(t *testing.T) {
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

func TestRouter_GetRates(t *testing.T) {
	ownerId := uuid.NewV4().String()
	req := GetRateRequest{
		UserId:   ownerId,
		Country:  testCountry,
		Provider: testProvider,
		To:       time.Now().UTC().Format(time.RFC3339),
		From:     time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		SimType:  testSimType,
	}

	w := httptest.NewRecorder()

	hreq, _ := http.NewRequest("GET", "/v1/rates/users/"+ownerId+"/rate", nil)
	q := hreq.URL.Query()
	q.Add("country", req.Country)
	q.Add("provider", req.Provider)
	q.Add("to", req.To)
	q.Add("from", req.From)
	q.Add("sim_type", req.SimType)
	hreq.URL.RawQuery = q.Encode()

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.GetRateRequest{
		OwnerId:  req.UserId,
		Country:  req.Country,
		Provider: req.Provider,
		To:       req.To,
		From:     req.From,
		SimType:  req.SimType,
	}

	pResp := &rpb.GetRateResponse{
		Rates: []*bpb.Rate{
			{
				X2G:         true,
				X3G:         true,
				Country:     req.Country,
				Data:        testDataRate,
				EffectiveAt: testEffectiveAt,
				Imsi:        testImsiValue,
				Lte:         true,
				SimType:     req.SimType,
				SmsMo:       testSmsMoRate,
				SmsMt:       testSmsMtRate,
				Vpmn:        testVpmn,
			},
		},
	}

	m.On("GetRate", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetUserMarkup(t *testing.T) {
	ownerId := uuid.NewV4().String()
	req := GetMarkupRequest{
		OwnerId: ownerId,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/markup/users/"+ownerId, nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.GetMarkupRequest{
		OwnerId: req.OwnerId,
	}

	pResp := &rpb.GetMarkupResponse{
		OwnerId: req.OwnerId,
		Markup:  testMarkupValue,
	}

	m.On("GetMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), ownerId)
	m.AssertExpectations(t)
}

func TestRouter_DeleteUserMarkup(t *testing.T) {
	ownerId := uuid.NewV4().String()
	req := GetMarkupRequest{
		OwnerId: ownerId,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("DELETE", "/v1/markup/users/"+ownerId, nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.DeleteMarkupRequest{
		OwnerId: req.OwnerId,
	}

	pResp := &rpb.DeleteMarkupResponse{}

	m.On("DeleteMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_SetUserMarkup(t *testing.T) {
	ownerId := uuid.NewV4().String()
	req := SetMarkupRequest{
		OwnerId: ownerId,
		Markup:  testMarkupValue,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("POST", "/v1/markup/"+strconv.FormatFloat(req.Markup, 'f', 'g', 64)+"/users/"+ownerId, nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.UpdateMarkupRequest{
		OwnerId: req.OwnerId,
		Markup:  req.Markup,
	}

	pResp := &rpb.UpdateMarkupResponse{}

	m.On("UpdateMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetUserMarkupHistory(t *testing.T) {
	ownerId := uuid.NewV4().String()
	req := GetMarkupRequest{
		OwnerId: ownerId,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/markup/users/"+ownerId+"/history", nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.GetMarkupHistoryRequest{
		OwnerId: req.OwnerId,
	}

	pResp := &rpb.GetMarkupHistoryResponse{
		OwnerId: req.OwnerId,
		MarkupRates: []*rpb.MarkupRates{
			{
				CreatedAt: testCreatedAt1,
				DeletedAt: testDeletedAt,
				Markup:    testMarkupHistory1,
			},
			{
				CreatedAt: testCreatedAt2,
				Markup:    testMarkupValue,
			},
		},
	}

	m.On("GetMarkupHistory", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), ownerId)
	assert.Contains(t, w.Body.String(), pResp.MarkupRates[0].CreatedAt)
	assert.Contains(t, w.Body.String(), pResp.MarkupRates[1].CreatedAt)
	m.AssertExpectations(t)
}

func TestRouter_SetDefaultMarkup(t *testing.T) {
	req := SetDefaultMarkupRequest{
		Markup: testMarkupValue,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("POST", "/v1/markup/"+strconv.FormatFloat(req.Markup, 'f', 'g', 64)+"/default", nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.UpdateDefaultMarkupRequest{
		Markup: req.Markup,
	}

	pResp := &rpb.UpdateDefaultMarkupResponse{}

	m.On("UpdateDefaultMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetDefaultMarkup(t *testing.T) {
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/markup/default", nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.GetDefaultMarkupRequest{}

	pResp := &rpb.GetDefaultMarkupResponse{
		Markup: testMarkupValue,
	}

	m.On("GetDefaultMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), strconv.FormatFloat(pResp.Markup, 'f', -1, 64))
	m.AssertExpectations(t)
}

func TestRouter_GetDefaultMarkupHistory(t *testing.T) {
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/markup/default/history", nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.GetDefaultMarkupHistoryRequest{}

	pResp := &rpb.GetDefaultMarkupHistoryResponse{
		MarkupRates: []*rpb.MarkupRates{
			{
				CreatedAt: testCreatedAt1,
				DeletedAt: testDeletedAt,
				Markup:    testMarkupHistory1,
			},
			{
				CreatedAt: testCreatedAt2,
				Markup:    testMarkupValue,
			},
		},
	}

	m.On("GetDefaultMarkupHistory", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), pResp.MarkupRates[0].CreatedAt)
	assert.Contains(t, w.Body.String(), pResp.MarkupRates[1].CreatedAt)
	m.AssertExpectations(t)
}

func TestRouter_GetBaseRatesById(t *testing.T) {
	id := uuid.NewV4()
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/baserates/"+id.String(), nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &bpb.GetBaseRatesByIdRequest{
		Uuid: id.String(),
	}

	pResp := &bpb.GetBaseRatesByIdResponse{
		Rate: &bpb.Rate{
			Uuid: id.String(),
		},
	}

	b.On("GetBaseRatesById", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), id.String())
	m.AssertExpectations(t)
}

func TestRouter_UploadBaseRates(t *testing.T) {
	ureq := UploadBaseRatesRequest{
		FileURL:     testFileURL,
		EffectiveAt: testEffectiveAtTime,
		EndAt:       testEndAtTime,
		SimType:     testSimType,
	}

	jreq, err := json.Marshal(&ureq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("POST", "/v1/baserates/upload", bytes.NewReader(jreq))

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &bpb.UploadBaseRatesRequest{
		FileURL:     ureq.FileURL,
		EffectiveAt: ureq.EffectiveAt,
		SimType:     ureq.SimType,
		EndAt:       ureq.EndAt,
	}

	pResp := &bpb.UploadBaseRatesResponse{
		Rate: []*bpb.Rate{
			{
				Uuid: uuid.NewV4().String(),
			},
		},
	}

	b.On("UploadBaseRates", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetBaseRates(t *testing.T) {
	t.Run("ByCountry", func(t *testing.T) {
		ureq := GetBaseRatesByCountryRequest{
			Country:  testCountry,
			Provider: testProvider,
			SimType:  testSimType,
		}

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("GET", "/v1/baserates", nil)
		q := hreq.URL.Query()
		q.Add("country", ureq.Country)
		q.Add("provider", ureq.Provider)
		q.Add("sim_type", ureq.SimType)
		hreq.URL.RawQuery = q.Encode()

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &bpb.GetBaseRatesByCountryRequest{
			Country:  ureq.Country,
			Provider: ureq.Provider,
			SimType:  ureq.SimType,
		}

		pResp := &bpb.GetBaseRatesResponse{
			Rates: []*bpb.Rate{
				{
					Uuid: uuid.NewV4().String(),
				},
			},
		}

		b.On("GetBaseRatesByCountry", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("HistoryByCountry", func(t *testing.T) {
		ureq := GetBaseRatesByCountryRequest{
			Country:  testCountry,
			Provider: testProvider,
			SimType:  testSimType,
		}

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("GET", "/v1/baserates/history", nil)
		q := hreq.URL.Query()
		q.Add("country", ureq.Country)
		q.Add("provider", ureq.Provider)
		q.Add("sim_type", ureq.SimType)
		hreq.URL.RawQuery = q.Encode()

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &bpb.GetBaseRatesByCountryRequest{
			Country:  ureq.Country,
			Provider: ureq.Provider,
			SimType:  ureq.SimType,
		}

		pResp := &bpb.GetBaseRatesResponse{
			Rates: []*bpb.Rate{
				{
					Uuid: uuid.NewV4().String(),
				},
			},
		}

		b.On("GetBaseRatesHistoryByCountry", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("ByCountryForPeriod", func(t *testing.T) {
		ureq := GetBaseRatesForPeriodRequest{
			Country:  testCountry,
			Provider: testProvider,
			SimType:  testSimType,
			To:       testToTime,
			From:     testFromTime,
		}

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("GET", "/v1/baserates/period", nil)
		q := hreq.URL.Query()
		q.Add("country", ureq.Country)
		q.Add("provider", ureq.Provider)
		q.Add("sim_type", ureq.SimType)
		q.Add("to", ureq.To)
		q.Add("from", ureq.From)
		hreq.URL.RawQuery = q.Encode()

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &bpb.GetBaseRatesByPeriodRequest{
			Country:  ureq.Country,
			Provider: ureq.Provider,
			SimType:  ureq.SimType,
			From:     ureq.From,
			To:       ureq.To,
		}

		pResp := &bpb.GetBaseRatesResponse{
			Rates: []*bpb.Rate{
				{
					Uuid: uuid.NewV4().String(),
				},
			},
		}

		b.On("GetBaseRatesForPeriod", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
}

func TestRouter_Package(t *testing.T) {
	t.Run("GetPackage", func(t *testing.T) {
		ureq := PackagesRequest{
			Uuid: uuid.NewV4().String(),
		}

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("GET", "/v1/packages/"+ureq.Uuid, nil)

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &ppb.GetPackageRequest{
			Uuid: ureq.Uuid,
		}

		pResp := &ppb.GetPackageResponse{
			Package: &ppb.Package{
				Uuid: ureq.Uuid,
			},
		}

		p.On("Get", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("GetPackageDetails", func(t *testing.T) {
		ureq := PackagesRequest{
			Uuid: uuid.NewV4().String(),
		}

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("GET", "/v1/packages/"+ureq.Uuid+"/details", nil)

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &ppb.GetPackageRequest{
			Uuid: ureq.Uuid,
		}

		pResp := &ppb.GetPackageResponse{
			Package: &ppb.Package{
				Uuid: ureq.Uuid,
			},
		}

		p.On("GetDetails", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("GetPackages", func(t *testing.T) {
		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("GET", "/v1/packages", nil)

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &ppb.GetAllRequest{}

		pResp := &ppb.GetAllResponse{
			Packages: []*ppb.Package{
				{
					Name: "my-pack",
				},
			},
		}

		p.On("GetAll", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("DeletePackage", func(t *testing.T) {
		ureq := PackagesRequest{
			Uuid: uuid.NewV4().String(),
		}

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("DELETE", "/v1/packages/"+ureq.Uuid, nil)

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &ppb.DeletePackageRequest{
			Uuid: ureq.Uuid,
		}

		pResp := &ppb.DeletePackageResponse{}

		p.On("Delete", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("AddPackage", func(t *testing.T) {
		ureq := AddPackageRequest{
			Name:          testPackageName,
			From:          testFromDate,
			To:            testToDate,
			OwnerId:       uuid.NewV4().String(),
			SimType:       testSimType,
			SmsVolume:     testSmsVolume,
			DataVolume:    testDataVolume,
			DataUnit:      testDataUnit,
			VoiceUnit:     testVoiceUnit,
			Duration:      testDuration,
			Type:          testPackageType,
			Flatrate:      false,
			Amount:        testAmount,
			Markup:        testMarkup,
			Apn:           testApn,
			Active:        true,
			VoiceVolume:   testVoiceVolume,
			BaserateId:    uuid.NewV4().String(),
			Overdraft:     testOverdraft,
			TrafficPolicy: testTrafficPolicy,
			Networks:      testNetworks,
			Country:       testCountry,
			Currency:      testCurrency,
		}

		jreq, err := json.Marshal(&ureq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("POST", "/v1/packages", bytes.NewReader(jreq))
		hreq.Header.Set("Content-Type", "application/json")

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &ppb.AddPackageRequest{
			Name:          ureq.Name,
			From:          ureq.From,
			To:            ureq.To,
			OwnerId:       ureq.OwnerId,
			SimType:       ureq.SimType,
			SmsVolume:     ureq.SmsVolume,
			DataVolume:    ureq.DataVolume,
			DataUnit:      ureq.DataUnit,
			VoiceUnit:     ureq.VoiceUnit,
			Duration:      ureq.Duration,
			Type:          ureq.Type,
			Flatrate:      ureq.Flatrate,
			Amount:        ureq.Amount,
			Markup:        ureq.Markup,
			Apn:           ureq.Apn,
			Active:        ureq.Active,
			VoiceVolume:   ureq.VoiceVolume,
			BaserateId:    ureq.BaserateId,
			Overdraft:     ureq.Overdraft,
			TrafficPolicy: ureq.TrafficPolicy,
			Networks:      ureq.Networks,
			Country:       ureq.Country,
			Currency:      ureq.Currency,
		}

		pResp := &ppb.AddPackageResponse{
			Package: &ppb.Package{
				Uuid:    uuid.NewV4().String(),
				Name:    ureq.Name,
				OwnerId: ureq.OwnerId,
			},
		}

		p.On("Add", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), ureq.Name)
		assert.Contains(t, w.Body.String(), ureq.OwnerId)
		m.AssertExpectations(t)
	})

	t.Run("UpdatePackage", func(t *testing.T) {
		packageUuid := uuid.NewV4().String()
		ureq := UpdatePackageRequest{
			Uuid:   packageUuid,
			Name:   testUpdatedName,
			Active: false,
		}

		jreq, err := json.Marshal(&ureq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("PATCH", "/v1/packages/"+packageUuid, bytes.NewReader(jreq))
		hreq.Header.Set("Content-Type", "application/json")

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &ppb.UpdatePackageRequest{
			Uuid:   ureq.Uuid,
			Name:   ureq.Name,
			Active: ureq.Active,
		}

		pResp := &ppb.UpdatePackageResponse{
			Package: &ppb.Package{
				Uuid:   ureq.Uuid,
				Name:   ureq.Name,
				Active: ureq.Active,
			},
		}

		p.On("Update", mock.Anything, pReq).Return(pResp, nil)

		r := NewRouter(&Clients{
			r: client.NewRateClientFromClient(m),
			b: client.NewBaseRateClientFromClient(b),
			p: client.NewPackageFromClient(p),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()
		// act
		r.ServeHTTP(w, hreq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), ureq.Name)
		assert.Contains(t, w.Body.String(), ureq.Uuid)
		m.AssertExpectations(t)
	})
}

func TestRouter_GetBaseRatesForPackage(t *testing.T) {
	ureq := GetBaseRatesForPeriodRequest{
		Country:  testCountry,
		Provider: testProvider,
		SimType:  testSimType,
		To:       testToTime,
		From:     testFromTime,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/baserates/package", nil)
	q := hreq.URL.Query()
	q.Add("country", ureq.Country)
	q.Add("provider", ureq.Provider)
	q.Add("sim_type", ureq.SimType)
	q.Add("to", ureq.To)
	q.Add("from", ureq.From)
	hreq.URL.RawQuery = q.Encode()

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &bpb.GetBaseRatesByPeriodRequest{
		Country:  ureq.Country,
		Provider: ureq.Provider,
		SimType:  ureq.SimType,
		From:     ureq.From,
		To:       ureq.To,
	}

	pResp := &bpb.GetBaseRatesResponse{
		Rates: []*bpb.Rate{
			{
				Uuid:     uuid.NewV4().String(),
				Country:  ureq.Country,
				Provider: ureq.Provider,
				SimType:  ureq.SimType,
			},
		},
	}

	b.On("GetBaseRatesForPackage", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), ureq.Country)
	assert.Contains(t, w.Body.String(), ureq.Provider)
	m.AssertExpectations(t)
}
