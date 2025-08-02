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

	"github.com/ukama/ukama/systems/api/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/config"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

var REDIRECT_URI = "https://subscriber.dev.ukama.com/swagger/#/"

type Router struct {
	f       *fizz.Fizz
	network client.Network
	pkg     client.Package
	sim     client.Sim
	node    client.Node
	config  *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *crest.HttpConfig
	auth       *config.Auth
}

func NewRouter(network client.Network, pkg client.Package, sim client.Sim, node client.Node,
	config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
	r := &Router{
		network: network,
		pkg:     pkg,
		sim:     sim,
		node:    node,
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
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		auth:       svcConf.Auth,
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
	r.f = crest.NewFizzRouter(r.config.serverConf, pkg.SystemName,
		version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	auth := r.f.Group("/v1", "Ukama API GW ", "API system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")

			return
		}

		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)

		err := f(ctx, r.config.auth.AuthServerUrl)
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
		// network routes
		networks := auth.Group("/networks", "Network", "Networks")
		networks.POST("", formatDoc("Create Network", "Create a new network"), tonic.Handler(r.postNetwork, http.StatusPartialContent))
		networks.GET("/:network_id", formatDoc("Get Network", "Get a specific network"), tonic.Handler(r.getNetwork, http.StatusOK))

		// package routes
		packages := auth.Group("/packages", "Package", "Packages")
		packages.POST("", formatDoc("Add Package", "Add a new package"), tonic.Handler(r.postPackage, http.StatusPartialContent))
		packages.GET("/:package_id", formatDoc("Get Package", "Get a specific package"), tonic.Handler(r.getPackage, http.StatusOK))

		// sim routes
		sims := auth.Group("/sims", "Sim", "sims")
		sims.POST("", formatDoc("Configure Sim", "Configure a new sim"), tonic.Handler(r.postSim, http.StatusPartialContent))
		sims.GET("/:id", formatDoc("Get Sim", "Get a specific sim"), tonic.Handler(r.getSim, http.StatusOK))

		// node routes
		nodes := auth.Group("/nodes", "Node", "Operations on Nodes")
		nodes.GET("/:node_id", formatDoc("Get Node", "Get a specific node"), tonic.Handler(r.getNode, http.StatusOK))
		nodes.POST("", formatDoc("Add Node", "Add a new Node to an organization"), tonic.Handler(r.postNode, http.StatusCreated))
		nodes.DELETE("/:node_id", formatDoc("Delete Node", "Remove node from org"), tonic.Handler(r.deleteNode, http.StatusOK))
		nodes.POST("/:node_id/attach", formatDoc("Attach Node", "Group nodes"), tonic.Handler(r.attachNode, http.StatusCreated))
		nodes.DELETE("/:node_id/attach", formatDoc("Dettach Node", "Move node out of group"), tonic.Handler(r.detachNode, http.StatusOK))
		nodes.POST("/:node_id/sites", formatDoc("Add To Site", "Add node to site"), tonic.Handler(r.postNodeToSite, http.StatusCreated))
		nodes.DELETE("/:node_id/sites", formatDoc("Release From Site", "Release node from site"), tonic.Handler(r.deleteNodeFromSite, http.StatusOK))
	}
}

func (r *Router) postNetwork(c *gin.Context, req *AddNetworkReq) (*creg.NetworkInfo, error) {
	return r.network.CreateNetwork(req.OrgName, req.NetName, req.AllowedCountries,
		req.AllowedNetworks, req.Budget, req.Overdraft, req.TrafficPolicy, req.PaymentLinks)
}

func (r *Router) getNetwork(c *gin.Context, req *GetNetworkReq) (*creg.NetworkInfo, error) {
	return r.network.GetNetwork(req.NetworkId)
}

func (r *Router) postPackage(c *gin.Context, req *AddPackageReq) (*cdplan.PackageInfo, error) {
	return r.pkg.AddPackage(req.Name, req.OrgId, req.OwnerId, req.From, req.To, req.BaserateId, req.Active,
		req.Flatrate, req.SmsVolume, req.VoiceVolume, req.DataVolume, req.VoiceUnit, req.DataUnit, req.SimType,
		req.Apn, req.Type, req.Duration, req.Markup, req.Amount, req.Overdraft, req.TrafficPolicy, req.Networks)
}

func (r *Router) getPackage(c *gin.Context, req *GetPackageReq) (*cdplan.PackageInfo, error) {
	return r.pkg.GetPackage(req.PackageId)
}

func (r *Router) postSim(c *gin.Context, req *AddSimReq) (*csub.SimInfo, error) {
	return r.sim.ConfigureSim(req.SubscriberId, req.OrgId, req.NetworkId, req.Name,
		req.Email, req.PhoneNumber, req.Address, req.Dob, req.ProofOfIdentification,
		req.IdSerial, req.PackageId, req.SimType, req.SimToken, req.TrafficPolicy)
}

func (r *Router) getSim(c *gin.Context, req *GetSimReq) (*csub.SimInfo, error) {
	return r.sim.GetSim(req.Id)
}

func (r *Router) getNode(c *gin.Context, req *GetNodeRequest) (*creg.NodeInfo, error) {
	return r.node.GetNode(req.NodeId)
}

func (r *Router) postNode(c *gin.Context, req *AddNodeRequest) (*creg.NodeInfo, error) {
	return r.node.RegisterNode(req.NodeId, req.Name, req.OrgId, req.State)
}

func (r *Router) deleteNode(c *gin.Context, req *GetNodeRequest) error {
	return r.node.DeleteNode(req.NodeId)
}

func (r *Router) attachNode(c *gin.Context, req *AttachNodesRequest) error {
	return r.node.AttachNode(req.ParentNode, req.AmpNodeL, req.AmpNodeR)
}

func (r *Router) detachNode(c *gin.Context, req *GetNodeRequest) error {
	return r.node.DetachNode(req.NodeId)
}

func (r *Router) postNodeToSite(c *gin.Context, req *AddNodeToSiteRequest) error {
	return r.node.AddNodeToSite(req.NodeId, req.NetworkId, req.SiteId)
}

func (r *Router) deleteNodeFromSite(c *gin.Context, req *GetNodeRequest) error {
	return r.node.RemoveNodeFromSite(req.NodeId)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
