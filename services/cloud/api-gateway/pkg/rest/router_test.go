package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	userspb "github.com/ukama/ukama/services/cloud/users/pb/gen"
	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/services/common/ukama"

	"github.com/ukama/ukama/services/cloud/api-gateway/pkg/client"

	"github.com/ukama/ukama/services/cloud/api-gateway/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	netpb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	netmocks "github.com/ukama/ukama/services/cloud/network/pb/gen/mocks"
	nodepb "github.com/ukama/ukama/services/cloud/node/pb/gen"
	nodemocks "github.com/ukama/ukama/services/cloud/node/pb/gen/mocks"
	orgpb "github.com/ukama/ukama/services/cloud/org/pb/gen"
	orgmocks "github.com/ukama/ukama/services/cloud/org/pb/gen/mocks"
	usrmocks "github.com/ukama/ukama/services/cloud/users/pb/gen/mocks"
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
	r := NewRouter(NewDebugAuthMiddleware(), testClientSet, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestGetOrg_Unauthorized(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orgs/org-name", nil)

	r := NewRouter(NewDebugAuthMiddleware(), testClientSet, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetOrg_NotFound(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orgs/org-name", nil)
	req.Header.Set("token", "bearer 123")

	n := &netmocks.NetworkServiceClient{}

	o := &orgmocks.OrgServiceClient{}
	o.On("Get", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))
	nd := &nodemocks.NodeServiceClient{}

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
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
	req, _ := http.NewRequest("GET", "/orgs/"+orgName, nil)
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

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
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
	req, _ := http.NewRequest("GET", "/orgs/test-org/nodes", nil)
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

		r := NewRouter(NewDebugAuthMiddleware(), &Clients{
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

		r := NewRouter(NewDebugAuthMiddleware(), &Clients{
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
	req, _ := http.NewRequest("GET", "/orgs/test-org/nodes/"+nodeId.String(), nil)
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

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
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
	req, _ := http.NewRequest("PUT", "/orgs/test-org/nodes/"+nodeId, strings.NewReader(`{ "name": "test-name" }`))
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

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
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
	req, _ := http.NewRequest("PATCH", "/orgs/test-org/nodes/"+nodeId, strings.NewReader(`{ "name": "test-name" }`))
	req.Header.Set("token", "bearer 123")

	w := httptest.NewRecorder()
	net := &netmocks.NetworkServiceClient{}
	o := &orgmocks.OrgServiceClient{}
	nd := &nodemocks.NodeServiceClient{}

	nd.On("GetNode", mock.Anything, mock.Anything).Return(&nodepb.GetNodeResponse{
		Node: &nodepb.Node{NodeId: nodeId}}, nil)
	nd.On("UpdateNode", mock.Anything, mock.Anything).Return(&nodepb.UpdateNodeResponse{}, nil)

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
		Registry: client.NewRegistryFromClient(net, o, nd),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	fmt.Printf("Response: %s\n", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	net.AssertExpectations(t)
}

func Test_HssMethods(t *testing.T) {
	// arrange
	const orgName = "org-name"
	const userUuid = "93fcb344-c752-411d-9506-e27417224920"
	const firstName = "Joe"
	const simToken = "0000010000000001"

	m := usrmocks.UserServiceClient{}
	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
		User: client.NewTestHssFromClient(&m),
	}, routerConfig).f.Engine()

	body, err := json.Marshal(UserRequest{Name: firstName, SimToken: simToken})
	if err != nil {
		panic(err)
	}

	// tests go here
	t.Run("AddUser", func(t *testing.T) {
		m = usrmocks.UserServiceClient{}
		m.On("Add", mock.Anything, mock.MatchedBy(func(r *userspb.AddRequest) bool {
			return r.User.Name == firstName && r.SimToken == simToken
		})).Return(&userspb.AddResponse{
			User: &userspb.User{
				Name: firstName,
				Uuid: userUuid,
			},
			Iccid: "0000000000000000001",
		}, nil)

		req, _ := http.NewRequest("POST", "/orgs/"+orgName+"/users", bytes.NewReader(body))
		req.Header.Set("token", "bearer 123")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)

		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"uuid":"%s"`, userUuid))
		m.AssertExpectations(t)
	})

	t.Run("AddUserReturnsError", func(t *testing.T) {
		m = usrmocks.UserServiceClient{}
		m.On("Add", mock.Anything, mock.Anything).Return(nil, status.Error(codes.PermissionDenied, "some err"))

		req, _ := http.NewRequest("POST", "/orgs/"+orgName+"/users", bytes.NewReader(body))
		req.Header.Set("token", "bearer 123")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "some err")
	})

	t.Run("DeleteUser", func(t *testing.T) {
		m = usrmocks.UserServiceClient{}
		m.On("Delete", mock.Anything, mock.MatchedBy(func(r *userspb.DeleteRequest) bool {
			return r.UserId == userUuid
		})).Return(&userspb.DeleteResponse{}, nil)
		req, _ := http.NewRequest("DELETE", "/orgs/"+orgName+"/users/"+userUuid, nil)
		req.Header.Set("token", "bearer 123")
		w := httptest.NewRecorder()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("ListUser", func(t *testing.T) {
		m = usrmocks.UserServiceClient{}
		m.On("List", mock.Anything, mock.MatchedBy(func(r *userspb.ListRequest) bool {
			return r.Org == orgName
		})).Return(&userspb.ListResponse{
			Org: orgName,
			Users: []*userspb.User{
				{
					Name: firstName,
					Uuid: userUuid,
				},
			},
		}, nil)

		req, _ := http.NewRequest("GET", "/orgs/"+orgName+"/users", nil)
		req.Header.Set("token", "bearer 123")
		w := httptest.NewRecorder()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), userUuid)
		assert.Contains(t, w.Body.String(), orgName)
		assert.Contains(t, w.Body.String(), firstName)
		m.AssertExpectations(t)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		const changedName = "changed"
		m = usrmocks.UserServiceClient{}
		m.On("DeactivateUser", mock.Anything, mock.MatchedBy(func(r *userspb.DeactivateUserRequest) bool {
			return r.UserId == userUuid
		})).Return(&userspb.DeactivateUserResponse{}, nil)
		m.On("Update", mock.Anything, mock.MatchedBy(func(r *userspb.UpdateRequest) bool {
			return r.UserId == userUuid && r.User.Name == changedName
		})).Return(&userspb.UpdateResponse{User: &userspb.User{
			Name:          changedName,
			IsDeactivated: true,
		}}, nil)

		m.On("Get", mock.Anything, mock.MatchedBy(func(r *userspb.GetRequest) bool {
			return r.UserId == userUuid
		})).Return(&userspb.GetResponse{User: &userspb.User{
			Name:          changedName,
			IsDeactivated: true,
		}}, nil)
		updBody, err := json.Marshal(UpdateUserRequest{
			Name:          changedName,
			IsDeactivated: true,
		})
		if err != nil {
			assert.FailNow(t, "error marshaling request", err.Error())
		}

		req, _ := http.NewRequest("PATCH", "/orgs/"+orgName+"/users/"+userUuid, bytes.NewReader(updBody))
		req.Header.Set("token", "bearer 123")
		w := httptest.NewRecorder()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
}
