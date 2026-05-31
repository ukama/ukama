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
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

const (
	RouteOperationCompleted = "event.cloud.global.{{ .Org}}.operation.manager.operation.completed"
	RouteOperationFailed    = "event.cloud.global.{{ .Org}}.operation.manager.operation.failed"
)

type EventServer struct {
	epb.UnimplementedEventNotificationServiceServer
	orgName string
	repo    db.OperationRepo
}

func NewEventServer(orgName string, repo db.OperationRepo) *EventServer {
	return &EventServer{orgName: orgName, repo: repo}
}

func (e *EventServer) EventNotification(ctx context.Context, evt *epb.Event) (*epb.EventResponse, error) {
	if evt == nil {
		return nil, fmt.Errorf("nil event")
	}
	log.Infof("operation/manager event-consumer: received %s", evt.RoutingKey)

	switch evt.RoutingKey {
	case msgbus.PrepareRoute(e.orgName, RouteOperationCompleted):
		return e.handleCompleted(ctx, evt.Msg)
	case msgbus.PrepareRoute(e.orgName, RouteOperationFailed):
		return e.handleFailed(ctx, evt.Msg)
	default:
		log.Warnf("operation/manager event-consumer: unhandled routing key %s", evt.RoutingKey)
		return &epb.EventResponse{}, nil
	}
}

func (e *EventServer) handleCompleted(_ context.Context, raw *anypb.Any) (*epb.EventResponse, error) {
	opId, token, _, err := unmarshalOperationTerminalEvent(raw)
	if err != nil {
		return nil, err
	}
	_, err = e.repo.Terminate(opId, token, db.OperationSuccess, db.OperationAudit{
		Event: "completed",
	}, "")
	if err != nil {
		log.Errorf("operation/manager: terminate(success) failed for op %s: %v", opId, err)
		return nil, err
	}
	log.Infof("operation %s → success (lock released)", opId)
	return &epb.EventResponse{}, nil
}

func (e *EventServer) handleFailed(_ context.Context, raw *anypb.Any) (*epb.EventResponse, error) {
	opId, token, reason, err := unmarshalOperationTerminalEvent(raw)
	if err != nil {
		return nil, err
	}
	_, err = e.repo.Terminate(opId, token, db.OperationFailed, db.OperationAudit{
		Event:  "failed",
		Reason: reason,
	}, reason)
	if err != nil {
		log.Errorf("operation/manager: terminate(failed) failed for op %s: %v", opId, err)
		return nil, err
	}
	log.Infof("operation %s → failed (lock released): %s", opId, reason)
	return &epb.EventResponse{}, nil
}

// TODO: replace with epb.UnmarshalOperationTerminalEvent once
// OperationTerminalEvent is added to systems/common/pb/events.
func unmarshalOperationTerminalEvent(raw *anypb.Any) (uuid.UUID, uint64, string, error) {
	if raw == nil {
		return uuid.UUID{}, 0, "", fmt.Errorf("empty event payload")
	}
	return uuid.UUID{}, 0, "", fmt.Errorf("OperationTerminalEvent unmarshal not yet implemented")
}
