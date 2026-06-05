/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

type OperationServer struct {
	pb.UnimplementedOperationManagerServiceServer
	orgName    string
	orgId      string
	repo       db.OperationRepo
	msgbus     mb.MsgBusServiceClient
	routingKey msgbus.RoutingKeyBuilder
}

func NewOperationServer(orgName, orgId string, repo db.OperationRepo, msgBus mb.MsgBusServiceClient) *OperationServer {
	return &OperationServer{
		orgName: orgName,
		orgId:   orgId,
		repo:    repo,
		msgbus:  msgBus,
		routingKey: msgbus.NewRoutingKeyBuilder().
			SetCloudSource().
			SetGlobalScope().
			SetSystem(pkg.SystemName).
			SetOrgName(orgName).
			SetService(pkg.ServiceName),
	}
}

func (s *OperationServer) StartOperation(ctx context.Context, req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	if req.IdempotencyKey != "" {
		existing, err := s.repo.GetByIdempotencyKey(req.IdempotencyKey)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "idempotency lookup: %v", err)
		}
		if existing != nil {
			return &pb.StartOperationResponse{Operation: toPb(existing)}, nil
		}
	}

	lease := time.Duration(req.LeaseSeconds) * time.Second
	if lease == 0 {
		lease = pkg.DefaultLeaseTTL
	}

	op := &db.Operation{
		Id:             uuid.NewV4(),
		Type:           req.Type,
		System:         req.System,
		Status:         db.OperationPending,
		RequestedBy:    req.RequestedBy,
		ResourceKey:    req.ResourceKey,
		LeaseExpiresAt: time.Now().UTC().Add(lease),
	}
	if req.IdempotencyKey != "" {
		k := req.IdempotencyKey
		op.IdempotencyKey = &k
	}

	op, err := s.repo.Start(op, lease)
	if errors.Is(err, db.ErrLockConflict) {
		resp := &pb.StartOperationResponse{}
		if op != nil {
			resp.ConflictingOperation = toPb(op)
		}
		return resp, status.Error(codes.AlreadyExists, "resource is locked by an active operation")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "start operation: %v", err)
	}

	log.Infof("operation %s started: type=%s resource=%s token=%d", op.Id, op.Type, op.ResourceKey, op.FencingToken)
	return &pb.StartOperationResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) GetOperation(ctx context.Context, req *pb.GetOperationRequest) (*pb.GetOperationResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	op, err := s.repo.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get operation: %v", err)
	}
	return &pb.GetOperationResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) GetByResource(ctx context.Context, req *pb.GetByResourceRequest) (*pb.GetByResourceResponse, error) {
	op, err := s.repo.GetByResource(req.ResourceKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get by resource: %v", err)
	}
	if op == nil {
		return &pb.GetByResourceResponse{}, nil
	}
	return &pb.GetByResourceResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) MarkRunning(ctx context.Context, req *pb.MarkRunningRequest) (*pb.MarkRunningResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	op, err := s.repo.MarkRunning(id, req.FencingToken)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}
	return &pb.MarkRunningResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) CompleteOperation(ctx context.Context, req *pb.ForceUnlockRequest) (*pb.ForceUnlockResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	current, err := s.repo.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get operation: %v", err)
	}

	op, err := s.repo.Terminate(id, current.FencingToken, db.OperationSuccess, db.OperationAudit{
		Event:  "completed",
		Actor:  req.Actor,
		Reason: req.Reason,
	}, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "complete operation: %v", err)
	}

	log.Infof("operation %s completed by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.ForceUnlockResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) FailOperation(ctx context.Context, req *pb.ForceUnlockRequest) (*pb.ForceUnlockResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	current, err := s.repo.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get operation: %v", err)
	}

	op, err := s.repo.Terminate(id, current.FencingToken, db.OperationFailed, db.OperationAudit{
		Event:  "failed",
		Actor:  req.Actor,
		Reason: req.Reason,
	}, req.Reason)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail operation: %v", err)
	}

	log.Warnf("operation %s failed by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.ForceUnlockResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) ForceUnlock(ctx context.Context, req *pb.ForceUnlockRequest) (*pb.ForceUnlockResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	current, err := s.repo.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get operation: %v", err)
	}

	op, err := s.repo.Terminate(id, current.FencingToken, db.OperationCancelled, db.OperationAudit{
		Event:  "force_unlock",
		Actor:  req.Actor,
		Reason: req.Reason,
	}, "force unlocked: "+req.Reason)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "force unlock: %v", err)
	}

	log.Warnf("operation %s force-unlocked by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.ForceUnlockResponse{Operation: toPb(op)}, nil
}

func toPb(o *db.Operation) *pb.Operation {
	if o == nil {
		return nil
	}
	out := &pb.Operation{
		Id:             o.Id.String(),
		Type:           o.Type,
		System:         o.System,
		Status:         pb.OperationStatus(o.Status),
		FencingToken:   o.FencingToken,
		RequestedBy:    o.RequestedBy,
		ResourceKey:    o.ResourceKey,
		LeaseExpiresAt: timestamppb.New(o.LeaseExpiresAt),
		Error:          o.Error,
		CreatedAt:      timestamppb.New(o.CreatedAt),
	}
	if o.IdempotencyKey != nil {
		out.IdempotencyKey = *o.IdempotencyKey
	}
	if o.StartedAt != nil {
		out.StartedAt = timestamppb.New(*o.StartedAt)
	}
	if o.TerminalAt != nil {
		out.TerminalAt = timestamppb.New(*o.TerminalAt)
	}
	return out
}
