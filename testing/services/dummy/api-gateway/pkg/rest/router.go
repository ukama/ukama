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

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/cmd/version"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/pkg"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"

	"github.com/gin-gonic/gin"
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
	DController dcontroller
}
type dcontroller interface {
	Update(req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error)
	Start(req *pb.StartMetricsRequest) (*pb.StartMetricsResponse, error)
}



func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	dcontroller, err := client.NewController(endpoints.Controller, endpoints.Timeout)
	if err != nil {
		log.Fatalf("Failed to create controller client: %v", err)
	}
	c.DController = dcontroller
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

	group.GET("/ping", formatDoc("Ping API", "Ping to get status"), r.ping)
	dcontroller := group.Group("/controller", "Dummy controller service", "Dummy controller service")
	dcontroller.PUT("/update", formatDoc("Update metrics", "Update metrics"),tonic.Handler(r.updateHandler, http.StatusCreated))
	dcontroller.POST("/start", formatDoc("Start controller", "Start controller"), tonic.Handler(r.startHandler, http.StatusCreated))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (r *Router) updateHandler(c *gin.Context, req *UpdateReq) (*pb.UpdateMetricsResponse, error) {
	profileValue, ok := pb.Profile_value[req.Profile]
	if !ok {
		return nil, fmt.Errorf("invalid profile: %s", req.Profile)
	}

	scenarioValue, ok := pb.Scenario_value[req.Scenario]
	if !ok {
		return nil, fmt.Errorf("invalid scenario: %s", req.Scenario)
	}

	portUpdates := make([]*pb.PortUpdate, len(req.PortUpdates))
	for i, update := range req.PortUpdates {
		portUpdates[i] = &pb.PortUpdate{
			PortNumber: update.PortNumber,
			Status:     update.Status,
		}
	}

	return r.clients.DController.Update(&pb.UpdateMetricsRequest{
		SiteId:      req.SiteId,
		Profile:     pb.Profile(profileValue),
		Scenario:    pb.Scenario(scenarioValue),
		PortUpdates: portUpdates,
	})
}

func (r *Router) startHandler(c *gin.Context, req *StartReq) (*pb.StartMetricsResponse, error) {
	return r.clients.DController.Start(&pb.StartMetricsRequest{
		SiteId: req.SiteId,
		Profile: pb.Profile(pb.Profile_value[req.Profile]),
		Scenario: pb.Scenario(pb.Scenario_value[req.Scenario]),
	})
}
