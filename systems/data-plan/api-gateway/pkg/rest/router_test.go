package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/client"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	bmocks "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen/mocks"
	pmocks "github.com/ukama/ukama/systems/data-plan/package/pb/gen/mocks"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	rmocks "github.com/ukama/ukama/systems/data-plan/rate/pb/gen/mocks"
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
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_GetRates(t *testing.T) {
	ownerId := uuid.NewV4().String()
	req := GetRateRequest{
		OwnerId:     ownerId,
		Country:     "USA",
		Provider:    "ABC",
		To:          1680733308,
		From:        1680703308,
		SimType:     "ukama_data",
		EffectiveAt: "xx",
	}

	jReq, err := json.Marshal(req)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("POST", "/v1/rates/users/"+ownerId+"/rate", bytes.NewReader(jReq))

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}

	pReq := &rpb.GetRateRequest{
		OwnerId:     req.OwnerId,
		Country:     req.Country,
		Provider:    req.Provider,
		To:          req.To,
		From:        req.From,
		SimType:     req.SimType,
		EffectiveAt: req.EffectiveAt,
	}

	pResp := &rpb.GetRateResponse{
		Rates: []*bpb.Rate{
			{
				X2G:         true,
				X3G:         true,
				Apn:         "Manual entry required",
				Country:     req.Country,
				Data:        0.0014,
				EffectiveAt: "2023-10-10",
				Imsi:        1,
				Lte:         true,
				Network:     "Multi Tel",
				SimType:     req.SimType,
				SmsMo:       0.0100,
				SmsMt:       0.0001,
				Vpmn:        "TTC",
			},
		},
	}

	m.On("GetRate", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig).f.Engine()

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

	pReq := &rpb.GetMarkupRequest{
		OwnerId: req.OwnerId,
	}

	pResp := &rpb.GetMarkupResponse{
		OwnerId: req.OwnerId,
		Markup:  10,
	}

	m.On("GetMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig).f.Engine()

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

	pReq := &rpb.DeleteMarkupRequest{
		OwnerId: req.OwnerId,
	}

	pResp := &rpb.DeleteMarkupResponse{}

	m.On("DeleteMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig).f.Engine()

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
		Markup:  10,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("POST", "/v1/markup/"+strconv.FormatFloat(req.Markup, 'f', 'g', 64)+"/users/"+ownerId, nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}

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
	}, routerConfig).f.Engine()

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

	pReq := &rpb.GetMarkupHistoryRequest{
		OwnerId: req.OwnerId,
	}

	pResp := &rpb.GetMarkupHistoryResponse{
		OwnerId: req.OwnerId,
		MarkupRates: []*rpb.MarkupRates{
			{
				CreatedAt: "2021-11-12T11:45:26.371Z",
				DeletedAt: "2022-11-12T11:45:26.371Z",
				Markup:    5.5,
			},
			{
				CreatedAt: "2022-11-12T11:45:26.371Z",
				Markup:    10,
			},
		},
	}

	m.On("GetMarkupHistory", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig).f.Engine()

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
		Markup: 10,
	}

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("POST", "/v1/markup/"+strconv.FormatFloat(req.Markup, 'f', 'g', 64)+"/default", nil)

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}

	pReq := &rpb.UpdateDefaultMarkupRequest{
		Markup: req.Markup,
	}

	pResp := &rpb.UpdateDefaultMarkupResponse{}

	m.On("UpdateDefaultMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig).f.Engine()

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

	pReq := &rpb.GetDefaultMarkupRequest{}

	pResp := &rpb.GetDefaultMarkupResponse{
		Markup: 10,
	}

	m.On("GetDefaultMarkup", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig).f.Engine()

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

	pReq := &rpb.GetDefaultMarkupHistoryRequest{}

	pResp := &rpb.GetDefaultMarkupHistoryResponse{
		MarkupRates: []*rpb.MarkupRates{
			{
				CreatedAt: "2021-11-12T11:45:26.371Z",
				DeletedAt: "2022-11-12T11:45:26.371Z",
				Markup:    5.5,
			},
			{
				CreatedAt: "2022-11-12T11:45:26.371Z",
				Markup:    10,
			},
		},
	}

	m.On("GetDefaultMarkupHistory", mock.Anything, pReq).Return(pResp, nil)

	r := NewRouter(&Clients{
		r: client.NewRateClientFromClient(m),
		b: client.NewBaseRateClientFromClient(b),
		p: client.NewPackageFromClient(p),
	}, routerConfig).f.Engine()
	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), pResp.MarkupRates[0].CreatedAt)
	assert.Contains(t, w.Body.String(), pResp.MarkupRates[1].CreatedAt)
	m.AssertExpectations(t)
}
