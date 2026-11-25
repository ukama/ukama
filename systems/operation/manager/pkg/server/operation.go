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
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"

	log "github.com/sirupsen/logrus"
	cgrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

type OperationServer struct {
	pb.UnimplementedOperationManagerServiceServer
	orgName string
	orgId   string
	repos   Repos
}

func NewOperationServer(orgName, orgId string, repos Repos) *OperationServer {
	return &OperationServer{orgName: orgName, orgId: orgId, repos: repos}
}

const (
	maxResourceKeyLen = 255
	defaultLeaseTTL   = 5 * time.Minute
	maxLeaseTTL       = time.Hour
	sweeperInterval   = 30 * time.Second

	operationManagerSweeperActor = "operation-manager-sweeper"
	operationMonitorActor        = "operation-monitor"
	operationMonitorSweeperActor = "operation-monitor-sweeper"
)

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

func (s *OperationServer) StartOperation(ctx context.Context,
	req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
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
	if lease > maxLeaseTTL {
		return nil, status.Errorf(codes.InvalidArgument,
			"lease seconds must be at most %d", int(maxLeaseTTL.Seconds()))
	}
	if lease == 0 {
		lease = defaultLeaseTTL
	}

	if req.IdempotencyKey != "" {
		resp, err := s.replay(ctx, req, targets)
		if resp != nil || err != nil {
			return resp, err
		}
	}

	ops := newGroup(req, targets, lease)
	err = s.repos.Operations.Add(ctx, ops, func(tx *gorm.DB) error {
		if err := s.repos.Locks.Acquire(tx, ops, mustBeFree); err != nil {
			return err
		}
		for _, op := range ops {
			if err := s.repos.Audits.AddInTx(tx, &db.OperationAudit{
				OperationId: op.Id, ResourceKey: op.ResourceKey,
				Event: pkg.AuditLockAcquired, Actor: op.RequestedBy,
			}); err != nil {
				return err
			}
		}
		return nil
	})
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
	ops, err := s.groupByIdempotencyKey(ctx, req.IdempotencyKey)
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

func (s *OperationServer) groupByIdempotencyKey(ctx context.Context,
	key string) ([]*db.Operation, error) {
	leaders, err := s.repos.Operations.List(ctx, db.OperationFilter{IdempotencyKey: key})
	if err != nil || len(leaders) == 0 {
		return nil, err
	}
	members, err := s.repos.Operations.List(ctx, db.OperationFilter{
		GroupId: &leaders[0].Id, OrderBy: db.OrderByResourceKey,
	})
	if err != nil {
		return nil, err
	}
	return append([]*db.Operation{leaders[0]}, members...), nil
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

func (s *OperationServer) GetOperation(ctx context.Context,
	req *pb.GetOperationRequest) (*pb.GetOperationResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	op, err := s.repos.Operations.Get(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, repoErr(err)
	}
	return &pb.GetOperationResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) GetByResource(ctx context.Context,
	req *pb.GetByResourceRequest) (*pb.GetByResourceResponse, error) {
	if err := validateResourceKey(req.ResourceKey); err != nil {
		return nil, err
	}
	lock, err := s.repos.Locks.Holder(ctx, req.ResourceKey)
	if err != nil {
		return nil, repoErr(err)
	}
	if lock == nil {
		return &pb.GetByResourceResponse{}, nil
	}
	op, err := s.repos.Operations.Get(ctx, lock.OperationId)
	if err != nil {
		return nil, repoErr(err)
	}
	return &pb.GetByResourceResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) MarkRunning(ctx context.Context,
	req *pb.MarkRunningRequest) (*pb.MarkRunningResponse, error) {
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
	op, err := s.repos.markRunning(ctx, id, req.FencingToken, req.Actor)
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

	op, err := s.repos.terminate(ctx, id, req.FencingToken, db.OperationSuccess,
		db.OperationAudit{
			Event:  pkg.AuditCompleted,
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

	op, err := s.repos.terminate(ctx, id, req.FencingToken, db.OperationFailed,
		db.OperationAudit{
			Event:  pkg.AuditFailed,
			Actor:  req.Actor,
			Reason: req.Reason,
		}, req.Reason)
	if err != nil {
		return nil, repoErr(err)
	}

	log.Warnf("operation %s failed by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.FailOperationResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) ForceUnlock(ctx context.Context,
	req *pb.ForceUnlockRequest) (*pb.ForceUnlockResponse, error) {
	if err := requireActorAndReason(req.Actor, req.Reason); err != nil {
		return nil, err
	}
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	current, err := s.repos.Operations.Get(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "operation not found")
	}
	if err != nil {
		return nil, repoErr(err)
	}

	op, err := s.repos.terminate(ctx, id, current.FencingToken, db.OperationCancelled,
		db.OperationAudit{
			Event:  pkg.AuditForceUnlock,
			Actor:  req.Actor,
			Reason: req.Reason,
		}, "force unlocked: "+req.Reason)
	if err != nil {
		return nil, repoErr(err)
	}

	log.Warnf("operation %s force-unlocked by %s: %s", op.Id, req.Actor, req.Reason)
	return &pb.ForceUnlockResponse{Operation: toPb(op)}, nil
}

func (s *OperationServer) DeleteOperation(ctx context.Context,
	req *pb.DeleteOperationRequest) (*pb.DeleteOperationResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	if err := requireText("actor", req.Actor); err != nil {
		return nil, err
	}
	err = s.repos.Operations.Delete(ctx, id, func(op *db.Operation, tx *gorm.DB) error {
		if !op.Status.IsTerminal() {
			return fmt.Errorf("operation %s is %s and cannot be deleted: %w",
				id, op.Status, store.ErrPrecondition)
		}
		return s.repos.Audits.AddInTx(tx, &db.OperationAudit{
			OperationId: op.Id, ResourceKey: op.ResourceKey, Event: pkg.AuditDeleted,
			Actor: req.Actor,
		})
	})
	if err != nil {
		return nil, repoErr(err)
	}

	log.Infof("operation %s deleted by %s", id, req.Actor)
	return &pb.DeleteOperationResponse{}, nil
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

type Repos struct {
	Operations db.OperationRepo
	Locks      db.LockRepo
	Audits     db.AuditRepo
	Results    ResultPublisher
}

type ResultPublisher interface {
	Publish(op *db.Operation, audit db.OperationAudit)
}

type resultPublisher struct {
	msgbus     mb.MsgBusServiceClient
	routingKey msgbus.RoutingKeyBuilder
}

func NewResultPublisher(orgName string, msgBus mb.MsgBusServiceClient) ResultPublisher {
	return &resultPublisher{
		msgbus: msgBus,
		routingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetGlobalScope().
			SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (p *resultPublisher) Publish(op *db.Operation, audit db.OperationAudit) {
	action, msg := resultEvent(op, audit)
	if msg == nil {
		return
	}
	route := p.routingKey.SetObject("operation").SetAction(action).MustBuild()
	if err := p.msgbus.PublishRequest(route, msg); err != nil {
		log.Errorf("operation/manager: publish %s result for op %s failed: %v", action,
			op.Id, err)
	}
}

func resultEvent(op *db.Operation, audit db.OperationAudit) (string,
	protoreflect.ProtoMessage) {
	id, token, key := op.Id.String(), op.FencingToken, op.ResourceKey
	at := timestamppb.New(audit.At)
	switch op.Status {
	case db.OperationSuccess:
		return "completed", &epb.OperationCompletedEvent{OperationId: id,
			FencingToken: token, ResourceKey: key, CompletedAt: at, Actor: audit.Actor}
	case db.OperationFailed:
		return "failed", &epb.OperationFailedEvent{OperationId: id, FencingToken: token,
			ResourceKey: key, Reason: op.Error, FailedAt: at, Actor: audit.Actor}
	case db.OperationTimeout:
		return "timeout", &epb.OperationTimedOutEvent{OperationId: id,
			FencingToken: token, ResourceKey: key, Reason: op.Error, TimedOutAt: at,
			Actor: audit.Actor}
	case db.OperationCancelled:
		return "cancelled", &epb.OperationCancelledEvent{OperationId: id,
			FencingToken: token, ResourceKey: key, Reason: op.Error, CancelledAt: at,
			Actor: audit.Actor}
	}
	return "", nil
}

type refusal struct {
	op     db.Operation
	prefix string
	reason string
	cause  error
}

func (r *refusal) Error() string {
	if r.cause != nil {
		return r.cause.Error()
	}
	return r.prefix + "transition"
}

func reject(op *db.Operation, cause error) *refusal {
	return &refusal{op: *op, prefix: pkg.AuditRejected, cause: cause}
}

func ignore(op *db.Operation, cause error) *refusal {
	return &refusal{op: *op, prefix: pkg.AuditIgnored, cause: cause}
}

func (r Repos) markRunning(ctx context.Context, id uuid.UUID, fencingToken uint64,
	actor string) (*db.Operation, error) {
	op, err := r.Operations.Update(ctx, id, func(op *db.Operation, tx *gorm.DB) error {
		now := time.Now().UTC()
		var ref *refusal
		switch {
		case op.FencingToken != fencingToken:
			ref = reject(op, fmt.Errorf("fencing token mismatch (op=%d, given=%d): %w",
				op.FencingToken, fencingToken, store.ErrPrecondition))
		case op.Status.IsTerminal():
			ref = reject(op, fmt.Errorf("operation %s is already %s: %w",
				id, op.Status, store.ErrPrecondition))
		case !op.LeaseExpiresAt.After(now):
			ref = reject(op, fmt.Errorf("operation %s lease expired at %s: %w",
				id, op.LeaseExpiresAt.Format(time.RFC3339), store.ErrPrecondition))
		case op.Status == db.OperationRunning:
			return ignore(op, nil)
		}
		if ref != nil {
			ref.reason = ref.cause.Error()
			return ref
		}

		op.Status = db.OperationRunning
		op.StartedAt = &now
		return r.Audits.AddInTx(tx, &db.OperationAudit{
			OperationId: op.Id, ResourceKey: op.ResourceKey, Event: pkg.AuditRunning,
			Actor: actor, At: now,
		})
	})
	return r.settle(ctx, op, err, db.OperationAudit{Event: pkg.AuditRunning, Actor: actor})
}

func (r Repos) terminate(ctx context.Context, id uuid.UUID, fencingToken uint64,
	status db.OperationStatus, audit db.OperationAudit,
	opErr string) (*db.Operation, error) {
	if !status.IsTerminal() {
		return nil, fmt.Errorf("terminate: %s is not a terminal status", status)
	}

	var final db.OperationAudit
	op, err := r.Operations.Update(ctx, id, func(op *db.Operation, tx *gorm.DB) error {
		if op.FencingToken != fencingToken {
			return reject(op, fmt.Errorf("fencing token mismatch (op=%d, given=%d): %w",
				op.FencingToken, fencingToken, store.ErrPrecondition))
		}
		if op.Status.IsTerminal() {
			var cause error
			if op.Status != status {
				cause = fmt.Errorf("operation %s is already %s, cannot move to %s: %w",
					id, op.Status, status, store.ErrPrecondition)
			}
			return ignore(op, cause)
		}
		if status == db.OperationSuccess && op.Status == db.OperationPending {
			return reject(op, fmt.Errorf("operation %s is pending, cannot move to %s: %w",
				id, status, store.ErrPrecondition))
		}

		at := endTime(op, audit.At)
		op.Status = status
		op.TerminalAt = &at
		if opErr != "" {
			op.Error = opErr
		}
		if err := r.Locks.Release(tx, op); err != nil {
			return err
		}
		final = audit
		final.OperationId = op.Id
		final.ResourceKey = op.ResourceKey
		final.At = at
		return r.Audits.AddInTx(tx, &final)
	})
	if err == nil && r.Results != nil {
		r.Results.Publish(op, final)
	}
	return r.settle(ctx, op, err, audit)
}

func endTime(op *db.Operation, observed time.Time) time.Time {
	now := time.Now().UTC()
	if observed.IsZero() || observed.After(now) || observed.Before(op.CreatedAt) {
		return now
	}
	return observed.UTC()
}

func (r Repos) settle(ctx context.Context, op *db.Operation, err error,
	audit db.OperationAudit) (*db.Operation, error) {
	var ref *refusal
	if !errors.As(err, &ref) {
		return op, err
	}

	audit.OperationId = ref.op.Id
	audit.ResourceKey = ref.op.ResourceKey
	audit.Event = ref.prefix + audit.Event
	if ref.reason != "" {
		audit.Reason = ref.reason
	}
	if aerr := r.Audits.Add(ctx, &audit); aerr != nil {
		return nil, aerr
	}
	if ref.cause != nil {
		return nil, ref.cause
	}
	return &ref.op, nil
}

type Sweeper struct {
	repos    Repos
	interval time.Duration
	batch    int
}

func NewSweeper(repos Repos) *Sweeper {
	return &Sweeper{
		repos:    repos,
		interval: sweeperInterval,
		batch:    100,
	}
}

func (s *Sweeper) Run(ctx context.Context) {
	log.Infof("operation/manager sweeper started (interval=%s)", s.interval)
	t := time.NewTicker(s.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("operation/manager sweeper stopped")
			return
		case <-t.C:
			s.sweepOnce(ctx)
		}
	}
}

func (s *Sweeper) sweepOnce(parent context.Context) {
	ctx, cancel := context.WithTimeout(parent, s.interval)
	defer cancel()

	now := time.Now().UTC()
	expired, err := s.repos.Operations.List(ctx, db.OperationFilter{
		Statuses:           []db.OperationStatus{db.OperationPending, db.OperationRunning},
		LeaseExpiresBefore: &now,
		OrderBy:            db.OrderByLeaseExpiry,
		Limit:              s.batch,
	})
	if err != nil {
		log.Errorf("sweeper: list expired operations: %v", err)
		return
	}
	if len(expired) == 0 {
		return
	}
	log.Infof("sweeper: timing out %d expired operation(s)", len(expired))
	for _, op := range expired {
		_, err := s.repos.terminate(ctx, op.Id, op.FencingToken, db.OperationTimeout,
			db.OperationAudit{
				Event:  pkg.AuditTimeout,
				Actor:  operationManagerSweeperActor,
				Reason: pkg.ReasonLeaseExpired,
			}, pkg.ReasonLeaseExpired)
		if err != nil {
			log.Warnf("sweeper: terminate(timeout) for %s failed: %v", op.Id, err)
		}
	}
}
