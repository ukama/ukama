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
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/init/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg/client"

	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
)

const ORG_URL_PARAMETER = "org"

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
	l lookup
}

type lookup interface {
	AddOrg(req *pb.AddOrgRequest) (*pb.AddOrgResponse, error)
	UpdateOrg(req *pb.UpdateOrgRequest) (*pb.UpdateOrgResponse, error)
	GetOrg(req *pb.GetOrgRequest) (*pb.GetOrgResponse, error)
	GetOrgs(req *pb.GetOrgsRequest) (*pb.GetOrgsResponse, error)
	AddNodeForOrg(req *pb.AddNodeRequest) (*pb.AddNodeResponse, error)
	GetNodeForOrg(req *pb.GetNodeForOrgRequest) (*pb.GetNodeResponse, error)
	DeleteNodeForOrg(req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error)
	AddSystemForOrg(req *pb.AddSystemRequest) (*pb.AddSystemResponse, error)
	UpdateSystemForOrg(req *pb.UpdateSystemRequest) (*pb.UpdateSystemResponse, error)
	GetSystemForOrg(req *pb.GetSystemRequest) (*pb.GetSystemResponse, error)
	DeleteSystemForOrg(req *pb.DeleteSystemRequest) (*pb.DeleteSystemResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.l = client.Newlookup(endpoints.Lookup, endpoints.Timeout)
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
	logrus.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "API gateway", "Init system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			logrus.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})
	auth.Use()
	{
		o := auth.Group("/orgs", "Orgs", "looking for orgs credentials")
		o.GET("", formatDoc("Get Orgs name", ""), tonic.Handler(r.getOrgsHandler, http.StatusOK))

		const org = "/orgs/" + ":" + ORG_URL_PARAMETER
		orgs := auth.Group(org, "Orgs", "looking for orgs credentials")
		orgs.GET("", formatDoc("Get Org by name", ""), tonic.Handler(r.getOrgHandler, http.StatusOK))
		orgs.PUT("", formatDoc("Add Org and Credential", ""), tonic.Handler(r.putOrgHandler, http.StatusCreated))
		orgs.PATCH("", formatDoc("Update Orgs Credential", ""), tonic.Handler(r.patchOrgHandler, http.StatusOK))

		nodes := orgs.Group("/nodes", "Nodes", "Orgs credentials for Node")
		nodes.GET("/:node", formatDoc("Get Orgs credential for Node", ""), tonic.Handler(r.getNodeHandler, http.StatusOK))
		nodes.PUT("/:node", formatDoc("Add Node to Org", ""), tonic.Handler(r.putNodeHandler, http.StatusCreated))
		nodes.DELETE("/:node", formatDoc("Delete Node from Org", ""), tonic.Handler(r.deleteNodeHandler, http.StatusOK))

		systems := orgs.Group("/systems", "Systems", "Orgs System credentials")
		systems.GET("/:system", formatDoc("Get System credential for Org", ""), tonic.Handler(r.getSystemHandler, http.StatusOK))
		systems.PUT("/:system", formatDoc("Add or Update System credential for Org", ""), tonic.Handler(r.putSystemHandler, http.StatusCreated))
		systems.DELETE("/:system", formatDoc("Delete System credential for Org", ""), tonic.Handler(r.deleteSystemHandler, http.StatusOK))
		systems.PATCH("/:system", formatDoc("Update System Credential", ""), tonic.Handler(r.patchSystemHandler, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getOrgsHandler(c *gin.Context, req *GetOrgsRequest) (*pb.GetOrgsResponse, error) {
	return r.clients.l.GetOrgs(&pb.GetOrgsRequest{})
}

func (r *Router) getOrgHandler(c *gin.Context, req *GetOrgRequest) (*pb.GetOrgResponse, error) {
	return r.clients.l.GetOrg(&pb.GetOrgRequest{
		OrgName: req.OrgName,
	})
}

func (r *Router) putOrgHandler(c *gin.Context, req *AddOrgRequest) (*pb.AddOrgResponse, error) {
	return r.clients.l.AddOrg(&pb.AddOrgRequest{
		OrgName:     req.OrgName,
		OrgId:       req.OrgId,
		Certificate: req.Certificate,
		Ip:          req.Ip,
	})

}

func (r *Router) patchOrgHandler(c *gin.Context, req *UpdateOrgRequest) (*pb.UpdateOrgResponse, error) {

	return r.clients.l.UpdateOrg(&pb.UpdateOrgRequest{
		OrgName:     req.OrgName,
		Certificate: req.Certificate,
		Ip:          req.Ip,
	})
}

func (r *Router) putNodeHandler(c *gin.Context, req *AddNodeRequest) (*pb.AddNodeResponse, error) {
	return r.clients.l.AddNodeForOrg(&pb.AddNodeRequest{
		OrgName: req.OrgName,
		NodeId:  req.NodeId,
	})
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*pb.GetNodeResponse, error) {
	return r.clients.l.GetNodeForOrg(&pb.GetNodeForOrgRequest{
		OrgName: req.OrgName,
		NodeId:  req.NodeId,
	})
}

func (r *Router) deleteNodeHandler(c *gin.Context, req *DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {

	return r.clients.l.DeleteNodeForOrg(&pb.DeleteNodeRequest{
		OrgName: req.OrgName,
		NodeId:  req.NodeId,
	})
}

func (r *Router) putSystemHandler(c *gin.Context, req *AddSystemRequest) (*pb.AddSystemResponse, error) {
	return r.clients.l.AddSystemForOrg(&pb.AddSystemRequest{
		OrgName:     req.OrgName,
		SystemName:  req.SysName,
		Certificate: req.Certificate,
		Ip:          req.Ip,
		Port:        req.Port,
		Url:         req.URL,
		NodeGwIp:    req.NodeGwIp,
		NodeGwPort:  req.NodeGwPort,
		NodeGwUrl:   req.NodeGwURL,
	})

}

func (r *Router) patchSystemHandler(c *gin.Context, req *UpdateSystemRequest) (*pb.UpdateSystemResponse, error) {

	return r.clients.l.UpdateSystemForOrg(&pb.UpdateSystemRequest{
		OrgName:     req.OrgName,
		SystemName:  req.SysName,
		Certificate: req.Certificate,
		Ip:          req.Ip,
		Port:        req.Port,
		NodeGwIp:    req.NodeGwIp,
		NodeGwPort:  req.NodeGwPort,
		NodeGwUrl:   req.NodeGwURL,
	})
}

func (r *Router) getSystemHandler(c *gin.Context, req *GetSystemRequest) (*pb.GetSystemResponse, error) {
	return r.clients.l.GetSystemForOrg(&pb.GetSystemRequest{
		OrgName:    req.OrgName,
		SystemName: req.SysName,
	})
}

func (r *Router) deleteSystemHandler(c *gin.Context, req *DeleteSystemRequest) (*pb.DeleteSystemResponse, error) {

	return r.clients.l.DeleteSystemForOrg(&pb.DeleteSystemRequest{
		OrgName:    req.OrgName,
		SystemName: req.SysName,
	})
}
