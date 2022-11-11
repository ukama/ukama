package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"

	"github.com/ukama/ukama/systems/registry/api-gateway/pkg/client"

	"github.com/ukama/ukama/systems/registry/api-gateway/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	netmocks "github.com/ukama/ukama/systems/registry/network/pb/gen/mocks"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	nodemocks "github.com/ukama/ukama/systems/registry/node/pb/gen/mocks"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	orgmocks "github.com/ukama/ukama/systems/registry/org/pb/gen/mocks"
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

func TestPingRoute(t *testing.T) {
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

func TestGetOrg_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name", nil)
	req.Header.Set("token", "bearer 123")

	n := &netmocks.NetworkServiceClient{}

	o := &orgmocks.OrgServiceClient{}
	o.On("Get", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))
	nd := &nodemocks.NodeServiceClient{}

	r := NewRouter(&Clients{
		Registry: client.NewRegistryFromClient(n, o, nd),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	n.AssertExpectations(t)
}

func TestGetOrg(t *testing.T) {
	// arrange
	const orgName = "org-name"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/"+orgName, nil)
	req.Header.Set("token", "bearer 123")

	n := &netmocks.NetworkServiceClient{}

	o := &orgmocks.OrgServiceClient{}
	nd := &nodemocks.NodeServiceClient{}

	o.On("Get", mock.Anything, mock.Anything).Return(&orgpb.GetResponse{
		Org: &orgpb.Organization{
			Name:  orgName,
			Owner: "owner",
		},
	}, nil)

	r := NewRouter(&Clients{
		Registry: client.NewRegistryFromClient(n, o, nd),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	o.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"org-name"`)
}

func TestGetNodes(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1/orgs/test-org/nodes", nil)
	req.Header.Set("token", "bearer 123")
	nodeId := ukama.NewVirtualNodeId("homenode")

	t.Run("GetNodesSucceeded", func(t *testing.T) {
		w := httptest.NewRecorder()
		m := &netmocks.NetworkServiceClient{}
		o := &orgmocks.OrgServiceClient{}
		nd := &nodemocks.NodeServiceClient{}
		m.On("GetNodes", mock.Anything, mock.MatchedBy(func(r *netpb.GetNodesRequest) bool {
			return r.OrgName == "test-org"
		})).Return(&netpb.GetNodesResponse{
			Nodes: []*netpb.Node{
				{NodeId: nodeId.String()}},
			OrgName: "test-org",
		}, nil)

		r := NewRouter(&Clients{
			Registry: client.NewRegistryFromClient(m, o, nd),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"nodeId":"%s"`, nodeId.String()))
	})

	t.Run("NoNodesReturned", func(t *testing.T) {
		w := httptest.NewRecorder()
		m := &netmocks.NetworkServiceClient{}
		o := &orgmocks.OrgServiceClient{}
		nd := &nodemocks.NodeServiceClient{}
		m.On("GetNodes", mock.Anything, mock.Anything).Return(&netpb.GetNodesResponse{
			Nodes: []*netpb.Node{},
		}, nil)

		r := NewRouter(&Clients{
			Registry: client.NewRegistryFromClient(m, o, nd),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
		body := w.Body.String()
		assert.Contains(t, body, "nodes")
	})

}

func TestGetNode(t *testing.T) {

	const attachedNodeId = "nodeId1"

	// arrange
	nodeId := ukama.NewVirtualNodeId("homenode")
	req, _ := http.NewRequest("GET", "/v1/orgs/test-org/nodes/"+nodeId.String(), nil)
	req.Header.Set("token", "bearer 123")

	w := httptest.NewRecorder()
	m := &netmocks.NetworkServiceClient{}
	o := &orgmocks.OrgServiceClient{}
	nd := &nodemocks.NodeServiceClient{}
	nd.On("GetNode", mock.Anything, mock.Anything).Return(&nodepb.GetNodeResponse{
		Node: &nodepb.Node{
			NodeId: nodeId.String(),
			Attached: []*nodepb.Node{
				{NodeId: attachedNodeId},
				{NodeId: "nodeId2"},
			}},
	}, nil)

	r := NewRouter(&Clients{
		Registry: client.NewRegistryFromClient(m, o, nd),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	node := NodeExtended{}
	err := json.Unmarshal(w.Body.Bytes(), &node)
	if assert.NoError(t, err) {
		assert.Equal(t, nodeId.StringLowercase(), node.NodeId)
		assert.Equal(t, 2, len(node.Attached))
		assert.Equal(t, attachedNodeId, node.Attached[0].NodeId)
	}

}

func TestAddNode(t *testing.T) {
	// arrange
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	req, _ := http.NewRequest("PUT", "/v1/orgs/test-org/nodes/"+nodeId, strings.NewReader(`{ "name": "test-name" }`))
	req.Header.Set("token", "bearer 123")

	w := httptest.NewRecorder()
	net := &netmocks.NetworkServiceClient{}
	o := &orgmocks.OrgServiceClient{}
	nd := &nodemocks.NodeServiceClient{}

	nd.On("GetNode", mock.Anything, mock.Anything).Return(&nodepb.GetNodeResponse{
		Node: &nodepb.Node{
			NodeId: nodeId,
			Name:   "test-name",
		}}, nil)
	nd.On("AddNode", mock.Anything, mock.Anything).Return(&nodepb.AddNodeResponse{
		Node: &nodepb.Node{
			NodeId: nodeId,
			Name:   "test-name",
		},
	}, nil)

	net.On("AddNode", mock.Anything, mock.Anything).Return(&netpb.AddNodeResponse{}, nil)

	r := NewRouter(&Clients{
		Registry: client.NewRegistryFromClient(net, o, nd),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	fmt.Printf("Response: %s\n", w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code)
	net.AssertExpectations(t)

}

func Test_UpdateNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	req, _ := http.NewRequest("PATCH", "/v1/orgs/test-org/nodes/"+nodeId, strings.NewReader(`{ "name": "test-name" }`))
	req.Header.Set("token", "bearer 123")

	w := httptest.NewRecorder()
	net := &netmocks.NetworkServiceClient{}
	o := &orgmocks.OrgServiceClient{}
	nd := &nodemocks.NodeServiceClient{}

	nd.On("GetNode", mock.Anything, mock.Anything).Return(&nodepb.GetNodeResponse{
		Node: &nodepb.Node{NodeId: nodeId}}, nil)
	nd.On("UpdateNode", mock.Anything, mock.Anything).Return(&nodepb.UpdateNodeResponse{}, nil)

	r := NewRouter(&Clients{
		Registry: client.NewRegistryFromClient(net, o, nd),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	fmt.Printf("Response: %s\n", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	net.AssertExpectations(t)
}
