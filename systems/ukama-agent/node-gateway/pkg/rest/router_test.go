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
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	amocks "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen/mocks"
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/pkg/client"
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

var iccid = "012345678901234567891"
var imsi = "012345678912345"

// var orgId = "880f7c63-eb57-461a-b514-248ce91e9b3e"
var packageId = "8adcdfb4-ed30-405d-b32f-d0b2dda4a1e0"

var sub = pb.ReadResp{
	Record: &pb.Record{
		Iccid:       iccid,
		SimId:       "880f7c63-eb57-461a-b514-248ce91e9b3e",
		Imsi:        imsi,
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

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout: 1 * time.Second,
	})
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

func TestRouter_PostGuti(t *testing.T) {
	w := httptest.NewRecorder()
	req := UpdateGutiReq{
		Guti: Guti{
			PlmnId: "00101",
			Mmegi:  3200,
			Mmec:   100,
			Mtmsi:  1,
		},
		UpdatedAt: uint32(time.Now().Unix()),
	}

	body, _ := json.Marshal(req)

	hreq, _ := http.NewRequest("POST", "/v1/subscriber/"+imsi+"/guti",
		bytes.NewBuffer(body))

	m := &amocks.AsrRecordServiceClient{}
	pReq := &pb.UpdateGutiReq{
		Imsi:      imsi,
		UpdatedAt: req.UpdatedAt,
		Guti: &pb.Guti{
			PlmnId: req.Guti.PlmnId,
			Mmegi:  req.Guti.Mmegi,
			Mmec:   req.Guti.Mmec,
			Mtmsi:  req.Guti.Mtmsi,
		},
	}
	m.On("UpdateGuti", mock.Anything, pReq).Return(&pb.UpdateGutiResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_PostTai(t *testing.T) {
	w := httptest.NewRecorder()
	req := UpdateTaiReq{
		PlmnId:    "00101",
		Tac:       1,
		UpdatedAt: uint32(time.Now().Unix()),
	}

	body, _ := json.Marshal(req)

	hreq, _ := http.NewRequest("POST", "/v1/subscriber/"+imsi+"/tai",
		bytes.NewBuffer(body))

	m := &amocks.AsrRecordServiceClient{}
	pReq := &pb.UpdateTaiReq{
		Imsi:      imsi,
		UpdatedAt: req.UpdatedAt,
		Tac:       req.Tac,
	}
	m.On("UpdateTai", mock.Anything, pReq).Return(&pb.UpdateTaiResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_GetSubscriber(t *testing.T) {
	t.Run("SubscriberExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/subscriber/"+imsi, nil)

		m := &amocks.AsrRecordServiceClient{}

		pReq := &pb.ReadReq{
			Id: &pb.ReadReq_Imsi{
				Imsi: imsi,
			},
		}
		m.On("Read", mock.Anything, pReq).Return(nil, fmt.Errorf("not found"))

		r := NewRouter(&Clients{
			a: client.NewAsrFromClient(m),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("SubscriberDoesn'tExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/subscriber/"+imsi, nil)

		m := &amocks.AsrRecordServiceClient{}

		pReq := &pb.ReadReq{
			Id: &pb.ReadReq_Imsi{
				Imsi: imsi,
			},
		}

		m.On("Read", mock.Anything, pReq).Return(&sub, nil)

		r := NewRouter(&Clients{
			a: client.NewAsrFromClient(m),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
}
