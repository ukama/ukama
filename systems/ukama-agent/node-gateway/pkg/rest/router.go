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
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/pkg/client"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	cpb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
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
	c cdr
}

type asr interface {
	UpdateGuti(req *pb.UpdateGutiReq) (*pb.UpdateGutiResp, error)
	UpdateTai(req *pb.UpdateTaiReq) (*pb.UpdateTaiResp, error)
	Read(req *pb.ReadReq) (*pb.ReadResp, error)
}

type cdr interface {
	PostCDR(req *cpb.CDR) (*cpb.CDRResp, error)
	GetCDR(req *cpb.RecordReq) (*cpb.RecordResp, error)
	GetUsage(req *cpb.UsageReq) (*cpb.UsageResp, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.a = client.NewAsr(endpoints.Asr, endpoints.Timeout)
	c.c = client.NewCDR(endpoints.CDR, endpoints.Timeout)
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
	auth := r.f.Group("/v1", "ukama-agent-node-gateway ", "Ukama-agent system", func(ctx *gin.Context) {
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
		asr.GET("/:imsi", formatDoc("Get Subscriber", ""), tonic.Handler(r.getActiveSubscriber, http.StatusOK))
		asr.POST("/:imsi/guti", formatDoc("GUTI update for subscriber", ""), tonic.Handler(r.postGuti, http.StatusOK))
		asr.POST("/:imsi/tai", formatDoc("TAI update for subscriber", ""), tonic.Handler(r.postTai, http.StatusOK))
		asr.GET("/:imsi/usage", formatDoc("Get total usage", ""), tonic.Handler(r.getUsage, http.StatusOK))

		cdr := auth.Group("/cdr", "CDR", "Call Detail Record")
		cdr.POST("/:imsi", formatDoc("Post CDR", ""), tonic.Handler(r.postCDR, http.StatusCreated))
		cdr.GET("/:imsi", formatDoc("Get CDR", ""), tonic.Handler(r.getCDR, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) postCDR(c *gin.Context, req *PostCDRReq) (*cpb.CDRResp, error) {

	return r.clients.c.PostCDR(&cpb.CDR{
		Session:       req.Session,
		Imsi:          req.Imsi,
		Policy:        req.Policy,
		ApnName:       req.ApnName,
		NodeId:        req.NodeId,
		Ip:            req.Ip,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		TxBytes:       req.TxBytes,
		RxBytes:       req.RxBytes,
		TotalBytes:    req.TotalBytes,
		LastUpdatedAt: req.LastUpdatedAt,
	})
}

func (r *Router) getCDR(c *gin.Context, req *GetCDRReq) (*cpb.RecordResp, error) {

	return r.clients.c.GetCDR(&cpb.RecordReq{
		Imsi:      req.Imsi,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Policy:    req.Policy,
		SessionId: req.SessionId,
	})

}

func (r *Router) getUsage(c *gin.Context, req *GetUsageReq) (*cpb.UsageResp, error) {

	return r.clients.c.GetUsage(&cpb.UsageReq{
		Imsi:      req.Imsi,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Policy:    req.Policy,
		SessionId: req.SessionId,
	})

}

func (r *Router) postGuti(c *gin.Context, req *UpdateGutiReq) (*pb.UpdateGutiResp, error) {

	return r.clients.a.UpdateGuti(&pb.UpdateGutiReq{
		Imsi:      req.Imsi,
		UpdatedAt: req.UpdatedAt,
		Guti: &pb.Guti{
			PlmnId: req.Guti.PlmnId,
			Mmegi:  req.Guti.Mmegi,
			Mmec:   req.Guti.Mmec,
			Mtmsi:  req.Guti.Mtmsi,
		},
	})

}

func (r *Router) postTai(c *gin.Context, req *UpdateTaiReq) (*pb.UpdateTaiResp, error) {

	return r.clients.a.UpdateTai(&pb.UpdateTaiReq{
		Imsi:      req.Imsi,
		UpdatedAt: req.UpdatedAt,
		Tac:       req.Tac,
	})
}

func (r *Router) getActiveSubscriber(c *gin.Context, req *GetSubscriberReq) (*pb.ReadResp, error) {
	return r.clients.a.Read(&pb.ReadReq{
		Id: &pb.ReadReq_Imsi{
			Imsi: req.Imsi,
		},
	})
}
