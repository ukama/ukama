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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	log "github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/ukama/ukama/systems/operation/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/operation/api-gateway/pkg"
	"github.com/ukama/ukama/systems/operation/api-gateway/pkg/client"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
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
	grpcEndpoints *pkg.GrpcEndpoints
	descriptions  *pkg.ServiceDescriptions
}

type Clients struct {
	Manager manager
}

type manager interface {
	Start(*pb.StartOperationRequest) (*pb.StartOperationResponse, error)
	Get(id string) (*pb.GetOperationResponse, error)
	GetByResource(resourceKey string) (*pb.GetByResourceResponse, error)
	MarkRunning(id string, fencingToken uint64, actor string) (*pb.MarkRunningResponse,
		error)
	ForceUnlock(id, actor, reason string) (*pb.ForceUnlockResponse, error)
	Complete(id string, fencingToken uint64, actor,
		reason string) (*pb.CompleteOperationResponse, error)
	Fail(id string, fencingToken uint64, actor,
		reason string) (*pb.FailOperationResponse, error)
	Delete(id, actor string) (*pb.DeleteOperationResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	return &Clients{
		Manager: client.NewManager(endpoints.Manager, endpoints.Timeout),
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
		httpEndpoints: &svcConf.Http,
		serverConf:    &svcConf.Server,
		grpcEndpoints: &svcConf.Services,
		descriptions:  &svcConf.Descriptions,
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

	desc := r.config.descriptions
	if desc == nil {
		desc = &pkg.ServiceDescriptions{}
	}

	if r.config.grpcEndpoints != nil {
		rest.RegisterStatusEndpoint(r.f, pkg.SystemName, map[string]rest.StatusTarget{
			"manager": {Host: r.config.grpcEndpoints.Manager, Description: desc.Manager},
		}, r.config.grpcEndpoints.Timeout)
	}

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
		ops.POST("", formatDoc("Start operations", "Locks every resource_keys entry and creates one pending operation per key, in one transaction, provided no must_be_free key is held by an active operation. Returns 409 if any key is busy."), tonic.Handler(r.postStartHandler, http.StatusCreated))
		ops.GET("", formatDoc("Get operation by resource", "Returns explicit lock state for resource_key. When free: locked=false and operation=null. When locked: locked=true and operation is populated."), tonic.Handler(r.getByResourceHandler, http.StatusOK))
		ops.GET("/:id", formatDoc("Get operation by id", "Returns the current state of an operation."), tonic.Handler(r.getOperationHandler, http.StatusOK))
		ops.PATCH("/:id", formatDoc("Change operation status",
			"Moves an operation to RUNNING, SUCCESS or FAILED. Caller must pass the "+
				"fencing token; SUCCESS and FAILED also need a reason. Returns 409 if "+
				"the token, lease or current status does not allow it."),
			tonic.Handler(r.patchOperationHandler, http.StatusOK))
		ops.DELETE("/:id/lock", formatDoc("Release an operation lock",
			"Privileged. Cancels the operation and releases its lock without a "+
				"fencing token; actor and reason are recorded in the audit."),
			tonic.Handler(r.deleteLockHandler, http.StatusOK))
		ops.DELETE("/:id", formatDoc("Delete an operation",
			"Soft-deletes a finished operation. Returns 409 while the operation is "+
				"pending or running."),
			tonic.Handler(r.deleteOperationHandler, http.StatusOK))
	}
}

func (r *Router) postStartHandler(c *gin.Context, req *StartOperationRequest) (*StartOperationResponse, error) {
	resp, err := r.clients.Manager.Start(&pb.StartOperationRequest{
		Type:           req.Type,
		System:         req.System,
		ResourceKeys:   req.ResourceKeys,
		MustBeFree:     req.MustBeFree,
		RequestedBy:    req.RequestedBy,
		IdempotencyKey: req.IdempotencyKey,
		LeaseSeconds:   req.LeaseSeconds,
	})
	if err != nil {
		return nil, err
	}

	ops := make([]*Operation, 0, len(resp.Operations))
	for _, op := range resp.Operations {
		ops = append(ops, operationFromProto(op))
	}
	return &StartOperationResponse{Operations: ops}, nil
}

func (r *Router) getOperationHandler(c *gin.Context, req *GetOperationRequest) (*GetOperationResponse, error) {
	resp, err := r.clients.Manager.Get(req.Id)
	if err != nil {
		return nil, err
	}

	return &GetOperationResponse{
		Operation: operationFromProto(resp.Operation),
	}, nil
}

func (r *Router) getByResourceHandler(c *gin.Context, req *GetByResourceRequest) (*GetByResourceResponse, error) {
	resp, err := r.clients.Manager.GetByResource(req.ResourceKey)
	if err != nil {
		return nil, err
	}

	op := operationFromProto(resp.Operation)

	return &GetByResourceResponse{
		Locked:    op != nil,
		Operation: op,
	}, nil
}

func (r *Router) patchOperationHandler(c *gin.Context,
	req *UpdateOperationRequest) (*UpdateOperationResponse, error) {
	var (
		op  *pb.Operation
		err error
	)
	switch req.Status {
	case pb.OperationStatus_RUNNING.String():
		var resp *pb.MarkRunningResponse
		if resp, err = r.clients.Manager.MarkRunning(req.Id, req.FencingToken,
			req.Actor); err == nil {
			op = resp.Operation
		}
	case pb.OperationStatus_SUCCESS.String():
		var resp *pb.CompleteOperationResponse
		if resp, err = r.clients.Manager.Complete(req.Id, req.FencingToken, req.Actor,
			req.Reason); err == nil {
			op = resp.Operation
		}
	case pb.OperationStatus_FAILED.String():
		var resp *pb.FailOperationResponse
		if resp, err = r.clients.Manager.Fail(req.Id, req.FencingToken, req.Actor,
			req.Reason); err == nil {
			op = resp.Operation
		}
	default:
		return nil, rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: fmt.Sprintf("status %q cannot be set", req.Status)}
	}
	if err != nil {
		return nil, restErr(err)
	}

	return &UpdateOperationResponse{Operation: operationFromProto(op)}, nil
}

func (r *Router) deleteLockHandler(c *gin.Context,
	req *ReleaseLockRequest) (*UpdateOperationResponse, error) {
	resp, err := r.clients.Manager.ForceUnlock(req.Id, req.Actor, req.Reason)
	if err != nil {
		return nil, restErr(err)
	}

	return &UpdateOperationResponse{Operation: operationFromProto(resp.Operation)}, nil
}

func (r *Router) deleteOperationHandler(c *gin.Context,
	req *DeleteOperationRequest) error {
	if _, err := r.clients.Manager.Delete(req.Id, req.Actor); err != nil {
		return restErr(err)
	}
	return nil
}

func restErr(err error) error {
	if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
		return rest.HttpError{HttpCode: http.StatusConflict, Message: st.Message()}
	}
	return err
}

func operationFromProto(op *pb.Operation) *Operation {
	if op == nil {
		return nil
	}

	return &Operation{
		Id:             op.Id,
		GroupId:        op.GroupId,
		Type:           op.Type,
		System:         op.System,
		Status:         op.Status.String(),
		FencingToken:   op.FencingToken,
		RequestedBy:    op.RequestedBy,
		IdempotencyKey: op.IdempotencyKey,
		ResourceKey:    op.ResourceKey,
		LeaseExpiresAt: timestampAsTime(op.LeaseExpiresAt),
		Error:          op.Error,
		StartedAt:      timestampAsTime(op.StartedAt),
		TerminalAt:     timestampAsTime(op.TerminalAt),
		CreatedAt:      timestampAsTime(op.CreatedAt),
	}
}

func timestampAsTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil || !ts.IsValid() {
		return nil
	}

	t := ts.AsTime()
	return &t
}

func formatDoc(summary, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
