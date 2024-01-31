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
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
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
		pcrf := auth.Group("/pcrf", "PCRF", "pcrf")

		s := pcrf.Group("/session", "session", "session")
		s.POST("", formatDoc("Create session", "Create a new session"), tonic.Handler(r.createSession, http.StatusAccepted))
		s.POST("", formatDoc("End session", "End a session"), tonic.Handler(r.endSession, http.StatusAccepted))
		s.GET("/:id", formatDoc("Get session", "Get a particular session"), tonic.Handler(r.getSessionByID, http.StatusOK))

		// cdr
		cdr := pcrf.Group("/cdr", "CDR", "cdr")
		cdr.GET("/:imsi", formatDoc("Get CDR", "Get CDR for imsi"), tonic.Handler(r.getCDRByImsi, http.StatusOK))
		cdr.GET("/:imsi/id/:id", formatDoc("Get a CDR", "Get a specific CDR for imsi"), tonic.Handler(r.getCDRById, http.StatusOK))

		// policy
		policy := pcrf.Group("/policy", "Policy", "policy")
		policy.POST("", formatDoc("Configure Policy", "Configure a new policy"), tonic.Handler(r.addPolicy, http.StatusCreated))
		policy.GET("/:imsi", formatDoc("Get Sim", "Get a specific sim"), tonic.Handler(r.getPolicy, http.StatusOK))
		policy.DELETE("", formatDoc("Delete Policy", "Delete a policy"), tonic.Handler(r.removePolicy, http.StatusOK))

		// re-route
		route := pcrf.Group("/reroute", "Reroute", "rerouting IP address")
		route.GET("", formatDoc("reroute address", "Get a rerouting IP Address"), tonic.Handler(r.getReroute, http.StatusOK))
		route.POST("", formatDoc("Add or Update Node", "Update a rerouting IP Address"), tonic.Handler(r.updateReroute, http.StatusCreated))

	}
}

func (r *Router) createSession(c *gin.Context, req *controller.CreateSession) error {
	return r.controller.CreateSession(req)
}

func (r *Router) endSession(c *gin.Context, req *controller.EndSession) error {
	return r.controller.EndSession(req)
}

func (r *Router) getSessionByID(c *gin.Context, req *controller.GetSessionByID) (*store.Session, error) {
	return r.controller.EndSession(req)
}

func (r *Router) getCDRById(c *gin.Context, req *controller.GetCDRById) (**controller.CDR, error) {
	return r.controller.GetCDRById(req)
}

func (r *Router) getCDRByImsi(c *gin.Context, req *controller.GetCDRByImsi) (**controller.CDR, error) {
	return r.controller.GetCDRByImsi(req)
}

func (r *Router) getPolicy(c *gin.Context, req *controller.PolicyByImsi) (*controller.Policy, error) {
	return r.controller.GetPolicy(req.Imsi)
}

func (r *Router) addPolicy(c *gin.Context, req *controller.AddPolicyByImsi) error {
	return r.controller.AddPolicy(req)
}

func (r *Router) removePolicy(c *gin.Context, req *controller.PolicyByImsi) error {
	return r.controller.RemovePolicy(req)
}

func (r *Router) getReroute(c *gin.Context) (*controller.Reroute, error) {
	return r.controller.GetPolicy()
}

func (r *Router) updateReroute(c *gin.Context, req *controller.Reroute) error {
	return r.controller.UpdateReroute(req)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
