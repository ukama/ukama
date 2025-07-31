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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg/client"

	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	amocks "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen/mocks"
)

var (
	testClientSet *Clients

	defaultCors = cors.Config{
		AllowAllOrigins: true,
	}

	routerConfig = &RouterConfig{
		serverConf: &rest.HttpConfig{
			Cors: defaultCors,
		},
		httpEndpoints: &pkg.HttpEndpoints{
			NodeMetrics: "localhost:8080",
		},
		auth: &config.Auth{
			BypassAuthMode: true,
		},
	}
)

var (
	iccid   = "012345678901234567891"
	network = "40987edb-ebb6-4f84-a27c-99db7c136127"

	// orgId = "880f7c63-eb57-461a-b514-248ce91e9b3e"
	packageId = "8adcdfb4-ed30-405d-b32f-d0b2dda4a1e0"

	sub = pb.ReadResp{
		Record: &pb.Record{
			Iccid:       iccid,
			Imsi:        "012345678912345",
			Op:          []byte("0123456789012345"),
			Key:         []byte("0123456789012345"),
			Amf:         []byte("800"),
			AlgoType:    1,
			UeDlAmbrBps: 2000000,
			UeUlAmbrBps: 2000000,
			Sqn:         1,
			CsgIdPrsent: false,
			CsgId:       0,
			PackageId:   packageId,
		},
	}
)

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout: 1 * time.Second,
		Asr:     "localhost:9090",
	})
}

func TestRouter_PingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_PutSubscriber(t *testing.T) {
	w := httptest.NewRecorder()

	httpreq := ActivateReq{
		Iccid:     iccid,
		NetworkId: network,
		PackageId: packageId,
	}

	jReq, err := json.Marshal(httpreq)
	assert.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/v1/asr/"+iccid, bytes.NewReader(jReq))

	m := &amocks.AsrRecordServiceClient{}

	pReq := &pb.ActivateReq{
		Iccid:     iccid,
		NetworkId: network,
		PackageId: packageId,
	}

	m.On("Activate", mock.Anything, pReq).Return(&pb.ActivateResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_DeleteSubscriber(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/asr/"+iccid, nil)

	m := &amocks.AsrRecordServiceClient{}

	pReq := &pb.InactivateReq{
		Iccid: iccid,
	}

	m.On("Inactivate", mock.Anything, pReq).Return(&pb.InactivateResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_PatchSubscriber(t *testing.T) {
	httpreq := UpdatePackageReq{
		Iccid:     iccid,
		PackageId: packageId,
	}

	jReq, err := json.Marshal(httpreq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/asr/"+iccid, bytes.NewReader(jReq))

	m := &amocks.AsrRecordServiceClient{}
	pReq := &pb.UpdatePackageReq{
		Iccid:     iccid,
		PackageId: packageId,
	}
	m.On("UpdatePackage", mock.Anything, pReq).Return(&pb.UpdatePackageResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_GetSubscriber(t *testing.T) {
	t.Run("SubscriberExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/asr/"+iccid, nil)

		m := &amocks.AsrRecordServiceClient{}

		pReq := &pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: iccid,
			},
		}
		m.On("Read", mock.Anything, pReq).Return(nil, fmt.Errorf("not found"))

		r := NewRouter(&Clients{
			a: client.NewAsrFromClient(m),
		}, routerConfig, nil).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("SubscriberDoesn'tExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/asr/"+iccid, nil)

		m := &amocks.AsrRecordServiceClient{}

		pReq := &pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: iccid,
			},
		}

		m.On("Read", mock.Anything, pReq).Return(&sub, nil)

		r := NewRouter(&Clients{
			a: client.NewAsrFromClient(m),
		}, routerConfig, nil).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
}

func TestRouter_GetUsage(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/asr/"+iccid+"/usage", nil)

	m := &amocks.AsrRecordServiceClient{}
	pReq := &pb.UsageReq{
		Id: &pb.UsageReq_Iccid{
			Iccid: iccid,
		},
	}
	m.On("GetUsage", mock.Anything, pReq).Return(&pb.UsageResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_GetUsageForPeriod(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/asr/"+iccid+"/period?start_time=1714008143&end_time=1714539344", nil)

	m := &amocks.AsrRecordServiceClient{}
	pReq := &pb.UsageForPeriodReq{
		Id: &pb.UsageForPeriodReq_Iccid{
			Iccid: iccid,
		},
		StartTime: 1714008143,
		EndTime:   1714539344,
	}
	m.On("GetUsageForPeriod", mock.Anything, pReq).Return(&pb.UsageResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}
