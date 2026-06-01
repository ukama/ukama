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
	"fmt"

	log "github.com/sirupsen/logrus"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
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
	opId, err := uuid.FromString(msg.OperationId)
	if err != nil {
		return nil, fmt.Errorf("invalid operation id %q: %w", msg.OperationId, err)
	}
	if _, err := e.repo.Terminate(opId, msg.FencingToken, db.OperationSuccess,
		db.OperationAudit{Event: "completed"}, ""); err != nil {
		log.Errorf("operation/manager: terminate(success) failed for op %s: %v", opId, err)
		return nil, err
	}
	log.Infof("operation %s → success (lock released)", opId)
	return &epb.EventResponse{}, nil
}

func (e *EventServer) handleFailed(_ context.Context, event *epb.Event) (*epb.EventResponse, error) {
	msg, err := epb.UnmarshalOperationFailedEvent(event.Msg, "OperationFailedEvent")
	if err != nil {
		return nil, err
	}
	opId, err := uuid.FromString(msg.OperationId)
	if err != nil {
		return nil, fmt.Errorf("invalid operation id %q: %w", msg.OperationId, err)
	}
	if _, err := e.repo.Terminate(opId, msg.FencingToken, db.OperationFailed,
		db.OperationAudit{Event: "failed", Reason: msg.Reason}, msg.Reason); err != nil {
		log.Errorf("operation/manager: terminate(failed) failed for op %s: %v", opId, err)
		return nil, err
	}
	log.Infof("operation %s → failed (lock released): %s", opId, msg.Reason)
	return &epb.EventResponse{}, nil
}
