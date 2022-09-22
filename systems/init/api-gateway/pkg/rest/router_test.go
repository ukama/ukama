package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/services/common/ukama"

	"github.com/ukama/ukama/systems/init/api-gateway/pkg"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	lmocks "github.com/ukama/ukama/systems/init/lookup/pb/gen/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestRouter_GetOrg_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lookup/orgs/org-name", nil)

	m := &lmocks.LookupServiceClient{}

	m.On("GetOrg", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetOrg(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lookup/orgs/org-name", nil)

	m := &lmocks.LookupServiceClient{}

	m.On("GetOrg", mock.Anything, mock.Anything).Return(&pb.GetOrgResponse{
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
	}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"orgName":"org-name"`)
}

func TestRouter_AddOrg(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/lookup/orgs/org-name",
		strings.NewReader(`{"Certificate": "helloOrg","Ip": "0.0.0.0"}`))

	m := &lmocks.LookupServiceClient{}

	org := &pb.AddOrgRequest{
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
	}
	m.On("AddOrg", mock.Anything, org).Return(&pb.AddOrgResponse{}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_GetNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lookup/orgs/org-name/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	nodeReq := &pb.GetNodeForOrgRequest{
		OrgName: "org-name",
		NodeId:  nodeId,
	}

	m.On("GetNodeForOrg", mock.Anything, nodeReq).Return(&pb.GetNodeResponse{
		NodeId:      nodeId,
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
	}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), strings.ToLower(nodeId))
}

func TestRouter_AddNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/lookup/orgs/org-name/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	nodeReq := &pb.AddNodeRequest{
		OrgName: "org-name",
		NodeId:  nodeId,
	}

	m.On("AddNodeForOrg", mock.Anything, nodeReq).Return(&pb.AddNodeResponse{}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_DeleteNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/lookup/orgs/org-name/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	nodeReq := &pb.DeleteNodeRequest{
		OrgName: "org-name",
		NodeId:  nodeId,
	}

	m.On("DeleteNodeForOrg", mock.Anything, nodeReq).Return(&pb.DeleteNodeResponse{}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetSystem(t *testing.T) {
	sys := "sys"
	sysId := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lookup/orgs/org-name/systems/"+sys, nil)

	m := &lmocks.LookupServiceClient{}

	sysReq := &pb.GetSystemRequest{
		OrgName:    "org-name",
		SystemName: sys,
	}

	m.On("GetSystemForOrg", mock.Anything, sysReq).Return(&pb.GetSystemResponse{
		SystemName:  sys,
		SystemId:    sysId,
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
		Port:        100,
	}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), strings.ToLower(sysId))
	assert.Contains(t, w.Body.String(), strings.ToLower(sys))
}

func TestRouter_AddSystem(t *testing.T) {
	sys := "sys"
	sysId := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/lookup/orgs/org-name/systems/"+sys,
		strings.NewReader(`{ "org":"org-name", "system":"sys", "ip":"0.0.0.0", "certificate":"certs", "port":100}`))

	m := &lmocks.LookupServiceClient{}

	sysReq := &pb.UpdateSystemRequest{
		SystemName:  sys,
		OrgName:     "org-name",
		Certificate: "certs",
		Ip:          "0.0.0.0",
		Port:        100,
	}

	m.On("UpdateSystemForOrg", mock.Anything, sysReq).Return(&pb.UpdateSystemResponse{
		SystemId: sysId,
	}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), strings.ToLower(sysId))
}

func TestRouter_DeleteSystem(t *testing.T) {
	sys := "sys"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/lookup/orgs/org-name/systems/"+sys, nil)

	m := &lmocks.LookupServiceClient{}

	sysReq := &pb.DeleteSystemRequest{
		OrgName:    "org-name",
		SystemName: sys,
	}

	m.On("DeleteSystemForOrg", mock.Anything, sysReq).Return(&pb.DeleteSystemResponse{}, nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}
