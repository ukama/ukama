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
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

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
	default:
		log.Warnf("operation-monitor: unhandled routing key %s", event.RoutingKey)
		return &epb.EventResponse{}, nil
	}
}

func (e *EventServer) handleStateTransition(_ context.Context, event *epb.Event) (*epb.EventResponse, error) {
	msg, err := epb.UnmarshalNodeStateChangeEvent(event.Msg, "NodeStateChangeEvent")
	if err != nil {
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
		if !ruleMatches(intent.CompletionRule, transition) {
			continue
		}
		if _, err := e.monitor.repo.MarkTerminal(intent.OperationId, db.IntentCompleted); err != nil {
			log.Errorf("operation-monitor: mark %s completed: %v", intent.OperationId, err)
			continue
		}
		if err := e.publishCompleted(intent); err != nil {
			log.Errorf("operation-monitor: publish completed for %s: %v", intent.OperationId, err)
			continue
		}
		log.Infof("operation-monitor: intent %s satisfied (rule=%q matched %v)",
			intent.OperationId, intent.CompletionRule, transition)
	}
	return &epb.EventResponse{}, nil
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

// ruleMatches evaluates a key=value,key2=value2 rule against the transition map.
// All k=v pairs must match. Empty rule matches anything (defensive — RegisterIntent
// rejects empty rules so this should never happen in practice).
//
// TODO: replace with a real expression parser if we ever need OR / NOT / globs.
func ruleMatches(rule string, transition map[string]string) bool {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return true
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
