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
	"sort"
	"strings"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	cgrpc "github.com/ukama/ukama/systems/common/grpc"
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

const maxResourceKeyLen = 255

func validateResourceKey(key string) error {
	if key == "" || len(key) > maxResourceKeyLen {
		return status.Errorf(codes.InvalidArgument,
			"resource key must be 1 to %d bytes", maxResourceKeyLen)
	}
	for _, r := range key {
		if r == unicode.ReplacementChar || !unicode.IsPrint(r) || unicode.IsSpace(r) {
			return status.Error(codes.InvalidArgument,
				"resource key must contain only printable, non-space characters")
		}
	}
	return nil
}

func (s *OperationServer) StartOperation(ctx context.Context, req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	targets, err := validateResourceKeys("resourceKeys", req.ResourceKeys)
	if err != nil {
		return nil, err
	}
	if len(targets) == 0 {
		return nil, status.Error(codes.InvalidArgument,
			"resourceKeys must contain at least one key")
	}
	mustBeFree, err := validateResourceKeys("mustBeFree", req.MustBeFree)
	if err != nil {
		return nil, err
	}
	if err := requireText("requestedBy", req.RequestedBy); err != nil {
		return nil, err
	}

	lease := time.Duration(req.LeaseSeconds) * time.Second
	if lease > pkg.MaxLeaseTTL {
		return nil, status.Errorf(codes.InvalidArgument,
			"lease seconds must be at most %d", int(pkg.MaxLeaseTTL.Seconds()))
	}
	if lease == 0 {
		lease = pkg.DefaultLeaseTTL
	}

	if req.IdempotencyKey != "" {
		resp, err := s.replay(ctx, req, targets)
		if resp != nil || err != nil {
			return resp, err
		}
	}

	ops, err := s.repo.Start(ctx, newGroup(req, targets, lease), mustBeFree, lease)
	if req.IdempotencyKey != "" && (errors.Is(err, db.ErrIdempotencyKeyTaken) ||
		errors.Is(err, db.ErrLockConflict)) {
		resp, rerr := s.replay(ctx, req, targets)
		if resp != nil || rerr != nil {
			return resp, rerr
		}
		if errors.Is(err, db.ErrIdempotencyKeyTaken) {
			return nil, status.Errorf(codes.Aborted,
				"idempotency key %q is in use by a concurrent start", req.IdempotencyKey)
		}
	}
	if errors.Is(err, db.ErrLockConflict) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return nil, repoErr(err)
	}

	for _, op := range ops {
		log.Infof("operation %s started: type=%s resource=%s token=%d", op.Id, op.Type,
			op.ResourceKey, op.FencingToken)
	}
	return &pb.StartOperationResponse{Operations: toPbList(ops)}, nil
}

func newGroup(req *pb.StartOperationRequest, targets []string,
	lease time.Duration) []*db.Operation {
	expires := time.Now().UTC().Add(lease)
	ops := make([]*db.Operation, 0, len(targets))
	for _, key := range targets {
		ops = append(ops, &db.Operation{
			Id:             uuid.NewV4(),
			Type:           req.Type,
			System:         req.System,
			Status:         db.OperationPending,
			RequestedBy:    req.RequestedBy,
			ResourceKey:    key,
			LeaseExpiresAt: expires,
		})
	}

	leader := ops[0]
	if req.IdempotencyKey != "" {
		k := req.IdempotencyKey
		leader.IdempotencyKey = &k
	}
	for _, member := range ops[1:] {
		member.GroupId = &leader.Id
	}
	return ops
}

func (s *OperationServer) replay(ctx context.Context, req *pb.StartOperationRequest,
	targets []string) (*pb.StartOperationResponse, error) {
	ops, err := s.repo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return nil, repoErr(err)
	}
	if len(ops) == 0 {
		return nil, nil
	}
	if !sameRequest(ops, req, targets) {
		keys := make([]string, 0, len(ops))
		for _, op := range ops {
			keys = append(keys, op.ResourceKey)
		}
		return nil, status.Errorf(codes.InvalidArgument,
			"idempotency key %q belongs to %s %s operation on %s",
			req.IdempotencyKey, ops[0].System, ops[0].Type, strings.Join(keys, ","))
	}
	return &pb.StartOperationResponse{Operations: toPbList(ops)}, nil
}

func sameRequest(ops []*db.Operation, req *pb.StartOperationRequest,
	targets []string) bool {
	if len(ops) != len(targets) {
		return false
	}
	for i, op := range ops {
		if op.Type != req.Type || op.System != req.System ||
			op.ResourceKey != targets[i] {
			return false
		}
	}
	return true
}

func validateResourceKeys(field string, keys []string) ([]string, error) {
	seen := make(map[string]bool, len(keys))
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		if err := validateResourceKey(key); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%s: %s", field,
				status.Convert(err).Message())
		}
		if !seen[key] {
			seen[key] = true
			out = append(out, key)
		}
	}
	sort.Strings(out)
	return out, nil
}

func (s *OperationServer) GetOperation(ctx context.Context, req *pb.GetOperationRequest) (*pb.GetOperationResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	op, err := s.repo.Get(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, repoErr(err)
	}
	return &pb.GetOperationResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) GetByResource(ctx context.Context, req *pb.GetByResourceRequest) (*pb.GetByResourceResponse, error) {
	if err := validateResourceKey(req.ResourceKey); err != nil {
		return nil, err
	}
	op, err := s.repo.GetByResource(ctx, req.ResourceKey)
	if err != nil {
		return nil, repoErr(err)
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
	if err := requireText("actor", req.Actor); err != nil {
		return nil, err
	}
	if req.FencingToken == 0 {
		return nil, status.Error(codes.InvalidArgument, "fencing token must be set")
	}
	op, err := s.repo.MarkRunning(ctx, id, req.FencingToken, req.Actor)
	if err != nil {
		return nil, repoErr(err)
	}
	return &pb.MarkRunningResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) CompleteOperation(ctx context.Context,
	req *pb.CompleteOperationRequest) (*pb.CompleteOperationResponse, error) {
	if err := requireActorAndReason(req.Actor, req.Reason); err != nil {
		return nil, err
	}
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	op, err := s.repo.Terminate(ctx, id, req.FencingToken, db.OperationSuccess,
		db.OperationAudit{
			Event:  "completed",
			Actor:  req.Actor,
			Reason: req.Reason,
		}, "")
	if err != nil {
		return nil, repoErr(err)
	}

	log.Infof("operation %s completed by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.CompleteOperationResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) FailOperation(ctx context.Context,
	req *pb.FailOperationRequest) (*pb.FailOperationResponse, error) {
	if err := requireActorAndReason(req.Actor, req.Reason); err != nil {
		return nil, err
	}
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	op, err := s.repo.Terminate(ctx, id, req.FencingToken, db.OperationFailed,
		db.OperationAudit{
			Event:  "failed",
			Actor:  req.Actor,
			Reason: req.Reason,
		}, req.Reason)
	if err != nil {
		return nil, repoErr(err)
	}

	log.Warnf("operation %s failed by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.FailOperationResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) ForceUnlock(ctx context.Context, req *pb.ForceUnlockRequest) (*pb.ForceUnlockResponse, error) {
	if err := requireActorAndReason(req.Actor, req.Reason); err != nil {
		return nil, err
	}
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	current, err := s.repo.Get(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, repoErr(err)
	}

	op, err := s.repo.Terminate(ctx, id, current.FencingToken, db.OperationCancelled,
		db.OperationAudit{
			Event:  "force_unlock",
			Actor:  req.Actor,
			Reason: req.Reason,
		}, "force unlocked: "+req.Reason)
	if err != nil {
		return nil, repoErr(err)
	}

	log.Warnf("operation %s force-unlocked by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.ForceUnlockResponse{Operation: toPb(op)}, nil
}

func requireText(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return status.Errorf(codes.InvalidArgument, "%s must not be blank", field)
	}
	return nil
}

func requireActorAndReason(actor, reason string) error {
	if err := requireText("actor", actor); err != nil {
		return err
	}
	return requireText("reason", reason)
}

func repoErr(err error) error {
	if errors.Is(err, db.ErrLockTimeout) {
		return status.Error(codes.Unavailable, err.Error())
	}
	return cgrpc.StoreErrorToGrpc(err, "operation")
}

func toPbList(ops []*db.Operation) []*pb.Operation {
	out := make([]*pb.Operation, 0, len(ops))
	for _, op := range ops {
		out = append(out, toPb(op))
	}
	return out
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
	if o.GroupId != nil {
		out.GroupId = o.GroupId.String()
	}
	if o.StartedAt != nil {
		out.StartedAt = timestamppb.New(*o.StartedAt)
	}
	if o.TerminalAt != nil {
		out.TerminalAt = timestamppb.New(*o.TerminalAt)
	}
	return out
}
