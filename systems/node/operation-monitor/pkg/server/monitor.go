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

	pb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"
)

type MonitorServer struct {
	pb.UnimplementedOperationMonitorServiceServer
	orgName        string
	orgId          string
	repo           db.IntentRepo
	msgbus         mb.MsgBusServiceClient
	publishBuilder msgbus.RoutingKeyBuilder
}

func NewMonitorServer(orgName, orgId string, repo db.IntentRepo, msgBus mb.MsgBusServiceClient) *MonitorServer {
	return &MonitorServer{
		orgName: orgName,
		orgId:   orgId,
		repo:    repo,
		msgbus:  msgBus,
		// Publishes operation.* events at GLOBAL scope as if from operation/manager,
		// so the global manager's event-consumer (subscribed to those keys) receives them.
		publishBuilder: msgbus.NewRoutingKeyBuilder().
			SetCloudSource().
			SetGlobalScope().
			SetSystem("operation").
			SetService("manager").
			SetOrgName(orgName),
	}
}

func (s *MonitorServer) RegisterIntent(ctx context.Context, req *pb.RegisterIntentRequest) (*pb.RegisterIntentResponse, error) {
	opId, err := uuid.FromString(req.OperationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation id: %v", err)
	}

	rule := req.CompletionRule
	if rule == "" {
		if defaultRule, ok := pkg.DefaultCompletionRule[req.ActionType]; ok {
			rule = defaultRule
		} else {
			return nil, status.Errorf(codes.InvalidArgument,
				"completion_rule is empty and no default for action %q", req.ActionType)
		}
	}

	deadline := time.Duration(req.DeadlineSeconds) * time.Second
	if deadline == 0 {
		deadline = pkg.DefaultDeadlineTTL
	}

	intent := &db.MonitoredIntent{
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
		return nil, status.Errorf(codes.Internal, "register intent: %v", err)
	}

	log.Infof("monitor: registered intent op=%s resource=%s rule=%q deadline=%s",
		intent.OperationId, intent.ResourceKey, intent.CompletionRule, intent.Deadline)
	return &pb.RegisterIntentResponse{Intent: toPb(intent)}, nil
}

func (s *MonitorServer) GetIntent(ctx context.Context, req *pb.GetIntentRequest) (*pb.GetIntentResponse, error) {
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

func (s *MonitorServer) CancelIntent(ctx context.Context, req *pb.CancelIntentRequest) (*pb.CancelIntentResponse, error) {
	opId, err := uuid.FromString(req.OperationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation id: %v", err)
	}
	if _, err := s.repo.MarkTerminal(opId, db.IntentCancelled); err != nil {
		return nil, status.Errorf(codes.Internal, "cancel intent: %v", err)
	}
	return &pb.CancelIntentResponse{}, nil
}

func toPb(i *db.MonitoredIntent) *pb.MonitoredIntent {
	if i == nil {
		return nil
	}
	return &pb.MonitoredIntent{
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
