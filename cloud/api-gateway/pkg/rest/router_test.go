package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"ukamaX/cloud/api-gateway/pkg"
	"ukamaX/cloud/api-gateway/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/generated"
	pbmocks "github.com/ukama/ukamaX/cloud/registry/pb/generated/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), NewClientsSet(&pkg.GrpcEndpoints{})).gin

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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), NewClientsSet(&pkg.GrpcEndpoints{})).gin

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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), &Clients{
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

	r := NewRouter(123456, true, NewDebugAuthMiddleware(), &Clients{
		Registry: client.NewRegistryFromClient(m),
	}).gin

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"name":"org-name"`)
}
