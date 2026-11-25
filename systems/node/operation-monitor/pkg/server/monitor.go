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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
)

const (
	appRulePrefix = "app."

	defaultDeadlineTTL = 4 * time.Minute
	sweeperInterval    = 30 * time.Second

	operationMonitorActor        = "operation-monitor"
	operationMonitorSweeperActor = "operation-monitor-sweeper"

	reasonDeadlineExceeded = "deadline exceeded"

	nodeResourcePrefix = "node:"

	ruleKeyState  = "state"
	ruleKeyNodeId = "node_id"
	substateOff   = "off"
)

// Action → completion rule fallback when caller doesn't supply one.
// A reboot takes the node offline then back online (substate: on→off→on)
// without changing the main state (Operational stays Operational), so
// completion is the node reporting substate=on again.
// TODO: move to config or per-action proto when we add more actions.
var defaultCompletionRule = map[string]string{
	"SendNodeCommand":      pkg.RuleSubstateOn,
	"RestartNode":          pkg.RuleSubstateOn,
	"ToggleRadio":          pkg.RuleSubstateOn,
	"ToggleInternetSwitch": pkg.RuleSubstateOn,
	"ToggleService":        pkg.RuleSubstateOn,
}

type MonitorServer struct {
	pb.UnimplementedOperationMonitorServiceServer
	orgName        string
	orgId          string
	repo           db.IntentRepo
	msgbus         mb.MsgBusServiceClient
	publishBuilder msgbus.RoutingKeyBuilder
}

func NewMonitorServer(orgName, orgId string, repo db.IntentRepo,
	msgBus mb.MsgBusServiceClient) *MonitorServer {
	return &MonitorServer{
		orgName: orgName,
		orgId:   orgId,
		repo:    repo,
		msgbus:  msgBus,
		publishBuilder: msgbus.NewRoutingKeyBuilder().
			SetCloudSource().
			SetGlobalScope().
			SetSystem(pkg.SystemName).
			SetService(pkg.ServiceName).
			SetOrgName(orgName),
	}
}

func (s *MonitorServer) RegisterIntent(ctx context.Context,
	req *pb.RegisterIntentRequest) (*pb.RegisterIntentResponse, error) {
	opId, err := uuid.FromString(req.OperationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation id: %v", err)
	}

	rule := req.CompletionRule
	if rule == "" {
		if defaultRule, ok := defaultCompletionRule[req.ActionType]; ok {
			rule = defaultRule
		} else {
			return nil, status.Errorf(codes.InvalidArgument,
				"completion_rule is empty and no default for action %q", req.ActionType)
		}
	}

	deadline := time.Duration(req.DeadlineSeconds) * time.Second
	if deadline == 0 {
		deadline = defaultDeadlineTTL
	}

	intent := &db.Intent{
		Id:             uuid.NewV4(),
		OperationId:    opId,
		ResourceKey:    req.ResourceKey,
		ActionType:     req.ActionType,
		FencingToken:   req.FencingToken,
		CompletionRule: rule,
		Status:         db.IntentWatching,
		Deadline:       time.Now().UTC().Add(deadline),
	}
	if err := s.repo.Add(intent); err != nil {
		return nil, grpc.StoreErrorToGrpc(err, "intent")
	}

	log.Infof("monitor: registered intent op=%s resource=%s rule=%q deadline=%s",
		intent.OperationId, intent.ResourceKey, intent.CompletionRule, intent.Deadline)
	return &pb.RegisterIntentResponse{Intent: toPb(intent)}, nil
}

func (s *MonitorServer) GetIntent(ctx context.Context,
	req *pb.GetIntentRequest) (*pb.GetIntentResponse, error) {
	opId, err := uuid.FromString(req.OperationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation id: %v", err)
	}
	intent, err := s.repo.Get(opId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "intent not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get intent: %v", err)
	}
	return &pb.GetIntentResponse{Intent: toPb(intent)}, nil
}

func (s *MonitorServer) CancelIntent(ctx context.Context,
	req *pb.CancelIntentRequest) (*pb.CancelIntentResponse, error) {
	opId, err := uuid.FromString(req.OperationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation id: %v", err)
	}
	if _, err := s.repo.Update(opId, closeAs(db.IntentCancelled)); err != nil {
		return nil, grpc.StoreErrorToGrpc(err, "intent")
	}
	return &pb.CancelIntentResponse{}, nil
}

func toPb(i *db.Intent) *pb.Intent {
	if i == nil {
		return nil
	}
	return &pb.Intent{
		Id:             i.Id.String(),
		OperationId:    i.OperationId.String(),
		ResourceKey:    i.ResourceKey,
		ActionType:     i.ActionType,
		FencingToken:   i.FencingToken,
		CompletionRule: i.CompletionRule,
		Status:         pb.IntentStatus(i.Status),
		Deadline:       timestamppb.New(i.Deadline),
		CreatedAt:      timestamppb.New(i.CreatedAt),
	}
}

type Sweeper struct {
	monitor  *MonitorServer
	interval time.Duration
	batch    int
}

func NewSweeper(monitor *MonitorServer) *Sweeper {
	return &Sweeper{
		monitor:  monitor,
		interval: sweeperInterval,
		batch:    100,
	}
}

func (s *Sweeper) Run(ctx context.Context) {
	log.Infof("operation-monitor sweeper started (interval=%s)", s.interval)
	t := time.NewTicker(s.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("operation-monitor sweeper stopped")
			return
		case <-t.C:
			s.sweepOnce()
		}
	}
}

func (s *Sweeper) sweepOnce() {
	expired, err := s.monitor.repo.List(expiredBefore(time.Now().UTC(), s.batch))
	if err != nil {
		log.Errorf("sweeper: list expired intents: %v", err)
		return
	}
	if len(expired) == 0 {
		return
	}
	log.Infof("sweeper: expiring %d intent(s) past deadline", len(expired))
	for i := range expired {
		intent := &expired[i]
		// Leave publication failures watching for the next sweep.
		if err := s.publishExpired(intent, reasonDeadlineExceeded); err != nil {
			log.Warnf("sweeper: publish failed for %s: %v", intent.OperationId, err)
			continue
		}
		if _, err := s.monitor.repo.Update(intent.OperationId,
			closeAs(db.IntentExpired)); err != nil {
			log.Warnf("sweeper: mark expired for %s failed: %v", intent.OperationId, err)
		}
	}
}

func (s *Sweeper) publishExpired(intent *db.Intent, reason string) error {
	route := s.monitor.publishBuilder.SetAction("expired").SetObject("intent").MustBuild()
	return s.monitor.msgbus.PublishRequest(route, &epb.IntentExpiredEvent{
		OperationId:  intent.OperationId.String(),
		FencingToken: intent.FencingToken,
		ResourceKey:  intent.ResourceKey,
		Reason:       reason,
		ExpiredAt:    timestamppb.Now(),
		Actor:        operationMonitorSweeperActor,
	})
}

func watchingOn(resourceKey string) db.IntentFilter {
	st := db.IntentWatching
	return db.IntentFilter{ResourceKey: resourceKey, Status: &st}
}

func expiredBefore(now time.Time, limit int) db.IntentFilter {
	st := db.IntentWatching
	return db.IntentFilter{Status: &st, DeadlineBefore: now, Limit: limit}
}

func closeAs(st db.IntentStatus) db.IntentUpdate {
	return db.IntentUpdate{Status: &st}
}

func arm() db.IntentUpdate {
	armed := true
	return db.IntentUpdate{Armed: &armed}
}
