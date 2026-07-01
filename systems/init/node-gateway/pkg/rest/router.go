/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/init/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/init/node-gateway/pkg"
	"github.com/ukama/ukama/systems/init/node-gateway/pkg/client"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	rpb "github.com/ukama/ukama/systems/init/reflector/pb/gen"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
}

type Clients struct {
	Bootstrap client.BootstrapEP
	Reflector client.ReflectorEP
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Bootstrap = client.NewBootstrap(endpoints.Bootstrap, endpoints.Timeout)
	c.Reflector = client.NewReflector(endpoints.Reflector, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig) *Router {
	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.Http,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName+" - Bootstrap", version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	v1 := r.f.Group("/v1", "Node gateway system ", "Node gateway system version v1")

	nodes := v1.Group("nodes", "Nodes", "Looking for nodes credentials")
	nodes.GET("/:nodeId", formatDoc("Get Nodes Credentials", ""), tonic.Handler(r.getNodeHandler, http.StatusOK))

	reflector := v1.Group("reflector", "Reflector", "Reflector service")
	reflector.GET("/ping", formatDoc("Ping Reflector", ""), tonic.Handler(r.reflectorPingHandler, http.StatusOK))
	reflector.GET("/", formatDoc("Get Reflector", ""), tonic.Handler(r.reflectorGetHandler, http.StatusOK))
	reflector.POST("/download", formatDoc("Download from Reflector", ""), tonic.Handler(r.reflectorDownloadHandler, http.StatusOK))
	reflector.POST("/upload", formatDoc("Upload to Reflector", ""), tonic.Handler(r.reflectorUploadHandler, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*pb.GetNodeCredentialsResponse, error) {
	return r.clients.Bootstrap.GetNodeCredentials(&pb.GetNodeCredentialsRequest{
		Id: req.NodeId,
	})
}

func (r *Router) reflectorPingHandler(c *gin.Context, req *ReflectorPingRequest) (*rpb.PingResponse, error) {
	return r.clients.Reflector.Ping(&rpb.PingRequest{})
}

func (r *Router) reflectorGetHandler(c *gin.Context, req *ReflectorGetRequest) (*rpb.GetResponse, error) {
	return r.clients.Reflector.Get(&rpb.GetRequest{
		NodeId: req.NodeId,
	})
}

func (r *Router) reflectorDownloadHandler(c *gin.Context, req *ReflectorDownloadRequest) (*rpb.DownloadResponse, error) {
	return r.clients.Reflector.Download(&rpb.DownloadRequest{
		NodeId:       req.NodeId,
		Bytes:        req.Bytes,
		ChunkBytes:   req.ChunkBytes,
		ChunkDelayMs: req.ChunkDelayMs,
	})
}

func (r *Router) reflectorUploadHandler(c *gin.Context, req *ReflectorUploadRequest) (*rpb.UploadResponse, error) {
	return r.clients.Reflector.Upload(&rpb.UploadRequest{
		NodeId:  req.NodeId,
		Payload: req.Payload,
	})
}
