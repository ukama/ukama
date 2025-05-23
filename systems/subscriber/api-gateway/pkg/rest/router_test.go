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
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	subPb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	submocks "github.com/ukama/ukama/systems/subscriber/registry/pb/gen/mocks"
	smPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	smmocks "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen/mocks"
	spPb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	spmocks "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen/mocks"
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
	arc := &providers.AuthRestClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

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
	arc := &providers.AuthRestClient{}
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
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

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
	arc := &providers.AuthRestClient{}
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
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	csp.AssertExpectations(t)

}

func TestRouter_addSimsToSimPool(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/simpool",
		strings.NewReader(`{"sim_info": [{ "iccid": "1234567890123456789", "sim_type": "ukama_data", "msisdn": "555-555-1234", "smdp_address": "http://example.com", "activation_code": "abc123", "qr_code": "qr123", "is_physical_sim": true}]}`))

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}
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
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

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
	arc := &providers.AuthRestClient{}
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
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	csp.AssertExpectations(t)

}

func TestRouter_Subscriber(t *testing.T) {
	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}
	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	s := &upb.Subscriber{
		SubscriberId:          "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a9f",
		Name:                  "John",
		NetworkId:             "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3",
		Email:                 "johndoe@example.com",
		PhoneNumber:           "1234567890",
		Gender:                "Male",
		Dob:                   "16-04-1995",
		Address:               "1 Main St",
		ProofOfIdentification: "Passport",
		IdSerial:              "123456789",
	}

	t.Run("getSubscriber", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/subscriber/"+s.SubscriberId,
			nil)

		preq := &subPb.GetSubscriberRequest{
			SubscriberId: s.SubscriberId,
		}
		csub.On("Get", mock.Anything, preq).Return(&subPb.GetSubscriberResponse{
			Subscriber: s,
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"subscriber_id":"`+s.SubscriberId+`"`)

		csp.AssertExpectations(t)
	})

	t.Run("putSubscriber", func(t *testing.T) {
		data := SubscriberAddReq{
			Name:                  "John",
			NetworkId:             "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3",
			Email:                 "johndoe@example.com",
			Phone:                 "1234567890",
			Gender:                "Male",
			Dob:                   "16-04-1995",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("PUT", "/v1/subscriber", bytes.NewReader(jdata))
		assert.NoError(t, err)

		preq := &subPb.AddSubscriberRequest{
			Name:                  data.Name,
			Email:                 data.Email,
			PhoneNumber:           data.Phone,
			Dob:                   data.Dob,
			Address:               data.Address,
			ProofOfIdentification: data.ProofOfIdentification,
			IdSerial:              data.IdSerial,
			NetworkId:             data.NetworkId,
			Gender:                data.Gender,
		}

		csub.On("Add", mock.Anything, preq).Return(&subPb.AddSubscriberResponse{
			Subscriber: s,
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"subscriber_id":"`+s.SubscriberId+`"`)
		csp.AssertExpectations(t)
	})

	t.Run("deleteSubscriber", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/subscriber/"+s.SubscriberId,
			nil)

		preq := &subPb.DeleteSubscriberRequest{
			SubscriberId: s.SubscriberId,
		}
		csub.On("Delete", mock.Anything, preq).Return(&subPb.DeleteSubscriberResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)

		csp.AssertExpectations(t)
	})

	t.Run("updateSubscriber", func(t *testing.T) {
		data := SubscriberUpdateReq{
			Name:                  "John",
			Phone:                 "1234567890",
			Address:               "1 Main St",
			ProofOfIdentification: "Passport",
			IdSerial:              "123456789",
		}

		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("PATCH", "/v1/subscriber/"+s.SubscriberId, bytes.NewReader(jdata))
		assert.NoError(t, err)

		preq := &subPb.UpdateSubscriberRequest{
			SubscriberId:          s.SubscriberId,
			Name:                  data.Name,
			PhoneNumber:           data.Phone,
			Address:               data.Address,
			ProofOfIdentification: data.ProofOfIdentification,
			IdSerial:              data.IdSerial,
		}
		csub.On("Update", mock.Anything, preq).Return(&subPb.UpdateSubscriberResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)

		csp.AssertExpectations(t)
	})
}

func TestRouter_SimManager(t *testing.T) {
	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}
	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()
	subscriberId := "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a9f"
	sim := &smPb.Sim{
		Id:           "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a11",
		SubscriberId: "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a9f",
		NetworkId:    "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3",
		Iccid:        "1234567890123456789",
		Msisdn:       "555-555-1234",
		Type:         "ukama_data",
		Imsi:         "01234567891234",
		IsPhysical:   false,
		Package: &smPb.Package{
			Id:        uuid.NewV4().String(),
			StartDate: time.Now().UTC().Format(time.RFC3339),
			EndDate:   time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		},
	}

	t.Run("getSims", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/sim/"+sim.Id,
			nil)

		preq := &smPb.GetSimRequest{
			SimId: sim.Id,
		}
		csm.On("GetSim", mock.Anything, preq).Return(&smPb.GetSimResponse{
			Sim: sim,
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":"`+sim.Id+`"`)

		csm.AssertExpectations(t)
	})

	t.Run("getSimsBySub", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/sim/subscriber/"+subscriberId,
			nil)

		preq := &smPb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberId,
		}
		csm.On("GetSimsBySubscriber", mock.Anything, preq).Return(&smPb.GetSimsBySubscriberResponse{
			Sims: []*smPb.Sim{sim},
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"subscriber_id":"`+subscriberId+`"`)

		csm.AssertExpectations(t)
	})

	t.Run("getPackagesForSim", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/sim/packages/"+sim.Id,
			nil)

		preq := &smPb.GetPackagesForSimRequest{
			SimId: sim.Id,
		}
		csm.On("GetPackagesForSim", mock.Anything, preq).Return(&smPb.GetPackagesForSimResponse{
			SimId:    sim.Id,
			Packages: []*smPb.Package{sim.Package},
		}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"sim_id":"`+sim.Id+`"`)
		assert.Contains(t, w.Body.String(), `"id":"`+sim.Package.Id+`"`)
		csm.AssertExpectations(t)
	})

	t.Run("addPkgForSim", func(t *testing.T) {
		p := PostPkgToSimReq{
			SimId:     sim.Id,
			PackageId: sim.Package.Id,
			StartDate: sim.Package.StartDate,
		}

		jdata, err := json.Marshal(&p)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/sim/package",
			bytes.NewReader(jdata))
		assert.NoError(t, err)

		preq := &smPb.AddPackageRequest{
			SimId:     p.SimId,
			PackageId: p.PackageId,
			StartDate: p.StartDate,
		}
		csm.On("AddPackageForSim", mock.Anything, preq).Return(&smPb.AddPackageResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		csm.AssertExpectations(t)

	})

	t.Run("allocateSim", func(t *testing.T) {
		p := AllocateSimReq{
			SubscriberId: sim.SubscriberId,
			SimToken:     "abcdef",
			PackageId:    sim.Package.Id,
			NetworkId:    sim.NetworkId,
			SimType:      sim.Type,
		}

		jdata, err := json.Marshal(&p)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/sim/",
			bytes.NewReader(jdata))
		assert.NoError(t, err)

		preq := &smPb.AllocateSimRequest{
			SubscriberId: p.SubscriberId,
			SimToken:     p.SimToken,
			PackageId:    p.PackageId,
			NetworkId:    p.NetworkId,
			SimType:      p.SimType,
		}

		csm.On("AllocateSim", mock.Anything, preq).Return(&smPb.AllocateSimResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		csm.AssertExpectations(t)
	})
	t.Run("updateSimStatus", func(t *testing.T) {
		p := ActivateDeactivateSimReq{
			Status: "active",
		}

		jdata, err := json.Marshal(&p)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("PATCH", "/v1/sim/"+sim.Id,
			bytes.NewReader(jdata))
		assert.NoError(t, err)

		preq := &smPb.ToggleSimStatusRequest{
			SimId:  sim.Id,
			Status: p.Status,
		}

		csm.On("ToggleSimStatus", mock.Anything, preq).Return(&smPb.ToggleSimStatusResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		csm.AssertExpectations(t)
	})
	t.Run("setActivePackageForSim", func(t *testing.T) {

		w := httptest.NewRecorder()
		req, err := http.NewRequest("PATCH", "/v1/sim/"+sim.Id+"/package/"+sim.Package.Id,
			nil)
		assert.NoError(t, err)

		preq := &smPb.SetActivePackageRequest{
			SimId:     sim.Id,
			PackageId: sim.Package.Id,
		}

		csm.On("SetActivePackageForSim", mock.Anything, preq).Return(&smPb.SetActivePackageResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		csm.AssertExpectations(t)
	})
	t.Run("removePkgForSim", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("DELETE", "/v1/sim/"+sim.Id+"/package/"+sim.Package.Id,
			nil)
		assert.NoError(t, err)

		preq := &smPb.RemovePackageRequest{
			SimId:     sim.Id,
			PackageId: sim.Package.Id,
		}

		csm.On("RemovePackageForSim", mock.Anything, preq).Return(&smPb.RemovePackageResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		csm.AssertExpectations(t)
	})
	t.Run("deleteSim", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("DELETE", "/v1/sim/"+sim.Id,
			nil)
		assert.NoError(t, err)

		preq := &smPb.TerminateSimRequest{
			SimId: sim.Id,
		}

		csm.On("TerminateSim", mock.Anything, preq).Return(&smPb.TerminateSimResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		csm.AssertExpectations(t)
	})
}
