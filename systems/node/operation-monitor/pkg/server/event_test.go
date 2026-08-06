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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/anypb"

	evt "github.com/ukama/ukama/systems/common/events"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/node/operation-monitor/mocks"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"
)

const (
	testOrg    = "test-org"
	testNodeId = "uk-983794-tnode-78-7830"
)

func transitionEvent(t *testing.T, state, substate string) *epb.Event {
	t.Helper()
	msg, err := anypb.New(&epb.NodeStateChangeEvent{
		NodeId:   testNodeId,
		State:    state,
		Substate: substate,
	})
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	return &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrg, evt.EventRoutingKey[evt.EventNodeStateTransition]),
		Msg:        msg,
	}
}

func watchingIntent(armed bool) db.MonitoredIntent {
	return db.MonitoredIntent{
		Id:             uuid.NewV4(),
		OperationId:    uuid.NewV4(),
		ResourceKey:    "node:" + testNodeId,
		ActionType:     "RestartNode",
		FencingToken:   1,
		CompletionRule: "substate=on",
		Status:         db.IntentWatching,
		Armed:          armed,
	}
}

func newEventServer(repo db.IntentRepo, mb *mbmocks.MsgBusServiceClient) *EventServer {
	return NewEventServer(testOrg, NewMonitorServer(testOrg, "org-id", repo, mb))
}

// A steady-state report matching the rule must NOT complete a fresh intent:
// the node never departed, so the action can't have run yet.
func TestHandleStateTransition_SteadyStateDoesNotComplete(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(false)

	repo.On("FindWatchingByResource", "node:"+testNodeId).
		Return([]db.MonitoredIntent{intent}, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, "Configured", "on"))

	assert.NoError(t, err)
	repo.AssertNotCalled(t, "MarkTerminal", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "Arm", mock.Anything)
	mb.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
}

// A transition outside the target state arms the intent.
func TestHandleStateTransition_DepartureArmsIntent(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(false)

	repo.On("FindWatchingByResource", "node:"+testNodeId).
		Return([]db.MonitoredIntent{intent}, nil).Once()
	repo.On("Arm", intent.OperationId).Return(nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, "Configured", "reboot"))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "MarkTerminal", mock.Anything, mock.Anything)
}

// An armed intent completes on the next rule match and publishes the
// operation-completed event that releases the manager lock.
func TestHandleStateTransition_ArmedIntentCompletes(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)

	repo.On("FindWatchingByResource", "node:"+testNodeId).
		Return([]db.MonitoredIntent{intent}, nil).Once()
	repo.On("MarkTerminal", intent.OperationId, db.IntentCompleted).
		Return(&intent, nil).Once()
	mb.On("PublishRequest", mock.MatchedBy(func(route string) bool {
		return strings.HasPrefix(route, "event.cloud.global.") &&
			strings.HasSuffix(route, ".operation.manager.operation.completed")
	}), mock.Anything).Return(nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, "Configured", "on"))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}

func TestRuleMatches(t *testing.T) {
	tr := map[string]string{"state": "Configured", "substate": "on"}

	assert.True(t, ruleMatches("substate=on", tr))
	assert.True(t, ruleMatches("state=Configured,substate=on", tr))
	assert.False(t, ruleMatches("state=Ready", tr))
	assert.False(t, ruleMatches("", tr))
	assert.False(t, ruleMatches("garbage", tr))
}
