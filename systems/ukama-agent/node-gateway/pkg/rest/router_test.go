package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
var network = "40987edb-ebb6-4f84-a27c-99db7c136127"

// var orgId = "880f7c63-eb57-461a-b514-248ce91e9b3e"
var packageId = "8adcdfb4-ed30-405d-b32f-d0b2dda4a1e0"

var sub = pb.ReadResp{
	Record: &pb.Record{
		Iccid:       iccid,
		SimId:       "880f7c63-eb57-461a-b514-248ce91e9b3e",
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

func TestRouter_PutSubscriber(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/subscriber/"+iccid,
		strings.NewReader(`{"network":"`+network+`","packageId":"`+packageId+`"}`))

	m := &amocks.AsrRecordServiceClient{}

	pReq := &pb.ActivateReq{
		Iccid:     iccid,
		Network:   network,
		PackageId: packageId,
	}

	m.On("Activate", mock.Anything, pReq).Return(&pb.ActivateResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_DeleteSubscriber(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/subscriber/"+iccid, nil)

	m := &amocks.AsrRecordServiceClient{}

	pReq := &pb.InactivateReq{
		Id: &pb.InactivateReq_Iccid{
			Iccid: iccid,
		},
	}

	m.On("Inactivate", mock.Anything, pReq).Return(&pb.InactivateResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_PatchSubscriber(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/subscriber/"+iccid,
		strings.NewReader(`{"packageId":"`+packageId+`"}`))

	m := &amocks.AsrRecordServiceClient{}
	pReq := &pb.UpdatePackageReq{
		Iccid:     iccid,
		PackageId: packageId,
	}
	m.On("UpdatePackage", mock.Anything, pReq).Return(&pb.UpdatePackageResp{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_GetSubscriber(t *testing.T) {
	t.Run("SubscriberExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/subscriber/"+iccid, nil)

		m := &amocks.AsrRecordServiceClient{}

		pReq := &pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: iccid,
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
		req, _ := http.NewRequest("GET", "/v1/subscriber/"+iccid, nil)

		m := &amocks.AsrRecordServiceClient{}

		pReq := &pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: iccid,
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
