package rest

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/ukama/ukamaX/common/ukama"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/client"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, NewClientsSet(&pkg.GrpcEndpoints{})).gin

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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, NewClientsSet(&pkg.GrpcEndpoints{})).gin

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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, &Clients{
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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, &Clients{
		Registry: client.NewRegistryFromClient(m),
	}).gin

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"org-name"`)
}

func TestGetNodes(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes", nil)
	req.Header.Set("token", "bearer 123")
	nodeId := ukama.NewVirtualNodeId("homenode")

	t.Run("HappyPath", func(t *testing.T) {
		m := &pbmocks.RegistryServiceClient{}
		m.On("GetNodes", mock.Anything, mock.Anything).Return(&pb.GetNodesResponse{
			Orgs: []*pb.NodesList{
				{
					Nodes: []*pb.Node{
						{NodeId: nodeId.String()},
					},
				},
			},
		}, nil)

		r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, &Clients{
			Registry: client.NewRegistryFromClient(m),
		}).gin

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"nodeId":"%s"`, nodeId.String()))
	})

	t.Run("NoNodesReturned", func(t *testing.T) {
		m := &pbmocks.RegistryServiceClient{}
		m.On("GetNodes", mock.Anything, mock.Anything).Return(&pb.GetNodesResponse{
			Orgs: []*pb.NodesList{},
		}, nil)

		r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, &Clients{
			Registry: client.NewRegistryFromClient(m),
		}).gin

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
		body := w.Body.String()
		assert.Contains(t, body, "orgs")
	})

	t.Run("NoNodesReturned", func(t *testing.T) {
		m := &pbmocks.RegistryServiceClient{}
		m.On("GetNodes", mock.Anything, mock.Anything).Return(&pb.GetNodesResponse{
			Orgs: []*pb.NodesList{},
		}, nil)

		r := NewRouter(123456, true, NewDebugAuthMiddleware(), defaultCors, &Clients{
			Registry: client.NewRegistryFromClient(m),
		}).gin

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
		body := w.Body.String()
		assert.Contains(t, body, "orgs")
	})

}
