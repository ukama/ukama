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
	"io"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/node/api-gateway/cmd/version"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/node/api-gateway/pkg"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
	cfgPb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	contPb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	spb "github.com/ukama/ukama/systems/node/software/pb/gen"
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
	Controller      controller
	Configurator    configurator
	SoftwareManager softwareManager
}

type controller interface {
	RestartSite(siteName, networkId string) (*contPb.RestartSiteResponse, error)
	RestartNode(nodeId string) (*contPb.RestartNodeResponse, error)
	RestartNodes(networkId string, nodeIds []string) (*contPb.RestartNodesResponse, error)
	ToggleInternetSwitch(status bool, port int32, siteId string) (*contPb.ToggleInternetSwitchResponse, error)
	PingNode(*contPb.PingNodeRequest) (*contPb.PingNodeResponse, error)
}

type configurator interface {
	ConfigEvent(b []byte) (*cfgPb.ConfigStoreEventResponse, error)
	ApplyConfig(commit string) (*cfgPb.ApplyConfigResponse, error)
	GetConfigVersion(nodeId string) (*cfgPb.ConfigVersionResponse, error)
}

type softwareManager interface {
	UpdateSoftware(space string, name string, tag string, nodeId string) (*spb.UpdateSoftwareResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Controller = client.NewController(endpoints.Controller, endpoints.Timeout)
	c.Configurator = client.NewConfigurator(endpoints.Configurator, endpoints.Timeout)
	c.SoftwareManager = client.NewSoftwareManager(endpoints.Software, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init(authfunc)
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.HttpServices,
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

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "API gateway", "node system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if err == nil {
			return
		}
	})
	auth.Use()
	{
		const cont = "/controller"
		controller := auth.Group(cont, "Controller", "Operations on controllers")
		controller.POST("/networks/:network_id/sites/:site_name/restart", formatDoc("Restart a site in an organization", "Restarting a site within an organization"), tonic.Handler(r.postRestartSiteHandler, http.StatusOK))
		controller.POST("/nodes/:node_id/restart", formatDoc("Restart a node", "Restarting a node"), tonic.Handler(r.postRestartNodeHandler, http.StatusOK))
		controller.POST("/networks/:network_id/restart-nodes", formatDoc("Restart multiple nodes within a network", "Restarting multiple nodes within a network"), tonic.Handler(r.postRestartNodesHandler, http.StatusOK))
		controller.POST("/sites/:site_id/toggle-internet-port", formatDoc("Toggle internet port for a site", "Turns the internet port on or off for a specific site"), tonic.Handler(r.postToggleInternetSwitchHandler, http.StatusOK))

		controller.POST("/nodes/:node_id/ping", formatDoc("Ping a node", "Ping a node"), tonic.Handler(r.postPingNodeHandler, http.StatusAccepted))

		const cfg = "/configurator"
		cfgS := auth.Group(cfg, "Configurator", "Config for nodes")
		cfgS.POST("/config", formatDoc("Event in config store", "push event has happened in config store"), tonic.Handler(r.postConfigEventHandler, http.StatusAccepted))
		cfgS.POST("/config/apply/:commit", formatDoc("Apply config version ", "Updated nodes to version"), tonic.Handler(r.postConfigApplyVersionHandler, http.StatusAccepted))
		cfgS.GET("/config/node/:node_id", formatDoc("Current ruunning config", "Read the cuurrent running version and status"), tonic.Handler(r.getRunningConfigVersionHandler, http.StatusOK))

		const soft = "/software"
		softS := auth.Group(soft, "Software manager", "Operations on software")
		softS.POST("/update/:space/:name/:tag/:node_id", formatDoc("Update software", "Update software"), tonic.Handler(r.postUpdateSoftwareHandler, http.StatusOK))
	}
}

func (r *Router) postPingNodeHandler(c *gin.Context, req *PingNodeRequest) (*contPb.PingNodeResponse, error) {
	return r.clients.Controller.PingNode(&contPb.PingNodeRequest{
		NodeId:    req.NodeId,
		RequestId: req.RequestId,
		Message:   req.Message,
		Timestamp: req.TimeStamp,
	})
}

func (r *Router) postRestartNodeHandler(c *gin.Context, req *RestartNodeRequest) (*contPb.RestartNodeResponse, error) {
	return r.clients.Controller.RestartNode(req.NodeId)
}

func (r *Router) postRestartSiteHandler(c *gin.Context, req *RestartSiteRequest) (*contPb.RestartSiteResponse, error) {
	return r.clients.Controller.RestartSite(req.SiteName, req.NetworkId)
}

func (r *Router) postUpdateSoftwareHandler(c *gin.Context, req *UpdateSoftwareRequest) (*spb.UpdateSoftwareResponse, error) {
	return r.clients.SoftwareManager.UpdateSoftware(req.Space, req.Name, req.Tag, req.NodeId)
}

func (r *Router) postConfigEventHandler(c *gin.Context) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Errorf("Failed to decode config event.Error %s", err.Error())
		return err
	}
	log.Infof("received config event with %+v", string(body))

	_, err = r.clients.Configurator.ConfigEvent(body)
	if err != nil {
		log.Errorf("Failed to configure nodes.Error %s", err.Error())
		return err
	}

	return nil
}

func (r *Router) postConfigApplyVersionHandler(c *gin.Context, req *ApplyConfigRequest) error {

	log.Infof("received apply config with %+v", req)

	_, err := r.clients.Configurator.ApplyConfig(req.Commit)
	if err != nil {
		log.Errorf("Failed to apply config version %s to nodes.Error %s", req.Commit, err.Error())
		return err
	}

	log.Infof("received apply config with %+v", req)

	_, err = r.clients.Configurator.ApplyConfig(req.Commit)
	if err != nil {
		log.Errorf("Failed to apply config version %s to nodes.Error %s", req.Commit, err.Error())
		return err
	}

	return nil
}

func (r *Router) getRunningConfigVersionHandler(c *gin.Context, req *GetConfigVersionRequest) (*cfgPb.ConfigVersionResponse, error) {

	log.Infof("Received get running config version.")

	cfg, err := r.clients.Configurator.GetConfigVersion(req.NodeId)
	if err != nil {
		log.Errorf("Failed to get config version for node %s.Error %s", req.NodeId, err.Error())
		return nil, err
	}

	return cfg, nil
}

func (r *Router) postRestartNodesHandler(c *gin.Context, req *RestartNodesRequest) (*contPb.RestartNodesResponse, error) {
	return r.clients.Controller.RestartNodes(req.NetworkId, req.NodeIds)
}
func (r *Router) postToggleInternetSwitchHandler(c *gin.Context, req *ToggleInternetSwitchRequest) (*contPb.ToggleInternetSwitchResponse, error) {
	return r.clients.Controller.ToggleInternetSwitch(req.Status, req.Port, req.SiteId)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
