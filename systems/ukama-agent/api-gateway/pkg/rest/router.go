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
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg/client"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
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
	a asr
}

type asr interface {
	Activate(req *pb.ActivateReq) (*pb.ActivateResp, error)
	Inactivate(req *pb.InactivateReq) (*pb.InactivateResp, error)
	UpdatePackage(req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error)
	Read(req *pb.ReadReq) (*pb.ReadResp, error)
	GetUsage(req *pb.UsageReq) (*pb.UsageResp, error)
	GetUsageForPeriod(req *pb.UsageForPeriodReq) (*pb.UsageResp, error)
	QueryUsage(req *pb.QueryUsageReq) (*pb.QueryUsageResp, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.a = client.NewAsr(endpoints.Asr, endpoints.Timeout)
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
	auth := r.f.Group("/v1", "ukama-agent ", "Ukama-agent system", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
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
		asr := auth.Group("/asr", "Asr", "Active susbcriber registry")
		asr.GET("/:iccid", formatDoc("Get Subscriber", ""), tonic.Handler(r.getActiveSubscriber, http.StatusOK))
		asr.PUT("/:iccid", formatDoc("Activate: Add a new subscriber", ""), tonic.Handler(r.putSubscriber, http.StatusCreated))
		asr.PATCH("/:iccid", formatDoc("Update package id", ""), tonic.Handler(r.patchPackageUpdate, http.StatusOK))
		asr.DELETE("/:iccid", formatDoc("Inactivate: Remove a susbcriber", ""), tonic.Handler(r.deleteSubscriber, http.StatusOK))
		asr.GET("/:iccid/usage", formatDoc("Get Subscriber usage for time", ""), tonic.Handler(r.getUsage, http.StatusOK))
		asr.GET("/:iccid/period", formatDoc("Get Subscriber usage package", ""), tonic.Handler(r.getUsageForPeriod, http.StatusOK))

		usage := auth.Group("/usage", "Data Usage", "Active susbcriber data usage")
		usage.GET("/:iccid", formatDoc("Query data usage for a givent ICCID with various filtering params", ""), tonic.Handler(r.queryUsage, http.StatusOK))
	}

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) putSubscriber(c *gin.Context, req *ActivateReq) (*pb.ActivateResp, error) {
	log.Infof("Received a add subscriber request: %v", req)
	return r.clients.a.Activate(&pb.ActivateReq{
		Iccid:        req.Iccid,
		Imsi:         req.Imsi,
		SimPackageId: req.SimPackageId,
		PackageId:    req.PackageId,
		NetworkId:    req.NetworkId,
	})
}

func (r *Router) deleteSubscriber(c *gin.Context, req *DeactivateReq) (*pb.InactivateResp, error) {
	log.Infof("Received a update subscriber request: %v", req)
	return r.clients.a.Inactivate(&pb.InactivateReq{
		Iccid:     req.Iccid,
		Imsi:      req.Imsi,
		PackageId: req.PackageId,
		NetworkId: req.NetworkId,
	})
}

func (r *Router) patchPackageUpdate(c *gin.Context, req *UpdatePackageReq) (*pb.UpdatePackageResp, error) {
	log.Infof("Received a delete subscriber request: %v", req)
	return r.clients.a.UpdatePackage(&pb.UpdatePackageReq{
		Iccid:        req.Iccid,
		Imsi:         req.Imsi,
		SimPackageId: req.SimPackageId,
		PackageId:    req.PackageId,
		NetworkId:    req.NetworkId,
	})
}

func (r *Router) getActiveSubscriber(c *gin.Context, req *ReadSubscriberReq) (*pb.ReadResp, error) {
	return r.clients.a.Read(&pb.ReadReq{
		Id: &pb.ReadReq_Iccid{
			Iccid: req.Iccid,
		},
	})
}

func (r *Router) getUsage(c *gin.Context, req *UsageRequest) (*pb.UsageResp, error) {
	return r.clients.a.GetUsage(&pb.UsageReq{
		Id: &pb.UsageReq_Iccid{
			Iccid: req.Iccid,
		},
	})
}

func (r *Router) getUsageForPeriod(c *gin.Context, req *UsageForPeriodRequest) (*pb.UsageResp, error) {
	return r.clients.a.GetUsageForPeriod(&pb.UsageForPeriodReq{
		Id: &pb.UsageForPeriodReq_Iccid{
			Iccid: req.Iccid,
		},
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
}

func (r *Router) queryUsage(c *gin.Context, req *QueryUsageRequest) (*pb.QueryUsageResp, error) {
	return r.clients.a.QueryUsage(&pb.QueryUsageReq{
		Iccid:   req.Iccid,
		NodeId:  req.NodeId,
		Session: req.Session,
		From:    req.From,
		To:      req.To,
		Count:   req.Count,
		Sort:    req.Sort,
	})
}
