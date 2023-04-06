package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		d: client.NewPackageFromClient(p, b),
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
	hreq, _ := http.NewRequest("POST", "/v1/rates/users/"+ownerId, bytes.NewReader(jReq))

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
		d: client.NewPackageFromClient(p, b),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}
