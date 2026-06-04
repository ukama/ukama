/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	log "github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"

	"github.com/ukama/ukama/systems/operation/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/operation/api-gateway/pkg"
	"github.com/ukama/ukama/systems/operation/api-gateway/pkg/client"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

// roles permitted to force-unlock an operation. Matches the role-gating pattern
// used by notification/distributor and event-notify (member lookup, not Keto).
var forceUnlockRoles = map[string]bool{
	"ROLE_OWNER": true,
	"ROLE_ADMIN": true,
}

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
	Manager manager
	Member  member
}

type manager interface {
	Start(*pb.StartOperationRequest) (*pb.StartOperationResponse, error)
	Get(id string) (*pb.GetOperationResponse, error)
	GetByResource(resourceKey string) (*pb.GetByResourceResponse, error)
	MarkRunning(id string, fencingToken uint64) (*pb.MarkRunningResponse, error)
	ForceUnlock(id, actor, reason string) (*pb.ForceUnlockResponse, error)
}

type member interface {
	GetByUserId(id string) (*creg.MemberInfoResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints, registryHost string) *Clients {
	return &Clients{
		Manager: client.NewManager(endpoints.Manager, endpoints.Timeout),
		Member:  creg.NewMemberClient(registryHost),
	}
}

func NewRouter(clients *Clients, conf *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
	r := &Router{clients: clients, config: conf}
	if !conf.debugMode {
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
	if err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port)); err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	auth := r.f.Group("/v1", "API gateway", "operation system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}
		meta := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", meta)
		if err := f(ctx, r.config.auth.AuthAPIGW); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})
	auth.Use()
	{
		ops := auth.Group("/operations", "Operations", "Lock + operation lifecycle")
		ops.POST("", formatDoc("Start an operation", "Acquires a lock and creates an operation in pending state. Returns 409 if the resource is already locked."), tonic.Handler(r.postStartHandler, http.StatusCreated))
		ops.GET("", formatDoc("Get operation by resource", "Returns the operation currently locking the given resource_key, or empty if free."), tonic.Handler(r.getByResourceHandler, http.StatusOK))
		ops.GET("/:id", formatDoc("Get operation by id", "Returns the current state of an operation."), tonic.Handler(r.getOperationHandler, http.StatusOK))
		ops.POST("/:id/run", formatDoc("Mark operation running", "Transitions a pending operation to running. Caller must pass the fencing token."), tonic.Handler(r.postMarkRunningHandler, http.StatusOK))
		ops.POST("/:id/force-unlock", formatDoc("Force-unlock an operation", "Privileged. Cancels the operation and releases its lock with audit reason."), tonic.Handler(r.postForceUnlockHandler, http.StatusOK))
	}
}

func (r *Router) postStartHandler(c *gin.Context, req *StartOperationRequest) (*pb.StartOperationResponse, error) {
	return r.clients.Manager.Start(&pb.StartOperationRequest{
		Type:           req.Type,
		System:         req.System,
		ResourceKey:    req.ResourceKey,
		RequestedBy:    req.RequestedBy,
		IdempotencyKey: req.IdempotencyKey,
		LeaseSeconds:   req.LeaseSeconds,
	})
}

func (r *Router) getOperationHandler(c *gin.Context, req *GetOperationRequest) (*pb.GetOperationResponse, error) {
	return r.clients.Manager.Get(req.Id)
}

func (r *Router) getByResourceHandler(c *gin.Context, req *GetByResourceRequest) (*pb.GetByResourceResponse, error) {
	return r.clients.Manager.GetByResource(req.ResourceKey)
}

func (r *Router) postMarkRunningHandler(c *gin.Context, req *MarkRunningRequest) (*pb.MarkRunningResponse, error) {
	return r.clients.Manager.MarkRunning(req.Id, req.FencingToken)
}

func (r *Router) postForceUnlockHandler(c *gin.Context, req *ForceUnlockRequest) (*pb.ForceUnlockResponse, error) {
	resp, err := r.clients.Member.GetByUserId(req.UserId)
	if err != nil {
		log.Errorf("ForceUnlock: failed to resolve member for user %s: %v", req.UserId, err)
		return nil, rest.HttpError{HttpCode: http.StatusForbidden, Message: "could not verify caller role"}
	}
	if !forceUnlockRoles[resp.Member.Role] {
		log.Warnf("ForceUnlock denied: user %s has role %s (owner/admin required)", req.UserId, resp.Member.Role)
		return nil, rest.HttpError{HttpCode: http.StatusForbidden, Message: "only org owner or admin may force-unlock"}
	}
	return r.clients.Manager.ForceUnlock(req.Id, req.UserId, req.Reason)
}

func formatDoc(summary, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
