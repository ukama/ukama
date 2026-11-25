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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"

	log "github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type EventServer struct {
	epb.UnimplementedEventNotificationServiceServer
	orgName string
	monitor *MonitorServer
}

func NewEventServer(orgName string, m *MonitorServer) *EventServer {
	return &EventServer{orgName: orgName, monitor: m}
}

func (e *EventServer) EventNotification(ctx context.Context,
	event *epb.Event) (*epb.EventResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("nil event")
	}
	log.Infof("operation-monitor: received %s", event.RoutingKey)

	switch event.RoutingKey {
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.local.{{ .Org}}.node.state.node.transition"):
		return e.handleStateTransition(ctx, event)
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.local.{{ .Org}}.node.health.apps.changed"):
		return e.handleAppsChanged(ctx, event)
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.global.{{ .Org}}.operation.manager.operation.completed"):
		return e.handleResult(event, db.IntentCompleted,
			operationIdOf[epb.OperationCompletedEvent])
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.global.{{ .Org}}.operation.manager.operation.failed"):
		return e.handleResult(event, db.IntentFailed, operationIdOf[epb.OperationFailedEvent])
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.global.{{ .Org}}.operation.manager.operation.timeout"):
		return e.handleResult(event, db.IntentExpired,
			operationIdOf[epb.OperationTimedOutEvent])
	case msgbus.PrepareRoute(e.orgName,
		"event.cloud.global.{{ .Org}}.operation.manager.operation.cancelled"):
		return e.handleResult(event, db.IntentCancelled,
			operationIdOf[epb.OperationCancelledEvent])
	default:
		log.Warnf("operation-monitor: unhandled routing key %s", event.RoutingKey)
		return &epb.EventResponse{}, nil
	}
}

func (e *EventServer) handleStateTransition(_ context.Context,
	event *epb.Event) (*epb.EventResponse, error) {
	msg, err := cpb.UnmarshalProtoEvent[epb.NodeStateChangeEvent](event.Msg)
	if err != nil {
		log.Errorf("operation-monitor: unmarshal NodeStateChangeEvent: %v", err)
		return nil, err
	}

	nodeId := strings.ToLower(msg.NodeId)
	resourceKey := nodeResourcePrefix + nodeId
	intents, err := e.monitor.repo.List(watchingOn(resourceKey))
	if err != nil {
		log.Errorf("operation-monitor: lookup intents for %s: %v", resourceKey, err)
		return nil, retryable("lookup intents for %s: %v", resourceKey, err)
	}
	if len(intents) == 0 {
		return &epb.EventResponse{}, nil
	}

	transition := map[string]string{
		ruleKeyState:        msg.State,
		pkg.RuleKeySubstate: msg.Substate,
		ruleKeyNodeId:       nodeId,
	}

	for i := range intents {
		intent := &intents[i]
		if isAppRule(intent.CompletionRule) {
			continue
		}
		if !ruleMatches(intent.CompletionRule, transition) {
			if intent.Armed || !ruleObserved(intent.CompletionRule, transition) {
				continue
			}
			if _, err := e.monitor.repo.Update(intent.OperationId, arm()); err != nil {
				log.Errorf("operation-monitor: arm %s: %v", intent.OperationId, err)
				return nil, retryable("arm %s: %v", intent.OperationId, err)
			}
			log.Infof("operation-monitor: intent %s armed by %v",
				intent.OperationId, transition)
			continue
		}
		if !intent.Armed {
			log.Infof("operation-monitor: intent %s matched %v but not armed yet (no departure "+
				"observed), skipping",
				intent.OperationId, transition)
			continue
		}
		if err := e.complete(intent, transition); err != nil {
			return nil, err
		}
	}
	return &epb.EventResponse{}, nil
}

func (e *EventServer) handleAppsChanged(_ context.Context,
	event *epb.Event) (*epb.EventResponse, error) {
	msg, err := cpb.UnmarshalProtoEvent[epb.HealthAppsChangedEvent](event.Msg)
	if err != nil {
		log.Errorf("operation-monitor: unmarshal HealthAppsChangedEvent: %v", err)
		return nil, err
	}

	nodeId := strings.ToLower(msg.NodeId)
	resourceKey := nodeResourcePrefix + nodeId
	intents, err := e.monitor.repo.List(watchingOn(resourceKey))
	if err != nil {
		log.Errorf("operation-monitor: lookup intents for %s: %v", resourceKey, err)
		return nil, retryable("lookup intents for %s: %v", resourceKey, err)
	}
	if len(intents) == 0 {
		return &epb.EventResponse{}, nil
	}

	inventory := make(map[string]string, len(msg.Apps))
	for _, app := range msg.Apps {
		inventory[appRulePrefix+app.Name] = app.Version
	}

	for i := range intents {
		intent := &intents[i]
		if !isAppRule(intent.CompletionRule) {
			continue
		}
		if !ruleMatches(intent.CompletionRule, inventory) {
			log.Infof("operation-monitor: intent %s rule=%q not met by %v",
				intent.OperationId, intent.CompletionRule, inventory)
			continue
		}
		if err := e.complete(intent, inventory); err != nil {
			return nil, err
		}
	}
	return &epb.EventResponse{}, nil
}

func operationIdOf[T any, P interface {
	*T
	proto.Message
	GetOperationId() string
}](event *epb.Event) (string, error) {
	msg, err := cpb.UnmarshalProtoEvent[T, P](event.Msg)
	if err != nil {
		return "", err
	}
	return P(msg).GetOperationId(), nil
}

func (e *EventServer) handleResult(event *epb.Event, st db.IntentStatus,
	operationId func(*epb.Event) (string, error)) (*epb.EventResponse, error) {
	id, err := operationId(event)
	if err != nil {
		log.Errorf("operation-monitor: unmarshal %s: %v", event.RoutingKey, err)
		return nil, err
	}
	opId, err := uuid.FromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid operation id %q: %w", id, err)
	}
	intent, err := e.monitor.repo.Update(opId, closeAs(st))
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return &epb.EventResponse{}, nil
	case err != nil:
		log.Errorf("operation-monitor: close intent %s: %v", opId, err)
		return nil, retryable("close intent %s: %v", opId, err)
	}
	log.Infof("operation-monitor: intent %s is %s after operation result", opId,
		intent.Status)
	return &epb.EventResponse{}, nil
}

func (e *EventServer) complete(intent *db.Intent, observed map[string]string) error {
	if err := e.publishSatisfied(intent); err != nil {
		log.Errorf("operation-monitor: publish satisfied for %s: %v", intent.OperationId, err)
		return retryable("publish satisfied for %s: %v", intent.OperationId, err)
	}
	if _, err := e.monitor.repo.Update(intent.OperationId,
		closeAs(db.IntentCompleted)); err != nil {
		log.Errorf("operation-monitor: mark %s completed: %v", intent.OperationId, err)
		return retryable("mark %s completed: %v", intent.OperationId, err)
	}
	log.Infof("operation-monitor: intent %s satisfied (rule=%q matched %v)",
		intent.OperationId, intent.CompletionRule, observed)
	return nil
}

func retryable(format string, args ...any) error {
	return status.Errorf(codes.Unavailable, format, args...)
}

func ruleObserved(rule string, transition map[string]string) bool {
	for _, pair := range strings.Split(strings.TrimSpace(rule), ",") {
		k, _, _ := strings.Cut(pair, "=")
		if transition[strings.TrimSpace(k)] == "" {
			return false
		}
	}
	return true
}

func isAppRule(rule string) bool {
	return strings.HasPrefix(strings.TrimSpace(rule), appRulePrefix)
}

func (e *EventServer) publishSatisfied(intent *db.Intent) error {
	route := e.monitor.publishBuilder.SetAction("satisfied").SetObject("intent").MustBuild()
	return e.monitor.msgbus.PublishRequest(route, &epb.IntentSatisfiedEvent{
		OperationId:  intent.OperationId.String(),
		FencingToken: intent.FencingToken,
		ResourceKey:  intent.ResourceKey,
		ObservedAt:   timestamppb.Now(),
		Actor:        operationMonitorActor,
	})
}

// TODO: replace with a real expression parser if we ever need OR / NOT / globs.
func ruleMatches(rule string, transition map[string]string) bool {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return false
	}
	for _, pair := range strings.Split(rule, ",") {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) != 2 {
			return false
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if transition[k] != v {
			return false
		}
	}
	return true
}
