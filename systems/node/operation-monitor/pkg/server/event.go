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
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"
)

type EventServer struct {
	epb.UnimplementedEventNotificationServiceServer
	orgName string
	monitor *MonitorServer
}

func NewEventServer(orgName string, m *MonitorServer) *EventServer {
	return &EventServer{orgName: orgName, monitor: m}
}

func (e *EventServer) EventNotification(ctx context.Context, event *epb.Event) (*epb.EventResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("nil event")
	}
	log.Infof("operation-monitor: received %s", event.RoutingKey)

	switch event.RoutingKey {
	case msgbus.PrepareRoute(e.orgName, evt.EventRoutingKey[evt.EventNodeStateTransition]):
		return e.handleStateTransition(ctx, event)
	case msgbus.PrepareRoute(e.orgName, "event.cloud.local.{{ .Org}}.node.health.apps.changed"):
		return e.handleAppsChanged(ctx, event)
	default:
		log.Warnf("operation-monitor: unhandled routing key %s", event.RoutingKey)
		return &epb.EventResponse{}, nil
	}
}

func (e *EventServer) handleStateTransition(_ context.Context, event *epb.Event) (*epb.EventResponse, error) {
	msg, err := cpb.UnmarshalProtoEvent[epb.NodeStateChangeEvent](event.Msg)
	if err != nil {
		log.Errorf("operation-monitor: unmarshal NodeStateChangeEvent: %v", err)
		return nil, err
	}

	resourceKey := "node:" + msg.NodeId
	intents, err := e.monitor.repo.FindWatchingByResource(resourceKey)
	if err != nil {
		log.Errorf("operation-monitor: lookup intents for %s: %v", resourceKey, err)
		return nil, err
	}
	if len(intents) == 0 {
		return &epb.EventResponse{}, nil
	}

	transition := map[string]string{
		"state":    msg.State,
		"substate": msg.Substate,
		"node_id":  msg.NodeId,
	}

	for i := range intents {
		intent := &intents[i]
		if isAppRule(intent.CompletionRule) {
			continue
		}
		if !ruleMatches(intent.CompletionRule, transition) {
			if !intent.Armed {
				if err := e.monitor.repo.Arm(intent.OperationId); err != nil {
					log.Errorf("operation-monitor: arm %s: %v", intent.OperationId, err)
					return nil, err
				} else {
					log.Infof("operation-monitor: intent %s armed by %v", intent.OperationId, transition)
				}
			}
			continue
		}
		if !intent.Armed {
			log.Infof("operation-monitor: intent %s matched %v but not armed yet (no departure observed), skipping",
				intent.OperationId, transition)
			continue
		}
		if err := e.complete(intent, transition); err != nil {
			return nil, err
		}
	}
	return &epb.EventResponse{}, nil
}

func (e *EventServer) handleAppsChanged(_ context.Context, event *epb.Event) (*epb.EventResponse, error) {
	msg, err := cpb.UnmarshalProtoEvent[epb.HealthAppsChangedEvent](event.Msg)
	if err != nil {
		log.Errorf("operation-monitor: unmarshal HealthAppsChangedEvent: %v", err)
		return nil, err
	}

	resourceKey := "node:" + msg.NodeId
	intents, err := e.monitor.repo.FindWatchingByResource(resourceKey)
	if err != nil {
		log.Errorf("operation-monitor: lookup intents for %s: %v", resourceKey, err)
		return nil, err
	}
	if len(intents) == 0 {
		return &epb.EventResponse{}, nil
	}

	inventory := make(map[string]string, len(msg.Apps))
	for _, app := range msg.Apps {
		inventory[pkg.AppRulePrefix+app.Name] = app.Version
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

func (e *EventServer) complete(intent *db.MonitoredIntent, observed map[string]string) error {
	if err := e.publishCompleted(intent); err != nil {
		log.Errorf("operation-monitor: publish completed for %s: %v", intent.OperationId, err)
		return err
	}
	if _, err := e.monitor.repo.MarkTerminal(intent.OperationId, db.IntentCompleted); err != nil {
		log.Errorf("operation-monitor: mark %s completed: %v", intent.OperationId, err)
		return err
	}
	log.Infof("operation-monitor: intent %s satisfied (rule=%q matched %v)",
		intent.OperationId, intent.CompletionRule, observed)
	return nil
}

func isAppRule(rule string) bool {
	return strings.HasPrefix(strings.TrimSpace(rule), pkg.AppRulePrefix)
}

func (e *EventServer) publishCompleted(intent *db.MonitoredIntent) error {
	route := e.monitor.publishBuilder.SetAction("completed").SetObject("operation").MustBuild()
	return e.monitor.msgbus.PublishRequest(route, &epb.OperationCompletedEvent{
		OperationId:  intent.OperationId.String(),
		FencingToken: intent.FencingToken,
		ResourceKey:  intent.ResourceKey,
		CompletedAt:  timestamppb.Now(),
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
