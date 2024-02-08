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

	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/cmd/version"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz/openapi"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
)

var REDIRECT_URI = "https://subscriber.dev.ukama.com/swagger/#/"

type Router struct {
	f          *fizz.Fizz
	controller *controller.Controller
	config     *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *crest.HttpConfig
	auth       *config.Auth
}

func NewRouter(ctr *controller.Controller, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
	r := &Router{
		controller: ctr,
		config:     config,
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

	auth := r.f.Group("/v1", "PCRF for Node", "API system version v1", func(ctx *gin.Context) {
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
		// pcrf routes
		pcrf := auth.Group("v1/pcrf", "PCRF", "pcrf")

		s := pcrf.Group("/session", "session", "session")
		s.POST("", formatDoc("Create session", "Create a new session"), tonic.Handler(r.createSession, http.StatusAccepted))
		s.POST("", formatDoc("End session", "End a session"), tonic.Handler(r.endSession, http.StatusAccepted))
		s.GET("/:id", formatDoc("Get session", "Get a particular session"), tonic.Handler(r.getSessionByID, http.StatusOK))

		// cdr
		cdr := pcrf.Group("/cdr", "CDR", "cdr")
		cdr.GET("/session/:id", formatDoc("Get a CDR", "Get CDR for sepcific session"), tonic.Handler(r.getCDRBySessionId, http.StatusOK))
		cdr.GET("/imsi/:imsi", formatDoc("Get CDR's", "Get CDR's for imsi"), tonic.Handler(r.getCDRByImsi, http.StatusOK))

		// policy
		policy := pcrf.Group("/policy", "Policy", "policy")
		policy.POST("", formatDoc("Configure Policy", "Configure a new policy"), tonic.Handler(r.addPolicy, http.StatusCreated))
		policy.GET("/imsi/:imsi", formatDoc("Get policy", "Get a policy for a specific sim"), tonic.Handler(r.getPolicyByImsi, http.StatusOK))
		policy.GET("/id/:id", formatDoc("Get policy", "Get a policy by specific ID"), tonic.Handler(r.getPolicyByID, http.StatusOK))

		// re-route
		route := pcrf.Group("/reroute", "Reroute", "rerouting IP address")
		route.GET("/imsi/:imsi", formatDoc("reroute address", "Get a rerouting IP Address for Imsi"), tonic.Handler(r.getRerouteForImsi, http.StatusOK))
		route.POST("", formatDoc("Add or Update Node", "Update a rerouting IP Address"), tonic.Handler(r.updateReroute, http.StatusCreated))

		// Subscriber
		sub := pcrf.Group("/subscriber", "Subscriber", "subscriber")
		sub.GET("/imsi/:imsi", formatDoc("Get subscriber", "Get a subscriber"), tonic.Handler(r.getSubscriber, http.StatusOK))
		sub.POST("/imsi/:imsi", formatDoc("Add or update subscriber", "Add or update subscriber"), tonic.Handler(r.updateSubscriber, http.StatusOK))
		sub.GET("/imsi/:imsi/flow", formatDoc("Get flows", "Get a subscriber UE data path"), tonic.Handler(r.getSubscriberFlows, http.StatusOK))

	}
}

func (r *Router) createSession(c *gin.Context, req *api.CreateSession) error {
	return r.controller.CreateSession(c, req)
}

func (r *Router) endSession(c *gin.Context, req *api.EndSession) error {
	return r.controller.EndSession(c, req)
}

func (r *Router) getSessionByID(c *gin.Context, req *api.GetSessionByID) (*api.SessionResponse, error) {
	return r.controller.GetSessionByID(c, req)
}

func (r *Router) getCDRBySessionId(c *gin.Context, req *api.GetCDRBySessionId) (*api.CDR, error) {
	return r.controller.GetCDRBySessionId(c, req)
}

func (r *Router) getCDRByImsi(c *gin.Context, req *api.GetCDRByImsi) ([]*api.CDR, error) {
	return r.controller.GetCDRByImsi(c, req)
}

func (r *Router) getPolicyByImsi(c *gin.Context, req *api.GetPolicyByImsi) (*api.PolicyResponse, error) {
	return r.controller.GetPolicyByImsi(c, req)
}

func (r *Router) getPolicyByID(c *gin.Context, req *api.GetPolicyByID) (*api.PolicyResponse, error) {
	return r.controller.GetPolicyByID(c, req)
}

func (r *Router) addPolicy(c *gin.Context, req *api.AddPolicyByImsi) error {
	return r.controller.AddPolicy(c, req)
}

func (r *Router) getRerouteForImsi(c *gin.Context, req *api.GetReRouteByImsi) (*api.ReRouteResponse, error) {
	return r.controller.GetReroute(c, req)
}

func (r *Router) updateReroute(c *gin.Context, req *api.UpdateRerouteById) error {
	return r.controller.UpdateReroute(c, req)
}

func (r *Router) updateSubscriber(c *gin.Context, req *api.CreateSubscriber) error {
	return r.controller.AddSubscriber(c, req)
}

func (r *Router) getSubscriber(c *gin.Context, req *api.RequestSubscriber) (*api.SubscriberResponse, error) {
	return r.controller.GetSubscriber(c, req)
}

func (r *Router) getSubscriberFlows(c *gin.Context, req *api.GetFlowsForImsi) ([]*api.FlowResponse, error) {
	return r.controller.GetFlowsForImsi(c, req)
}
func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
