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
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"

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

// Test data constants
const (
	testIccid           = "1234567890123456789"
	testMsisdn          = "555-555-1234"
	testMsisdn2         = "2345678901"
	testSimType         = "ukama_data"
	testSmDpAddress     = "http://example.com"
	testSmDpAddress2    = "http://localhost:8080"
	testActivationCode  = "abc123"
	testActivationCode2 = "123456"
	testQrCode          = "qr123"
	testSimToken        = "abcdef"
	testEmail           = "johndoe@example.com"
	testPhone           = "1234567890"
	testName            = "John"
	testGender          = "Male"
	testDob             = "16-04-1995"
	testAddress         = "1 Main St"
	testProofOfId       = "Passport"
	testIdSerial        = "123456789"
	testStatus          = "active"
	testCdrType         = "voice"
	testRegion          = "US"
	testFromDate        = "2023-01-01"
	testToDate          = "2023-01-31"
	testDataPlanId      = "plan123"
	testUsage           = 100.5
	testUsageUnit       = "minutes"
)

// Test UUIDs
var (
	testSubscriberId = "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a9f"
	testNetworkId    = "9e82c8b1-a746-4f2c-a80e-f4d14d863ea3"
	testSimId        = "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a11"
	testImsi         = "01234567891234"
)

// Test package data
var (
	testPackageId = uuid.NewV4().String()
	testStartDate = time.Now().UTC().Format(time.RFC3339)
	testEndDate   = time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
)

// Test CSV data
const (
	testCsvData  = "SIM,ICCID,MSISDN,SmDpAddress,ActivationCode,IsPhysical,QRCode\n8910300000003540855,880170124847571,1001.9.0.0.1,1010,TRUE,459081a\n"
	testCsvB64   = "U0lNLElDQ0lELE1TSVNETixTbURwQWRkcmVzcyxBY3RpdmF0aW9uQ29kZSxJc1BoeXNpY2FsLFFSQ29kZQo4OTEwMzAwMDAwMDAzNTQwODU1LDg4MDE3MDEyNDg0NzU3MSwxMDAxLjkuMC4wLjEsMTAxMCxUUlVFLDQ1OTA4MWEK"
	testCsvIccid = "8910300000003540855"
)

// Test stats data
const (
	testTotalSims     = 10
	testAvailableSims = 5
	testConsumedSims  = 5
	testFailedSims    = 0
)

// Test counts and flags
const (
	testCount         = 10
	testSort          = true
	testIsActive      = true
	testIsPhysical    = false
	testIsPhysicalSim = true
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
	req, _ := http.NewRequest("GET", "/v1/simpool/sim/"+testIccid, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}
	preq := &spPb.GetByIccidRequest{
		Iccid: testIccid,
	}
	csp.On("GetByIccid", mock.Anything, preq).Return(&spPb.GetByIccidResponse{Sim: &spPb.Sim{
		Id:             1,
		Iccid:          testIccid,
		Msisdn:         testMsisdn,
		SimType:        testSimType,
		SmDpAddress:    testSmDpAddress,
		IsAllocated:    false,
		ActivationCode: testActivationCode,
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
	req, _ := http.NewRequest("GET", "/v1/simpool/stats/"+testSimType, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}
	preq := &spPb.GetStatsRequest{
		SimType: testSimType,
	}
	csp.On("GetStats", mock.Anything, preq).Return(&spPb.GetStatsResponse{
		Total:     testTotalSims,
		Available: testAvailableSims,
		Consumed:  testConsumedSims,
		Failed:    testFailedSims,
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
		strings.NewReader(`{"sim_info": [{ "iccid": "`+testIccid+`", "sim_type": "`+testSimType+`", "msisdn": "`+testMsisdn+`", "smdp_address": "`+testSmDpAddress+`", "activation_code": "`+testActivationCode+`", "qr_code": "`+testQrCode+`", "is_physical_sim": true}]}`))

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}
	preq := &spPb.AddRequest{
		Sim: []*spPb.AddSim{
			{
				Iccid: testIccid, SimType: testSimType, Msisdn: testMsisdn, SmDpAddress: testSmDpAddress, ActivationCode: testActivationCode, QrCode: testQrCode, IsPhysical: true,
			},
		},
	}
	csp.On("Add", mock.Anything, preq).Return(&spPb.AddResponse{
		Sim: []*spPb.Sim{
			{
				Id:             1,
				Iccid:          testIccid,
				Msisdn:         testMsisdn,
				SimType:        testSimType,
				SmDpAddress:    testSmDpAddress,
				IsAllocated:    false,
				ActivationCode: testActivationCode,
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
		SubscriberId:          testSubscriberId,
		Name:                  testName,
		NetworkId:             testNetworkId,
		Email:                 testEmail,
		PhoneNumber:           testPhone,
		Gender:                testGender,
		Dob:                   testDob,
		Address:               testAddress,
		ProofOfIdentification: testProofOfId,
		IdSerial:              testIdSerial,
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
			Name:                  testName,
			NetworkId:             testNetworkId,
			Email:                 testEmail,
			Phone:                 testPhone,
			Gender:                testGender,
			Dob:                   testDob,
			Address:               testAddress,
			ProofOfIdentification: testProofOfId,
			IdSerial:              testIdSerial,
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
			Name:                  testName,
			Phone:                 testPhone,
			Address:               testAddress,
			ProofOfIdentification: testProofOfId,
			IdSerial:              testIdSerial,
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
	subscriberId := testSubscriberId
	sim := &smPb.Sim{
		Id:           testSimId,
		SubscriberId: testSubscriberId,
		NetworkId:    testNetworkId,
		Iccid:        testIccid,
		Msisdn:       testMsisdn,
		Type:         testSimType,
		Imsi:         testImsi,
		IsPhysical:   testIsPhysical,
		Package: &smPb.Package{
			Id:        testPackageId,
			StartDate: testStartDate,
			EndDate:   testEndDate,
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
			SimToken:     testSimToken,
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
			Status: testStatus,
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

func TestRouter_getSims(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/simpool/sims/"+testSimType, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}
	preq := &spPb.GetSimsRequest{
		SimType: testSimType,
	}
	csp.On("GetSims", mock.Anything, preq).Return(&spPb.GetSimsResponse{
		Sims: []*spPb.Sim{
			{
				Id:             1,
				Iccid:          testIccid,
				Msisdn:         testMsisdn,
				SimType:        testSimType,
				SmDpAddress:    testSmDpAddress,
				IsAllocated:    false,
				ActivationCode: testActivationCode,
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
	assert.Equal(t, http.StatusOK, w.Code)
	csp.AssertExpectations(t)
}

func TestRouter_uploadSimsToSimPool(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/simpool/upload",
		strings.NewReader(`{"sim_type": "`+testSimType+`", "data": "`+testCsvB64+`"}`))

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	preq := &spPb.UploadRequest{
		SimType: testSimType,
		SimData: []byte(testCsvData),
	}
	csp.On("Upload", mock.Anything, preq).Return(&spPb.UploadResponse{
		Iccid: []string{testCsvIccid},
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

func TestRouter_getSubscriberByEmail(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/subscriber/email/"+testEmail, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	s := &upb.Subscriber{
		SubscriberId:          testSubscriberId,
		Name:                  testName,
		NetworkId:             testNetworkId,
		Email:                 testEmail,
		PhoneNumber:           testPhone,
		Gender:                testGender,
		Dob:                   testDob,
		Address:               testAddress,
		ProofOfIdentification: testProofOfId,
		IdSerial:              testIdSerial,
	}

	preq := &subPb.GetSubscriberByEmailRequest{
		Email: testEmail,
	}
	csub.On("GetByEmail", mock.Anything, preq).Return(&subPb.GetSubscriberByEmailResponse{
		Subscriber: s,
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
	assert.Contains(t, w.Body.String(), `"email":"`+testEmail+`"`)
	csub.AssertExpectations(t)
}

func TestRouter_getSubscriberByNetwork(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/subscribers/networks/"+testNetworkId, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	subscribers := []*upb.Subscriber{
		{
			SubscriberId:          testSubscriberId,
			Name:                  testName,
			NetworkId:             testNetworkId,
			Email:                 testEmail,
			PhoneNumber:           testPhone,
			Gender:                testGender,
			Dob:                   testDob,
			Address:               testAddress,
			ProofOfIdentification: testProofOfId,
			IdSerial:              testIdSerial,
		},
	}

	preq := &subPb.GetByNetworkRequest{
		NetworkId: testNetworkId,
	}
	csub.On("GetByNetwork", mock.Anything, preq).Return(&subPb.GetByNetworkResponse{
		Subscribers: subscribers,
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
	assert.Contains(t, w.Body.String(), `"network_id":"`+testNetworkId+`"`)
	csub.AssertExpectations(t)
}

func TestRouter_listSims(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/sim?subscriber_id="+testSubscriberId+"&network_id="+testNetworkId+"&sim_type="+testSimType+"&count="+strconv.Itoa(testCount)+"&sort="+strconv.FormatBool(testSort), nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	sim := &smPb.Sim{
		Id:           testSimId,
		SubscriberId: testSubscriberId,
		NetworkId:    testNetworkId,
		Iccid:        testIccid,
		Msisdn:       testMsisdn,
		Type:         testSimType,
		Imsi:         testImsi,
		IsPhysical:   testIsPhysical,
	}

	preq := &smPb.ListSimsRequest{
		SubscriberId: testSubscriberId,
		NetworkId:    testNetworkId,
		SimType:      testSimType,
		Count:        testCount,
		Sort:         testSort,
	}
	csm.On("ListSims", mock.Anything, preq).Return(&smPb.ListSimsResponse{
		Sims: []*smPb.Sim{sim},
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
	assert.Contains(t, w.Body.String(), `"subscriber_id":"`+testSubscriberId+`"`)
	csm.AssertExpectations(t)
}

func TestRouter_getSimsByNetwork(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/sims/networks/"+testNetworkId, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	sim := &smPb.Sim{
		Id:           testSimId,
		SubscriberId: testSubscriberId,
		NetworkId:    testNetworkId,
		Iccid:        testIccid,
		Msisdn:       testMsisdn,
		Type:         testSimType,
		Imsi:         testImsi,
		IsPhysical:   testIsPhysical,
	}

	preq := &smPb.GetSimsByNetworkRequest{
		NetworkId: testNetworkId,
	}
	csm.On("GetSimsByNetwork", mock.Anything, preq).Return(&smPb.GetSimsByNetworkResponse{
		Sims: []*smPb.Sim{sim},
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
	assert.Contains(t, w.Body.String(), `"network_id":"`+testNetworkId+`"`)
	csm.AssertExpectations(t)
}

func TestRouter_listPackagesForSim(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/sim/"+testSimId+"/package?data_plan_id="+testDataPlanId+"&is_active="+strconv.FormatBool(testIsActive)+"&count="+strconv.Itoa(testCount)+"&sort="+strconv.FormatBool(testSort), nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	package1 := &smPb.Package{
		Id:        testPackageId,
		StartDate: testStartDate,
		EndDate:   testEndDate,
	}

	// Mock the gRPC client call with the correct request structure
	preq := &smPb.ListPackagesForSimRequest{
		SimId:      testSimId,
		DataPlanId: testDataPlanId,
		IsActive:   testIsActive,
		Count:      testCount,
		Sort:       testSort,
	}
	csm.On("ListPackagesForSim", mock.Anything, preq).Return(&smPb.ListPackagesForSimResponse{
		Packages: []*smPb.Package{package1},
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
	assert.Contains(t, w.Body.String(), `"id":"`+testPackageId+`"`)
	csm.AssertExpectations(t)
}

func TestRouter_addPackageForSim(t *testing.T) {
	simId := testSimId
	packageId := testPackageId
	startDate := testStartDate

	data := AddPkgToSimReq{
		PackageId: packageId,
		StartDate: startDate,
	}

	jdata, err := json.Marshal(&data)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/v1/sim/"+simId+"/package", bytes.NewReader(jdata))
	assert.NoError(t, err)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	preq := &smPb.AddPackageRequest{
		SimId:     simId,
		PackageId: packageId,
		StartDate: startDate,
	}
	csm.On("AddPackageForSim", mock.Anything, preq).Return(&smPb.AddPackageResponse{}, nil)

	r := NewRouter(&Clients{
		sp:  client.NewSimPoolFromClient(csp),
		sm:  client.NewSimManagerFromClient(csm),
		sub: client.NewRegistryFromClient(csub),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	csm.AssertExpectations(t)
}

func TestRouter_getUsages(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/usages?sim_id="+testSimId+"&sim_type="+testSimType+"&cdr_type="+testCdrType+"&from="+testFromDate+"&to="+testToDate+"&region="+testRegion, nil)

	csp := &spmocks.SimServiceClient{}
	csm := &smmocks.SimManagerServiceClient{}
	csub := &submocks.RegistryServiceClient{}
	arc := &providers.AuthRestClient{}

	// Create a mock usage response with structpb.Struct
	usageData := map[string]interface{}{
		"sim_id":   testSimId,
		"sim_type": testSimType,
		"cdr_type": testCdrType,
		"from":     testFromDate,
		"to":       testToDate,
		"region":   testRegion,
		"usage":    testUsage,
		"unit":     testUsageUnit,
	}

	usageStruct, _ := structpb.NewStruct(usageData)

	// Mock the gRPC client call with the correct request structure
	preq := &smPb.UsageRequest{
		SimId:   testSimId,
		SimType: testSimType,
		Type:    testCdrType,
		From:    testFromDate,
		To:      testToDate,
		Region:  testRegion,
	}
	csm.On("GetUsages", mock.Anything, preq).Return(&smPb.UsageResponse{
		Usage: usageStruct,
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
	assert.Contains(t, w.Body.String(), `"sim_id":"`+testSimId+`"`)
	csm.AssertExpectations(t)
}

func TestRouter_addReqToAddSimReqPb(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := &SimPoolAddSimReq{
			SimInfo: []SimInfo{
				{
					Iccid:          testIccid,
					SimType:        testSimType,
					Msisdn:         testMsisdn,
					SmDpAddress:    testSmDpAddress,
					ActivationCode: testActivationCode,
					QrCode:         testQrCode,
					IsPhysicalSim:  true,
				},
			},
		}

		result, err := addReqToAddSimReqPb(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Sim, 1)
		assert.Equal(t, testIccid, result.Sim[0].Iccid)
		assert.Equal(t, testSimType, result.Sim[0].SimType)
		assert.Equal(t, testMsisdn, result.Sim[0].Msisdn)
		assert.Equal(t, testSmDpAddress, result.Sim[0].SmDpAddress)
		assert.Equal(t, testActivationCode, result.Sim[0].ActivationCode)
		assert.Equal(t, testQrCode, result.Sim[0].QrCode)
		assert.True(t, result.Sim[0].IsPhysical)
	})

	t.Run("nil request", func(t *testing.T) {
		result, err := addReqToAddSimReqPb(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid add request")
	})
}
