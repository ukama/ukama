package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	userspb "github.com/ukama/ukama/services/cloud/users/pb/gen"
	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/services/common/ukama"

	"github.com/ukama/ukama/services/cloud/api-gateway/pkg/client"

	"github.com/ukama/ukama/services/cloud/api-gateway/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	pbmocks "github.com/ukama/ukama/services/cloud/network/pb/gen/mocks"
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

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(NewDebugAuthMiddleware(), NewClientsSet(&pkg.GrpcEndpoints{}), routerConfig).f.Engine()

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

	r := NewRouter(NewDebugAuthMiddleware(), NewClientsSet(&pkg.GrpcEndpoints{}), routerConfig).f.Engine()

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

	m := &pbmocks.RegistryServiceClient{}
	m.On("GetOrg", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
		Registry: client.NewRegistryFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	m.AssertExpectations(t)
}

func TestGetOrg(t *testing.T) {
	// arrange
	const orgName = "org-name"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orgs/"+orgName, nil)
	req.Header.Set("token", "bearer 123")

	m := &pbmocks.RegistryServiceClient{}
	m.On("GetOrg", mock.Anything, mock.Anything).Return(&pb.Organization{
		Name:  orgName,
		Owner: "owner",
	}, nil)

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
		Registry: client.NewRegistryFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"org-name"`)
}

func TestGetNodes(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/orgs/test-org/nodes", nil)
	req.Header.Set("token", "bearer 123")
	nodeId := ukama.NewVirtualNodeId("homenode")

	t.Run("GetNodesSucceeded", func(t *testing.T) {
		w := httptest.NewRecorder()
		m := &pbmocks.RegistryServiceClient{}
		m.On("GetNodes", mock.Anything, mock.MatchedBy(func(r *pb.GetNodesRequest) bool {
			return r.OrgName == "test-org"
		})).Return(&pb.GetNodesResponse{
			Nodes: []*pb.Node{
				{NodeId: nodeId.String()}},
			OrgName: "test-org",
		}, nil)

		r := NewRouter(NewDebugAuthMiddleware(), &Clients{
			Registry: client.NewRegistryFromClient(m),
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
		m := &pbmocks.RegistryServiceClient{}
		m.On("GetNodes", mock.Anything, mock.Anything).Return(&pb.GetNodesResponse{
			Nodes: []*pb.Node{},
		}, nil)

		r := NewRouter(NewDebugAuthMiddleware(), &Clients{
			Registry: client.NewRegistryFromClient(m),
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
	m := &pbmocks.RegistryServiceClient{}
	m.On("GetNode", mock.Anything, mock.Anything).Return(&pb.GetNodeResponse{
		Node: &pb.Node{
			NodeId: nodeId.String(),
			Attached: []*pb.Node{
				{NodeId: attachedNodeId},
				{NodeId: "nodeId2"},
			}},
	}, nil)

	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
		Registry: client.NewRegistryFromClient(m),
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

func TestAddUpdateNode(t *testing.T) {
	// arrange

	tests := []struct {
		name             string
		isCreated        bool
		expectedHttpCode int
	}{
		{
			name:             "created",
			isCreated:        true,
			expectedHttpCode: http.StatusCreated,
		},
		{
			name:             "updated",
			isCreated:        false,
			expectedHttpCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			nodeId := ukama.NewVirtualNodeId("homenode").String()
			req, _ := http.NewRequest("PUT", "/orgs/test-org/nodes/"+nodeId, strings.NewReader(`{ "name": "test-name" }`))
			req.Header.Set("token", "bearer 123")

			w := httptest.NewRecorder()
			m := &pbmocks.RegistryServiceClient{}

			if test.isCreated {
				m.On("GetNode", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, ""))
				m.On("AddNode", mock.Anything, mock.Anything).Return(&pb.AddNodeResponse{}, nil)
			} else {
				m.On("GetNode", mock.Anything, mock.Anything).Return(&pb.GetNodeResponse{
					Node: &pb.Node{NodeId: nodeId}}, nil)
				m.On("UpdateNode", mock.Anything, mock.Anything).Return(&pb.UpdateNodeResponse{}, nil)
			}

			r := NewRouter(NewDebugAuthMiddleware(), &Clients{
				Registry: client.NewRegistryFromClient(m),
			}, routerConfig).f.Engine()

			// act
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, test.expectedHttpCode, w.Code)
			m.AssertExpectations(t)
		})
	}
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
