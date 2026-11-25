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
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"

	log "github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type EventServer struct {
	epb.UnimplementedEventNotificationServiceServer
	orgName string
	repos   Repos
}

func NewEventServer(orgName string, repos Repos) *EventServer {
	return &EventServer{orgName: orgName, repos: repos}
}

func (e *EventServer) EventNotification(ctx context.Context,
	event *epb.Event) (*epb.EventResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("nil event")
	}
	log.Infof("operation/manager event-consumer: received %s", event.RoutingKey)

	switch event.RoutingKey {
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.global.{{ .Org}}.node.operationmonitor.intent.satisfied"):
		return e.handleSatisfied(ctx, event)
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.global.{{ .Org}}.node.operationmonitor.intent.expired"):
		return e.handleExpired(ctx, event)
	default:
		log.Warnf("operation/manager event-consumer: unhandled routing key %s",
			event.RoutingKey)
		return &epb.EventResponse{}, nil
	}
}

func (e *EventServer) handleSatisfied(ctx context.Context,
	event *epb.Event) (*epb.EventResponse, error) {
	msg, err := cpb.UnmarshalProtoEvent[epb.IntentSatisfiedEvent](event.Msg)
	if err != nil {
		return nil, err
	}
	return e.terminate(ctx, msg.OperationId, msg.FencingToken, db.OperationSuccess,
		db.OperationAudit{
			Event: pkg.AuditCompleted, Actor: eventActor(msg.Actor),
			At: eventTime(msg.ObservedAt),
		}, "")
}

func (e *EventServer) handleExpired(ctx context.Context,
	event *epb.Event) (*epb.EventResponse, error) {
	msg, err := cpb.UnmarshalProtoEvent[epb.IntentExpiredEvent](event.Msg)
	if err != nil {
		return nil, err
	}
	return e.terminate(ctx, msg.OperationId, msg.FencingToken, db.OperationTimeout,
		db.OperationAudit{
			Event: pkg.AuditTimeout, Actor: eventActor(msg.Actor), Reason: msg.Reason,
			At: eventTime(msg.ExpiredAt),
		}, msg.Reason)
}

func eventTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func eventActor(actor string) string {
	if strings.TrimSpace(actor) == "" {
		return operationMonitorActor
	}
	return actor
}

func (e *EventServer) terminate(ctx context.Context, operationId string,
	fencingToken uint64, st db.OperationStatus, audit db.OperationAudit,
	opErr string) (*epb.EventResponse, error) {
	opId, err := uuid.FromString(operationId)
	if err != nil {
		return nil, fmt.Errorf("invalid operation id %q: %w", operationId, err)
	}

	_, err = e.repos.terminate(ctx, opId, fencingToken, st, audit, opErr)
	switch {
	case err == nil:
		log.Infof("operation %s → %s (lock released)", opId, st)
	case errors.Is(err, store.ErrPrecondition):
		log.Warnf("operation/manager: ignoring %s event for op %s: %v",
			audit.Event, opId, err)
	case errors.Is(err, store.ErrNotFound):
		log.Warnf("operation/manager: ignoring %s event for unknown op %s",
			audit.Event, opId)
	default:
		log.Errorf("operation/manager: %s event for op %s failed: %v",
			audit.Event, opId, err)
		return nil, status.Errorf(codes.Unavailable, "terminate operation %s: %v",
			opId, err)
	}
	return &epb.EventResponse{}, nil
}
