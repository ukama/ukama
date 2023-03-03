package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/client"
	submocks "github.com/ukama/ukama/systems/subscriber/registry/pb/gen/mocks"
	smmocks "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen/mocks"
	spPb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	spmocks "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var Iccid = "1234567890123456789"
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
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout:    1 * time.Second,
		SimPool:    "0.0.0.0:9091",
		SimManager: "0.0.0.0:9092",
		Registry:   "0.0.0.0:9093",
	})
}

func TestPingRoute(t *testing.T) {
	// arrange
	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_getSimByIccid(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/simpool/sim/"+Iccid, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	preq := &spPb.GetByIccidRequest{
		Iccid: Iccid,
	}
	csp.On("GetByIccid", mock.Anything, preq).Return(&spPb.GetByIccidResponse{Sim: &spPb.Sim{
		Id:             1,
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
		SimType:        "ukama_data",
		SmDpAddress:    "http://localhost:8080",
		IsAllocated:    false,
		ActivationCode: "123456",
	}}, nil)

	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	csp.AssertExpectations(t)

}

func TestRouter_getSimPoolStats(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/simpool/stats/"+"ukama_data", nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	preq := &spPb.GetStatsRequest{
		SimType: "ukama_data",
	}
	csp.On("GetStats", mock.Anything, preq).Return(&spPb.GetStatsResponse{
		Total:     10,
		Available: 5,
		Consumed:  5,
		Failed:    0,
	}, nil)

	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	csp.AssertExpectations(t)

}

func TestRouter_addSimsToSimPool(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/simpool",
		strings.NewReader(`{"sim_info": [{ "iccid": "1234567890123456789", "sim_type": "ukama_data", "msidn": "555-555-1234", "smdp_address": "http://example.com", "activation_code": "abc123", "qr_code": "qr123", "is_physical_sim": true}]}`))

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	preq := &spPb.AddRequest{
		Sim: []*spPb.AddSim{
			{
				Iccid: "1234567890123456789", SimType: "ukama_data", Msisdn: "555-555-1234", SmDpAddress: "http://example.com", ActivationCode: "abc123", QrCode: "qr123", IsPhysical: true,
			},
		},
	}
	csp.On("Add", mock.Anything, preq).Return(&spPb.AddResponse{
		Sim: []*spPb.Sim{
			{
				Id:             1,
				Iccid:          "1234567890123456789",
				Msisdn:         "555-555-1234",
				SimType:        "ukama_data",
				SmDpAddress:    "http://localhost:8080",
				IsAllocated:    false,
				ActivationCode: "abc123",
			},
		},
	}, nil)

	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	csp.AssertExpectations(t)

}

func TestRouter_deleteSimFromSimPool(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/simpool/sim/1",
		nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	preq := &spPb.DeleteRequest{
		Id: []uint64{1},
	}
	csp.On("Delete", mock.Anything, preq).Return(&spPb.DeleteResponse{
		Id: []uint64{1},
	}, nil)

	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	csp.AssertExpectations(t)

}
