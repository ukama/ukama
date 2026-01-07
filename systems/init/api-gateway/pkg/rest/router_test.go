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
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
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
		Lookup:  "localhost:8080",
	})
}

func TestRouter_PingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	arc := &cmocks.AuthClient{}
	r := NewRouter(testClientSet, routerConfig, arc.AuthenticateUser).f.Engine()

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_GetOrg_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name", nil)

	m := &lmocks.LookupServiceClient{}
	m.On("GetOrg", mock.Anything, mock.Anything).Return(nil, status.Error(codes.NotFound, "org not found"))

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetOrg(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name", nil)

	m := &lmocks.LookupServiceClient{}

	m.On("GetOrg", mock.Anything, mock.Anything).Return(&pb.GetOrgResponse{
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
	}, nil)

	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), `"orgName":"org-name"`)
}

func TestRouter_AddOrg(t *testing.T) {
	w := httptest.NewRecorder()
	id := uuid.NewV4().String()
	rd := AddOrgRequest{
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
		OrgId:       id,
	}

	jd, err := json.Marshal(&rd)
	assert.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/v1/orgs/org-name", bytes.NewReader(jd))

	m := &lmocks.LookupServiceClient{}

	org := &pb.AddOrgRequest{
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
		OrgId:       id,
	}
	m.On("AddOrg", mock.Anything, org).Return(&pb.AddOrgResponse{}, nil)

	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_UpdateOrg(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/orgs/org-name",
		strings.NewReader(`{"Certificate": "updated_certs","Ip": "127.0.0.1"}`))

	m := &lmocks.LookupServiceClient{}

	org := &pb.UpdateOrgRequest{
		OrgName:     "org-name",
		Certificate: "updated_certs",
		Ip:          "127.0.0.1",
	}

	m.On("UpdateOrg", mock.Anything, org).Return(&pb.UpdateOrgResponse{}, nil)

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}

func TestRouter_GetNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	nodeReq := &pb.GetNodeForOrgRequest{
		OrgName: "org-name",
		NodeId:  nodeId,
	}

	m.On("GetNodeForOrg", mock.Anything, nodeReq).Return(&pb.GetNodeResponse{
		NodeId:      nodeId,
		OrgName:     "org-name",
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
	}, nil)

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), strings.ToLower(nodeId))
}

func TestRouter_AddNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/orgs/org-name/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	nodeReq := &pb.AddNodeRequest{
		OrgName: "org-name",
		NodeId:  nodeId,
	}

	m.On("AddNodeForOrg", mock.Anything, nodeReq).Return(&pb.AddNodeResponse{}, nil)

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_DeleteNode(t *testing.T) {
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/orgs/org-name/nodes/"+nodeId, nil)

	m := &lmocks.LookupServiceClient{}

	nodeReq := &pb.DeleteNodeRequest{
		OrgName: "org-name",
		NodeId:  nodeId,
	}

	m.On("DeleteNodeForOrg", mock.Anything, nodeReq).Return(&pb.DeleteNodeResponse{}, nil)

	arc := &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_GetSystem(t *testing.T) {
	sys := "sys"
	sysId := uuid.NewV4().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/orgs/org-name/systems/"+sys, nil)

	m := &lmocks.LookupServiceClient{}

	sysReq := &pb.GetSystemRequest{
		OrgName:    "org-name",
		SystemName: sys,
	}

	m.On("GetSystemForOrg", mock.Anything, sysReq).Return(&pb.GetSystemResponse{
		SystemName:  sys,
		SystemId:    sysId,
		Certificate: "helloOrg",
		Ip:          "0.0.0.0",
		Port:        100,
	}, nil)

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), strings.ToLower(sysId))
	assert.Contains(t, w.Body.String(), strings.ToLower(sys))
}

func TestRouter_AddSystem(t *testing.T) {
	sys := "sys"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/orgs/org-name/systems/"+sys,
		strings.NewReader(`{ "org":"org-name", "system":"sys", "apiGwIp":"0.0.0.0", "certificate":"certs", "apiGwPort":100, "apiGwUrl":"http://localhost:8080"}`))

	m := &lmocks.LookupServiceClient{}

	sysReq := &pb.AddSystemRequest{
		SystemName:  sys,
		OrgName:     "org-name",
		Certificate: "certs",
		Ip:          "0.0.0.0",
		Port:        100,
		Url:         "http://localhost:8080",
	}

	m.On("AddSystemForOrg", mock.Anything, sysReq).Return(&pb.AddSystemResponse{}, nil)

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_UpdateSystem(t *testing.T) {
	sys := "sys"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/orgs/org-name/systems/"+sys,
		strings.NewReader(`{ "ip":"127.0.0.1", "certificate":"updated_certs", "port":101}`))

	m := &lmocks.LookupServiceClient{}

	sysReq := &pb.UpdateSystemRequest{
		SystemName:  sys,
		OrgName:     "org-name",
		Certificate: "updated_certs",
		Ip:          "127.0.0.1",
		Port:        101,
	}

	m.On("UpdateSystemForOrg", mock.Anything, sysReq).Return(&pb.UpdateSystemResponse{}, nil)

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRouter_DeleteSystem(t *testing.T) {
	sys := "sys"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/orgs/org-name/systems/"+sys, nil)

	m := &lmocks.LookupServiceClient{}

	sysReq := &pb.DeleteSystemRequest{
		OrgName:    "org-name",
		SystemName: sys,
	}

	m.On("DeleteSystemForOrg", mock.Anything, sysReq).Return(&pb.DeleteSystemResponse{}, nil)

	arc := &cmocks.AuthClient{}
	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	r := NewRouter(&Clients{
		l: client.NewLookupFromClient(m),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}
