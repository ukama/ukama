package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/init/node-gateway/pkg"
	"github.com/ukama/ukama/systems/init/node-gateway/pkg/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	lmocks "github.com/ukama/ukama/systems/init/lookup/pb/gen/mocks"
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
	})
}

func TestRouter_GetNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	nodeReq := &pb.GetNodeRequest{
		NodeId: nodeId,
	}

	m.On("GetNode", mock.Anything, nodeReq).Return(&pb.GetNodeResponse{
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

func TestRouter_GetNode_NotFound(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	m.On("GetNode", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "node not found"))

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	m.AssertExpectations(t)
}
