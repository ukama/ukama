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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	amocks "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

var iccid = "1111111111-test-iccid"

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
	req, _ := http.NewRequest("PUT", "/v1/susbcriber"+iccid, 
        strings.NewReader(`{"Certificate": "helloOrg","Ip": "0.0.0.0"}`))

	m := &amocks.AsrRecordServiceClient{}

	req := pb.ActivateReq{
		Iccid:     req.Iccid,
		Network:   req.Network,
		PackageId: req.PackageId,
	}

	m.On("Activate", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_DeleteSubscriber(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name", 
	)))

	m := &amocks.AsrRecordServiceClient{}

	m.On("Inactivate", mock.Anything, mock.Anything).Return(&pb.GetOrgResponse{
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
	}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"orgName":"org-name"`)
}

func TestRouter_PatchSubscriber(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/orgs/org-name",
		strings.NewReader(`{"Certificate": "helloOrg","Ip": "0.0.0.0"}`))

	m := &amocks.AsrRecordServiceClient{}
	org := &pb.AddOrgRequest{
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
	}
	m.On("UpdatePackage", mock.Anything, org).Return(&pb.AddOrgResponse{}, nil)

	r := NewRouter(&Clients{
		a: client.NewAsrFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_GetSubscriber(t *testing.T) {
	t.Run("SubscriberExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/v1/orgs/org-name",
			strings.NewReader(`{"Certificate": "updated_certs","Ip": "127.0.0.1"}`))

		m := &amocks.AsrRecordServiceClient{}

		org := &pb.UpdateOrgRequest{
			OrgName:     "org-name",
			Certificate: "updated_certs",
			Ip:          "127.0.0.1",
		}
		m.On("Read", mock.Anything, org).Return(&pb.UpdateOrgResponse{}, nil)

		r := NewRouter(&Clients{
			a: client.NewAsrFromClient(m),
		}, routerConfig).f.Engine()
		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("SubscriberDoesn'tExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/v1/orgs/org-name",
			strings.NewReader(`{"Certificate": "updated_certs","Ip": "127.0.0.1"}`))

		m := &amocks.AsrRecordServiceClient{}
		org := &pb.UpdateOrgRequest{
			OrgName:     "org-name",
			Certificate: "updated_certs",
			Ip:          "127.0.0.1",
		}
		m.On("Read", mock.Anything, org).Return(&pb.UpdateOrgResponse{}, nil)

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
