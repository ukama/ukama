package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	hsspb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"github.com/ukama/ukamaX/common/config"
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

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, config.Metrics{}, NewClientsSet(&pkg.GrpcEndpoints{})).gin

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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, config.Metrics{}, NewClientsSet(&pkg.GrpcEndpoints{})).gin

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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, config.Metrics{}, &Clients{
		Registry: client.NewRegistryFromClient(m),
	}).gin

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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, config.Metrics{}, &Clients{
		Registry: client.NewRegistryFromClient(m),
	}).gin

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name": "org-name"`)
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
			Orgs: []*pb.NodesList{
				{
					Nodes: []*pb.Node{
						{NodeId: nodeId.String()},
					},
				},
			},
		}, nil)

		r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, config.Metrics{}, &Clients{
			Registry: client.NewRegistryFromClient(m),
		}).gin

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"nodeId": "%s"`, nodeId.String()))
	})

	t.Run("NoNodesReturned", func(t *testing.T) {
		w := httptest.NewRecorder()
		m := &pbmocks.RegistryServiceClient{}
		m.On("GetNodes", mock.Anything, mock.Anything).Return(&pb.GetNodesResponse{
			Orgs: []*pb.NodesList{},
		}, nil)

		r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, config.Metrics{}, &Clients{
			Registry: client.NewRegistryFromClient(m),
		}).gin

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
	const imsi = "0000010000000001"

	m := hssmocks.UserServiceClient{}
	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, config.Metrics{}, &Clients{
		Hss: client.NewTestHssFromClient(&m),
	}).gin

	body, err := json.Marshal(&hsspb.User{FirstName: firstName, Uuid: userUuid, Imsi: imsi})
	if err != nil {
		panic(err)
	}

	// tests go here
	t.Run("AddUser", func(t *testing.T) {
		m = hssmocks.UserServiceClient{}
		m.On("Add", mock.Anything, mock.MatchedBy(func(r *hsspb.AddUserRequest) bool {
			return r.User.FirstName == firstName && r.User.Imsi == imsi
		})).Return(&hsspb.AddUserResponse{
			User: &hsspb.User{
				FirstName: firstName,
				Uuid:      userUuid,
				Imsi:      imsi,
			},
		}, nil)

		req, _ := http.NewRequest("POST", "/orgs/"+orgName+"/users", bytes.NewReader(body))
		req.Header.Set("token", "bearer 123")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)

		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"uuid": "%s"`, userUuid))
		m.AssertExpectations(t)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		m = hssmocks.UserServiceClient{}
		m.On("Delete", mock.Anything, mock.MatchedBy(func(r *hsspb.DeleteUserRequest) bool {
			return r.UserUuid == userUuid && r.Org == orgName
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
					FirstName: firstName,
					Imsi:      userUuid,
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
