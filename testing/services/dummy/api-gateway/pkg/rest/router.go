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
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/cmd/version"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/pkg"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
	logger  *logrus.Logger
}

type RouterConfig struct {
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

type Clients struct {
	Dsubscriber dsubscriber
}

type dsubscriber interface {
	Update(req *pb.UpdateRequest) (*pb.UpdateResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Dsubscriber = client.NewDsubscriber(endpoints.Dsubscriber, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig) *Router {
	r := &Router{
		clients: clients,
		config:  config,
		logger:  log.New(),
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
	endpoint := r.f.Group("/v1", "API gateway", "Dummy system version v1")
	endpoint.GET("/ping", formatDoc("Ping the server", "Returns a response indicating that the server is running."), tonic.Handler(r.pingHandler, http.StatusOK))

	health := endpoint.Group("/dsubscriber", "Dsubscriber", "Dummy subscriber service")
	health.PUT("/update", formatDoc("Update dsubscriber coroutine", "Update dsubscriber coroutine for specific subscriber."), tonic.Handler(r.updateHandler, http.StatusCreated))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) pingHandler(c *gin.Context) error {
	response := make(map[string]string)
	response["status"] = pkg.SystemName + " is running"
	response["version"] = version.Version
	c.JSON(http.StatusOK, response)

	return nil
}

func (r *Router) updateHandler(c *gin.Context, req *UpdateReq) (*pb.UpdateResponse, error) {
	return r.clients.Dsubscriber.Update(&pb.UpdateRequest{
		Dsubscriber: &pb.Dsubscriber{
			Iccid:   req.Iccid,
			Profile: req.Profile,
		}})
}
