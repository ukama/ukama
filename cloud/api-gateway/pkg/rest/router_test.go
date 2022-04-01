package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	hsspb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"github.com/ukama/ukamaX/common/rest"
	"github.com/ukama/ukamaX/common/ukama"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/client"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	hssmocks "github.com/ukama/ukamaX/cloud/hss/pb/gen/mocks"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	pbmocks "github.com/ukama/ukamaX/cloud/registry/pb/gen/mocks"
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

	t.Run("HappyPath", func(t *testing.T) {
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

func Test_HssMethods(t *testing.T) {
	// arrange
	const orgName = "org-name"
	const userUuid = "93fcb344-c752-411d-9506-e27417224920"
	const firstName = "Joe"
	const simToken = "0000010000000001"

	m := hssmocks.UserServiceClient{}
	r := NewRouter(NewDebugAuthMiddleware(), &Clients{
		Hss: client.NewTestHssFromClient(&m),
	}, routerConfig).f.Engine()

	body, err := json.Marshal(UserRequest{Name: firstName, SimToken: simToken})
	if err != nil {
		panic(err)
	}

	// tests go here
	t.Run("AddUser", func(t *testing.T) {
		m = hssmocks.UserServiceClient{}
		m.On("Add", mock.Anything, mock.MatchedBy(func(r *hsspb.AddRequest) bool {
			return r.User.Name == firstName && r.SimToken == simToken
		})).Return(&hsspb.AddResponse{
			User: &hsspb.User{
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
		m = hssmocks.UserServiceClient{}
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
		m = hssmocks.UserServiceClient{}
		m.On("Delete", mock.Anything, mock.MatchedBy(func(r *hsspb.DeleteUserRequest) bool {
			return r.Uuid == userUuid
		})).Return(&hsspb.DeleteUserResponse{}, nil)
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
		m = hssmocks.UserServiceClient{}
		m.On("List", mock.Anything, mock.MatchedBy(func(r *hsspb.ListUsersRequest) bool {
			return r.Org == orgName
		})).Return(&hsspb.ListUsersResponse{
			Org: orgName,
			Users: []*hsspb.User{
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
}
