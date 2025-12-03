/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	nnmocks "github.com/ukama/ukama/systems/messaging/nns/pb/gen/mocks"
)

// Test data constants
const (
	testNodeId       = "uk-sa3333-uk-0001-0001"
	testNodeIp       = "192.168.1.100"
	testMeshIp       = "192.168.1.1"
	testNodePort     = int32(8080)
	testMeshPort     = int32(9090)
	testNetwork      = "test-network"
	testSite         = "test-site"
	testMeshHostName = "mesh-hostname"
)

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &pkg.HttpEndpoints{},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	n := &nnmocks.NnsClient{}
	testClientSet = &Clients{
		n: client.NewNnsFromClient(n),
	}
}

func TestRouter_PingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_GetNode(t *testing.T) {
	nodeId := testNodeId
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/nns/node/"+nodeId, nil)

	n := &nnmocks.NnsClient{}
	arc := &cmocks.AuthClient{}

	pReq := &pb.GetNodeRequest{
		NodeId: nodeId,
	}

	pResp := &pb.GetNodeResponse{
		NodeId:   nodeId,
		NodeIp:   testNodeIp,
		NodePort: testNodePort,
	}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	n.On("GetNode", mock.Anything, pReq, mock.Anything).Return(pResp, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(n),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), nodeId)
	assert.Contains(t, w.Body.String(), testNodeIp)
	n.AssertExpectations(t)
}

func TestRouter_GetMesh(t *testing.T) {
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/nns/mesh", nil)

	n := &nnmocks.NnsClient{}
	arc := &cmocks.AuthClient{}

	pReq := &pb.GetMeshRequest{}

	pResp := &pb.GetMeshResponse{
		MeshIp:   testMeshIp,
		MeshPort: testMeshPort,
	}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	n.On("GetMesh", mock.Anything, pReq, mock.Anything).Return(pResp, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(n),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testMeshIp)
	n.AssertExpectations(t)
}

func TestRouter_PutNode(t *testing.T) {
	nodeId := testNodeId
	ureq := SetNodeRequest{
		NodeId:       nodeId,
		NodeIp:       testNodeIp,
		MeshIp:       testMeshIp,
		NodePort:     testNodePort,
		MeshPort:     testMeshPort,
		Network:      testNetwork,
		Site:         testSite,
		MeshHostName: testMeshHostName,
	}

	jreq, err := json.Marshal(&ureq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("PUT", "/v1/nns/node", bytes.NewReader(jreq))
	hreq.Header.Set("Content-Type", "application/json")

	n := &nnmocks.NnsClient{}
	arc := &cmocks.AuthClient{}

	pReq := &pb.SetRequest{
		NodeId:       ureq.NodeId,
		NodeIp:       ureq.NodeIp,
		NodePort:     ureq.NodePort,
		MeshIp:       ureq.MeshIp,
		MeshPort:     ureq.MeshPort,
		Network:      ureq.Network,
		Site:         ureq.Site,
		MeshHostName: ureq.MeshHostName,
	}

	pResp := &pb.SetResponse{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	n.On("Set", mock.Anything, pReq, mock.Anything).Return(pResp, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(n),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	n.AssertExpectations(t)
}

func TestRouter_UpdateNode(t *testing.T) {
	nodeId := testNodeId
	ureq := UpdateNodeRequest{
		NodeId:   nodeId,
		NodeIp:   testNodeIp,
		NodePort: testNodePort,
	}

	jreq, err := json.Marshal(&ureq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("PUT", "/v1/nns/node/"+nodeId, bytes.NewReader(jreq))
	hreq.Header.Set("Content-Type", "application/json")

	n := &nnmocks.NnsClient{}
	arc := &cmocks.AuthClient{}

	pReq := &pb.UpdateNodeRequest{
		NodeId:   ureq.NodeId,
		NodeIp:   ureq.NodeIp,
		NodePort: ureq.NodePort,
	}

	pResp := &pb.UpdateNodeResponse{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	n.On("UpdateNode", mock.Anything, pReq, mock.Anything).Return(pResp, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(n),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	n.AssertExpectations(t)
}

func TestRouter_UpdateMesh(t *testing.T) {
	ureq := UpdateMeshRequest{
		MeshIp:   testMeshIp,
		MeshPort: testMeshPort,
	}

	jreq, err := json.Marshal(&ureq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("PUT", "/v1/nns/mesh", bytes.NewReader(jreq))
	hreq.Header.Set("Content-Type", "application/json")

	n := &nnmocks.NnsClient{}
	arc := &cmocks.AuthClient{}

	pReq := &pb.UpdateMeshRequest{
		MeshIp:   ureq.MeshIp,
		MeshPort: ureq.MeshPort,
	}

	pResp := &pb.UpdateMeshResponse{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	n.On("UpdateMesh", mock.Anything, pReq, mock.Anything).Return(pResp, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(n),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	n.AssertExpectations(t)
}

func TestRouter_DeleteNode(t *testing.T) {
	nodeId := testNodeId
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("DELETE", "/v1/nns/node/"+nodeId, nil)

	n := &nnmocks.NnsClient{}
	arc := &cmocks.AuthClient{}

	pReq := &pb.DeleteRequest{
		NodeId: nodeId,
	}

	pResp := &pb.DeleteResponse{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	n.On("Delete", mock.Anything, pReq, mock.Anything).Return(pResp, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(n),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	n.AssertExpectations(t)
}

func TestRouter_ListNodes(t *testing.T) {
	w := httptest.NewRecorder()
	hreq, _ := http.NewRequest("GET", "/v1/nns/list", nil)

	n := &nnmocks.NnsClient{}
	arc := &cmocks.AuthClient{}

	pReq := &pb.ListRequest{}

	pResp := &pb.ListResponse{
		List: []*pb.OrgMap{
			{
				NodeId:       testNodeId,
				NodeIp:       testNodeIp,
				NodePort:     testNodePort,
				MeshIp:       testMeshIp,
				MeshPort:     testMeshPort,
				Network:      testNetwork,
				Site:         testSite,
				MeshHostName: testMeshHostName,
			},
		},
	}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	n.On("List", mock.Anything, pReq, mock.Anything).Return(pResp, nil)

	r := NewRouter(&Clients{
		n: client.NewNnsFromClient(n),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, hreq)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testNodeId)
	assert.Contains(t, w.Body.String(), testNodeIp)
	n.AssertExpectations(t)
}
