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
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/messaging/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg/client"

	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
)

const ORG_URL_PARAMETER = "org"

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig      config.Metrics
	httpEndpoints      *pkg.HttpEndpoints
	debugMode          bool
	serverConf         *rest.HttpConfig
	auth               *config.Auth
	orgName            string
	nodeTargetTemplate string
	grpcEndpoints      *pkg.GrpcEndpoints
	descriptions       *pkg.ServiceDescriptions
}

type Clients struct {
	n nns
}

type nns interface {
	GetNodeRequest(req *pb.GetNodeRequest) (*pb.GetNodeResponse, error)
	GetMeshRequest(req *pb.GetMeshRequest) (*pb.GetMeshResponse, error)
	SetRequest(req *pb.SetRequest) (*pb.SetResponse, error)
	UpdateMeshRequest(req *pb.UpdateMeshRequest) (*pb.UpdateMeshResponse, error)
	UpdateNodeRequest(req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error)
	DeleteRequest(req *pb.DeleteRequest) (*pb.DeleteResponse, error)
	ListRequest(req *pb.ListRequest) (*pb.ListResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.n = client.NewNns(endpoints.Nns, endpoints.Timeout)
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
		httpEndpoints: &svcConf.Http,
		serverConf:    &svcConf.Server,
		grpcEndpoints: &svcConf.Services,
		descriptions:  &svcConf.Descriptions,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,

		orgName:            svcConf.OrgName,
		nodeTargetTemplate: svcConf.NodeTargetTemplate,
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

	desc := r.config.descriptions
	if desc == nil {
		desc = &pkg.ServiceDescriptions{}
	}

	if r.config.grpcEndpoints != nil {
		rest.RegisterStatusEndpoint(r.f, pkg.SystemName, map[string]rest.StatusTarget{
			"nns":         {Host: r.config.grpcEndpoints.Nns, Description: desc.Nns},
			"broadcaster": {Host: r.config.grpcEndpoints.Broadcaster, Description: desc.Broadcaster},
			"node-feeder": {Host: r.config.grpcEndpoints.NodeFeeder, Description: desc.NodeFeeder},
		}, r.config.grpcEndpoints.Timeout)
	}

	auth := r.f.Group("/v1", "Messaging system API Gateway", "Messaging system version v1", func(ctx *gin.Context) {
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
		nns := auth.Group("/nns", "Nns", "Looking for node Ip address")
		nns.GET("/node/:node_id", formatDoc("Get node Ip", ""), tonic.Handler(r.getNodeHandler, http.StatusOK))
		nns.GET("/mesh/:node_id", formatDoc("Get mesh ip", ""), tonic.Handler(r.getMeshHandler, http.StatusOK))
		nns.PUT("/node", formatDoc("Add node", ""), tonic.Handler(r.putNodeHandler, http.StatusCreated))
		nns.PUT("/node/:node_id", formatDoc("Update node", ""), tonic.Handler(r.updateNodeHandler, http.StatusCreated))
		nns.PUT("/mesh/:node_id", formatDoc("Update mesh", ""), tonic.Handler(r.updateMeshHandler, http.StatusOK))
		nns.DELETE("/node/:node_id", formatDoc("Remove node from dns", ""), tonic.Handler(r.deleteHandler, http.StatusOK))
		nns.GET("/list", formatDoc("Get all nodes", ""), tonic.Handler(r.listHandler, http.StatusOK))

		prom := auth.Group("/prometheus", "Prometheus target", "Target discovery endpoint")
		prom.GET("", formatDoc("Get targets to scrape",
			"Prometheus http_sd_config endpoint. Returns one target per known node, "+
				"labelled with node_id, org, network and site."),
			tonic.Handler(r.prometheusHandler, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*pb.GetNodeResponse, error) {

	return r.clients.n.GetNodeRequest(&pb.GetNodeRequest{
		NodeId: req.NodeId,
	})
}

func (r *Router) getMeshHandler(c *gin.Context, req *GetMeshRequest) (*pb.GetMeshResponse, error) {
	return r.clients.n.GetMeshRequest(&pb.GetMeshRequest{
		NodeId: req.NodeId,
	})
}

func (r *Router) putNodeHandler(c *gin.Context, req *SetNodeRequest) (*pb.SetResponse, error) {
	return r.clients.n.SetRequest(&pb.SetRequest{
		NodeId:       req.NodeId,
		NodeIp:       req.NodeIp,
		NodePort:     req.NodePort,
		MeshIp:       req.MeshIp,
		MeshPort:     req.MeshPort,
		Network:      req.Network,
		Site:         req.Site,
		MeshHostName: req.MeshHostName,
	})

}

func (r *Router) updateNodeHandler(c *gin.Context, req *UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	return r.clients.n.UpdateNodeRequest(&pb.UpdateNodeRequest{
		NodeId:   req.NodeId,
		NodeIp:   req.NodeIp,
		NodePort: req.NodePort,
	})
}

func (r *Router) updateMeshHandler(c *gin.Context, req *UpdateMeshRequest) (*pb.UpdateMeshResponse, error) {
	return r.clients.n.UpdateMeshRequest(&pb.UpdateMeshRequest{
		NodeId:   req.NodeId,
		MeshIp:   req.MeshIp,
		MeshPort: req.MeshPort,
	})
}

func (r *Router) deleteHandler(c *gin.Context, req *DeleteRequest) (*pb.DeleteResponse, error) {
	return r.clients.n.DeleteRequest(&pb.DeleteRequest{
		NodeId: req.NodeId,
	})
}

func (r *Router) listHandler(c *gin.Context, req *ListRequest) (*pb.ListResponse, error) {
	return r.clients.n.ListRequest(&pb.ListRequest{
		NodeId: req.NodeId,
	})
}

func (r *Router) prometheusHandler(c *gin.Context, req *PrometheusTargetsRequest) ([]PrometheusTarget, error) {
	resp, err := r.clients.n.ListRequest(&pb.ListRequest{
		NodeId: req.NodeId,
	})
	if err != nil {
		return nil, err
	}

	return buildPrometheusTargets(resp.GetList(), r.config.orgName, r.config.nodeTargetTemplate)
}

var nodeTargetFuncs = template.FuncMap{"lower": strings.ToLower}

// buildPrometheusTargets turns the nns mesh mapping into http_sd_config
// entries. A node whose address cannot be rendered is skipped rather than
// failing the whole response, so one bad record cannot blind Prometheus to the
// rest of the fleet. Empty network/site labels are omitted: a node that is not
// attached to a site should carry no site label at all.
func buildPrometheusTargets(nodes []*pb.NodeMeshInfo, orgName, tmplText string) ([]PrometheusTarget, error) {
	if strings.TrimSpace(tmplText) == "" {
		tmplText = pkg.DefaultNodeTargetTemplate
	}

	tmpl, err := template.New("nodeTarget").Funcs(nodeTargetFuncs).Parse(tmplText)
	if err != nil {
		return nil, fmt.Errorf("invalid node target template: %w", err)
	}

	targets := make([]PrometheusTarget, 0, len(nodes))

	for _, node := range nodes {
		if node.GetNodeId() == "" {
			continue
		}

		org := node.GetOrg()
		if org == "" {
			org = orgName
		}

		address := &bytes.Buffer{}

		err := tmpl.Execute(address, pkg.NodeTargetVars{
			OrgName:  org,
			NodeId:   node.GetNodeId(),
			NodeIp:   node.GetNodeIp(),
			NodePort: node.GetNodePort(),
			MeshIp:   node.GetMeshIp(),
			MeshPort: node.GetMeshPort(),
			Network:  node.GetNetwork(),
			Site:     node.GetSite(),
		})
		if err != nil {
			logrus.Errorf("Skipping node %s: failed to render scrape target. Error: %v", node.GetNodeId(), err)

			continue
		}

		if address.Len() == 0 {
			logrus.Errorf("Skipping node %s: rendered scrape target is empty", node.GetNodeId())

			continue
		}

		labels := map[string]string{"node_id": node.GetNodeId()}

		putLabel(labels, "org", org)
		putLabel(labels, "network", node.GetNetwork())
		putLabel(labels, "site", node.GetSite())

		targets = append(targets, PrometheusTarget{
			Targets: []string{address.String()},
			Labels:  labels,
		})
	}

	return targets, nil
}

func putLabel(labels map[string]string, key, value string) {
	if value == "" {
		return
	}

	labels[key] = value
}
