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
		OwnerId:  ownerId,
		Country:  "USA",
		Provider: "ABC",
		To:       time.Now().UTC().Format(time.RFC3339),
		From:     time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		SimType:  "ukama_data",
	}

	jReq, err := json.Marshal(req)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("POST", "/v1/rates/users/"+ownerId+"/rate", bytes.NewReader(jReq))

	m := &rmocks.RateServiceClient{}
	p := &pmocks.PackagesServiceClient{}
	b := &bmocks.BaseRatesServiceClient{}
	arc := &providers.AuthRestClient{}
	pReq := &rpb.GetRateRequest{
		OwnerId:  req.OwnerId,
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
				Apn:         "Manual entry required",
				Country:     req.Country,
				Data:        0.0014,
				EffectiveAt: "2023-10-10",
				Imsi:        1,
				Lte:         true,
				Provider:    "Multi Tel",
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
		Markup:  10,
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
		Markup:  10,
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
		Markup: 10,
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
		Markup: 10,
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
		FileURL:     "https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/base-rate/template/template.csv",
		EffectiveAt: "2023-10-12T07:20:50.52Z",
		EndAt:       "2043-10-12T07:20:50.52Z",
		SimType:     "ukama_data",
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
			Country:  "ABC",
			Provider: "XYZ",
			SimType:  "ukama_data",
		}

		jreq, err := json.Marshal(&ureq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("POST", "/v1/baserates/country/"+ureq.Country, bytes.NewReader(jreq))

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
			Country:  "ABC",
			Provider: "XYZ",
			SimType:  "ukama_data",
		}

		jreq, err := json.Marshal(&ureq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("POST", "/v1/baserates/country/"+ureq.Country+"/history", bytes.NewReader(jreq))

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
			Country:  "ABC",
			Provider: "XYZ",
			SimType:  "ukama_data",
			To:       "2023-10-12T07:20:50.52Z",
			From:     "2022-10-12T07:20:50.52Z",
		}

		jreq, err := json.Marshal(&ureq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("POST", "/v1/baserates/country/"+ureq.Country+"/period", bytes.NewReader(jreq))

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

	t.Run("GetPackageByOrgId", func(t *testing.T) {
		ureq := GetPackageByOrgRequest{
			OrgId: uuid.NewV4().String(),
		}

		w := httptest.NewRecorder()
		hreq, _ := http.NewRequest("GET", "/v1/packages/org/"+ureq.OrgId, nil)

		m := &rmocks.RateServiceClient{}
		p := &pmocks.PackagesServiceClient{}
		b := &bmocks.BaseRatesServiceClient{}
		arc := &providers.AuthRestClient{}
		pReq := &ppb.GetByOrgPackageRequest{
			OrgId: ureq.OrgId,
		}

		pResp := &ppb.GetByOrgPackageResponse{
			Packages: []*ppb.Package{
				{
					OrgId: ureq.OrgId,
				},
			},
		}

		p.On("GetByOrg", mock.Anything, pReq).Return(pResp, nil)

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

}
