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

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

type EventServer struct {
	epb.UnimplementedEventNotificationServiceServer
	orgName string
	repo    db.OperationRepo
}

func NewEventServer(orgName string, repo db.OperationRepo) *EventServer {
	return &EventServer{orgName: orgName, repo: repo}
}

func (e *EventServer) EventNotification(ctx context.Context, event *epb.Event) (*epb.EventResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("nil event")
	}
	log.Infof("operation/manager event-consumer: received %s", event.RoutingKey)

	switch event.RoutingKey {
	case msgbus.PrepareRoute(e.orgName, evt.EventRoutingKey[evt.EventOperationCompleted]):
		return e.handleCompleted(ctx, event)
	case msgbus.PrepareRoute(e.orgName, evt.EventRoutingKey[evt.EventOperationFailed]):
		return e.handleFailed(ctx, event)
	default:
		log.Warnf("operation/manager event-consumer: unhandled routing key %s", event.RoutingKey)
		return &epb.EventResponse{}, nil
	}
}

func (e *EventServer) handleCompleted(_ context.Context, event *epb.Event) (*epb.EventResponse, error) {
	msg, err := epb.UnmarshalOperationCompletedEvent(event.Msg, "OperationCompletedEvent")
	if err != nil {
		return nil, err
	}
	return e.terminate(msg.OperationId, msg.FencingToken, db.OperationSuccess,
		db.OperationAudit{Event: "completed"}, "")
}

func (e *EventServer) handleFailed(_ context.Context, event *epb.Event) (*epb.EventResponse, error) {
	msg, err := epb.UnmarshalOperationFailedEvent(event.Msg, "OperationFailedEvent")
	if err != nil {
		return nil, err
	}
	return e.terminate(msg.OperationId, msg.FencingToken, db.OperationFailed,
		db.OperationAudit{Event: "failed", Reason: msg.Reason}, msg.Reason)
}

func (e *EventServer) terminate(operationId string, fencingToken uint64,
	st db.OperationStatus, audit db.OperationAudit,
	opErr string) (*epb.EventResponse, error) {
	opId, err := uuid.FromString(operationId)
	if err != nil {
		return nil, fmt.Errorf("invalid operation id %q: %w", operationId, err)
	}

	_, err = e.repo.Terminate(opId, fencingToken, st, audit, opErr)
	switch {
	case err == nil:
		log.Infof("operation %s → %s (lock released)", opId, st)
	case errors.Is(err, store.ErrPrecondition):
		log.Warnf("operation/manager: ignoring %s event for op %s: %v",
			audit.Event, opId, err)
	case errors.Is(err, gorm.ErrRecordNotFound):
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
