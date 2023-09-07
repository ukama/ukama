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
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"

	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	nmocks "github.com/ukama/ukama/systems/messaging/nns/pb/gen/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		Timeout: 1 * time.Second,
		Nns:     "localhost:9090",
	})
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

func TestRouter_GetNodeIP(t *testing.T) {
	nodeId := ukama.NewVirtualHomeNodeId().String()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nns/node/"+nodeId, nil)
	arc := &providers.AuthRestClient{}
	m := &nmocks.NnsClient{}

	preq := &pb.GetNodeIPRequest{
		NodeId: nodeId,
	}
	m.On("Get", mock.Anything, preq).Return(&pb.GetNodeIPResponse{
		Ip: "1.1.1.1",
	}, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetNodeIPMapListRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nns/list", nil)
	arc := &providers.AuthRestClient{}
	m := &nmocks.NnsClient{}
	maps := []*pb.NodeIPMap{
		{
			NodeId: ukama.NewVirtualHomeNodeId().String(),
			NodeIp: "1.1.1.1",
		},
		{
			NodeId: ukama.NewVirtualHomeNodeId().String(),
			NodeIp: "1.1.1.2",
		},
		{
			NodeId: ukama.NewVirtualHomeNodeId().String(),
			NodeIp: "1.1.1.3",
		}}
	m.On("GetNodeIPMapList", mock.Anything, &pb.NodeIPMapListRequest{}).Return(&pb.NodeIPMapListResponse{
		Map: maps}, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, http.StatusOK, w.Code) {
		assert.Contains(t, w.Body.String(), strings.ToLower(maps[0].NodeId))
		assert.Contains(t, w.Body.String(), strings.ToLower(maps[0].NodeIp))
	}
	m.AssertExpectations(t)

}

func TestRouter_GetNodeOrgMapListRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nns/map", nil)
	arc := &providers.AuthRestClient{}
	m := &nmocks.NnsClient{}

	maps := []*pb.NodeOrgMap{
		{
			NodeId:     ukama.NewVirtualHomeNodeId().String(),
			NodeIp:     "1.1.1.1",
			NodePort:   1000,
			MeshPort:   2000,
			Org:        "ukama",
			Network:    "net",
			Domainname: "domain.name",
		},
		{
			NodeId:     ukama.NewVirtualHomeNodeId().String(),
			NodeIp:     "1.1.1.2",
			NodePort:   1000,
			MeshPort:   2000,
			Org:        "ukama",
			Network:    "net",
			Domainname: "domain.name",
		},
		{
			NodeId:     ukama.NewVirtualHomeNodeId().String(),
			NodeIp:     "1.1.1.3",
			NodePort:   1000,
			MeshPort:   2000,
			Org:        "ukama",
			Network:    "net",
			Domainname: "domain.name",
		}}
	m.On("GetNodeOrgMapList", mock.Anything, &pb.NodeOrgMapListRequest{}).Return(&pb.NodeOrgMapListResponse{
		Map: maps}, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, http.StatusOK, w.Code) {
		assert.Contains(t, w.Body.String(), strings.ToLower(maps[0].NodeId))
		assert.Contains(t, w.Body.String(), strings.ToLower(maps[0].NodeIp))
		assert.Contains(t, w.Body.String(), strings.ToLower(maps[0].Network))
		assert.Contains(t, w.Body.String(), strings.ToLower(maps[0].Org))
		assert.Contains(t, w.Body.String(), strings.ToLower(maps[0].Domainname))
	}
	m.AssertExpectations(t)

}

func TestRouter_PrometheusTargets(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/prometheus", nil)
	arc := &providers.AuthRestClient{}
	m := &nmocks.NnsClient{}
	NodeIds := []string{ukama.NewVirtualHomeNodeId().String(), ukama.NewVirtualHomeNodeId().String(), ukama.NewVirtualHomeNodeId().String()}
	nodeIpMap := []*pb.NodeIPMap{
		{
			NodeId: NodeIds[0],
			NodeIp: "1.1.1.1",
		},
		{
			NodeId: NodeIds[1],
			NodeIp: "1.1.1.2",
		},
		{
			NodeId: NodeIds[2],
			NodeIp: "1.1.1.3",
		}}

	nodeOrgMap := []*pb.NodeOrgMap{
		{
			NodeId:     NodeIds[0],
			NodeIp:     "1.1.1.1",
			NodePort:   1000,
			MeshPort:   2000,
			Org:        "ukama",
			Network:    "net",
			Domainname: "domain.name",
		},
		{
			NodeId:     NodeIds[1],
			NodeIp:     "1.1.1.2",
			NodePort:   1000,
			MeshPort:   2000,
			Org:        "ukama",
			Network:    "net",
			Domainname: "domain.name",
		},
		{
			NodeId:     NodeIds[2],
			NodeIp:     "1.1.1.3",
			NodePort:   1000,
			MeshPort:   2000,
			Org:        "ukama",
			Network:    "net",
			Domainname: "domain.name",
		}}

	m.On("GetNodeOrgMapList", mock.Anything, &pb.NodeOrgMapListRequest{}).Return(&pb.NodeOrgMapListResponse{
		Map: nodeOrgMap}, nil)
	m.On("GetNodeIPMapList", mock.Anything, &pb.NodeIPMapListRequest{}).Return(&pb.NodeIPMapListResponse{
		Map: nodeIpMap}, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, http.StatusOK, w.Code) {
		assert.Contains(t, w.Body.String(), strings.ToLower(NodeIds[0]))
		assert.Contains(t, w.Body.String(), strings.ToLower(NodeIds[1]))
		assert.Contains(t, w.Body.String(), strings.ToLower(NodeIds[2]))
	}
	m.AssertExpectations(t)

}

func TestRouter_SetNodeIP(t *testing.T) {
	w := httptest.NewRecorder()
	nodeId := ukama.NewVirtualHomeNodeId().String()
	hreq := SetNodeIPRequest{
		NodeId:   nodeId,
		NodeIp:   "2.2.2.2",
		MeshIp:   "1.1.1.1",
		NodePort: 1000,
		MeshPort: 2000,
		Org:      "ukama",
		Network:  "net",
	}

	b, err := json.Marshal(&hreq)
	assert.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/v1/nns/node/"+nodeId, bytes.NewReader(b))
	arc := &providers.AuthRestClient{}
	m := &nmocks.NnsClient{}

	preq := &pb.SetNodeIPRequest{
		NodeId:   hreq.NodeId,
		NodeIp:   hreq.NodeIp,
		NodePort: hreq.NodePort,
		MeshIp:   hreq.MeshIp,
		MeshPort: hreq.MeshPort,
		Org:      hreq.Org,
		Network:  hreq.Network,
	}
	m.On("Set", mock.Anything, preq).Return(&pb.SetNodeIPResponse{}, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)

	m.AssertExpectations(t)

}

func TestRouter_DeleteNodeIP(t *testing.T) {
	nodeId := ukama.NewVirtualHomeNodeId().String()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/nns/node/"+nodeId, nil)
	arc := &providers.AuthRestClient{}
	m := &nmocks.NnsClient{}

	preq := &pb.DeleteNodeIPRequest{
		NodeId: nodeId,
	}
	m.On("Delete", mock.Anything, preq).Return(&pb.DeleteNodeIPResponse{}, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}
