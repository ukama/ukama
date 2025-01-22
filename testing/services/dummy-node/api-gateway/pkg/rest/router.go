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
	Get() (string, error)
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
	group := r.f.Group("/v1", "API gateway", "Inventory system version v1")

	const dummy = "/dummy-node"
	d := group.Group(dummy, "Dummy node", "Operations on Dummy node", nil)
	d.GET("/", formatDoc("Get Dummy node", "Get dummy node id"), tonic.Handler(r.getDummyNodeId, http.StatusOK))

}

func (r *Router) getDummyNodeId(c *gin.Context, req *GetRequest) (string, error) {
	return r.clients.Node.Get()
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
