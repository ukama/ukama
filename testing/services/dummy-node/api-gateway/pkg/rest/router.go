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

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/testing/services/dummy-node/api-gateway/cmd/version"
	"github.com/ukama/ukama/testing/services/dummy-node/api-gateway/pkg"
	"github.com/ukama/ukama/testing/services/dummy-node/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/testing/services/dummy-node/node/pb/gen"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	debugMode     bool
	serverConf    *rest.HttpConfig
	httpEndpoints *pkg.HttpEndpoints
}

type Clients struct {
	Node node
}

type node interface {
	ResetNode(id string) (*pb.ResetResponse, error)
	TurnNodeOff(id string) (*pb.TurnNodeOffResponse, error)
	TurnRFOn(id string) (*pb.NodeRFOnResponse, error)
	TurnRFOff(id string) (*pb.NodeRFOffResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Node = client.NewNodeService(endpoints.Node, endpoints.Timeout)

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
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, "")
	group := r.f.Group("/v1", "API gateway", "Dummy node system version v1")

	const dummy = "/dummy-node"
	d := group.Group(dummy, "Dummy node", "Operations on Dummy node")
	d.GET("/reset-node/:id", formatDoc("Reset node", "Reset dummy node by id"), tonic.Handler(r.resetNodeById, http.StatusOK))
	d.GET("/node-rf-off/:id", formatDoc("Node RF OFF", "Turn node rf off by id"), tonic.Handler(r.turnRFOffByid, http.StatusOK))
	d.GET("/node-rf-on/:id", formatDoc("Node RF ON", "Turn node rf on by id"), tonic.Handler(r.turnRFOnByid, http.StatusOK))
	d.GET("/node-off/:id", formatDoc("Turn node OFF", "Turn node off by id"), tonic.Handler(r.turnNodeOffByid, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) resetNodeById(c *gin.Context, req *ReqNodeId) (interface{}, error) {
	return r.clients.Node.ResetNode(req.Id)
}

func (r *Router) turnNodeOffByid(c *gin.Context, req *ReqNodeId) (interface{}, error) {
	return r.clients.Node.TurnNodeOff(req.Id)
}

func (r *Router) turnRFOffByid(c *gin.Context, req *ReqNodeId) (interface{}, error) {
	return r.clients.Node.TurnRFOff(req.Id)
}

func (r *Router) turnRFOnByid(c *gin.Context, req *ReqNodeId) (interface{}, error) {
	return r.clients.Node.TurnRFOn(req.Id)
}
