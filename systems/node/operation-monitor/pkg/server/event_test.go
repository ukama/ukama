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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/operation-monitor/mocks"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"

	evt "github.com/ukama/ukama/systems/common/events"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

const (
	testOrg    = "test-org"
	testNodeId = "uk-983794-tnode-78-7830"
)

var (
	stateConfiguring = upb.NodeState_Configuring.String()
	stateOperational = upb.NodeState_Operational.String()
	stateFaulty      = upb.NodeState_Faulty.String()
)

func transitionEvent(t *testing.T, state, substate string) *epb.Event {
	t.Helper()
	return transitionEventFor(t, testNodeId, state, substate)
}

func transitionEventFor(t *testing.T, nodeId, state, substate string) *epb.Event {
	t.Helper()
	msg, err := anypb.New(&epb.NodeStateChangeEvent{
		NodeId:   nodeId,
		State:    state,
		Substate: substate,
	})
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	route := msgbus.PrepareRoute(testOrg, evt.EventRoutingKey[evt.EventNodeStateTransition])
	return &epb.Event{RoutingKey: route, Msg: msg}
}

func watchingIntent(armed bool) db.Intent {
	return db.Intent{
		Id:             uuid.NewV4(),
		OperationId:    uuid.NewV4(),
		ResourceKey:    nodeResourcePrefix + testNodeId,
		ActionType:     "RestartNode",
		FencingToken:   1,
		CompletionRule: pkg.RuleSubstateOn,
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

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateConfiguring,
		pkg.SubstateOn))

	assert.NoError(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)
	repo.AssertNotCalled(t, "Update", mock.Anything, arm())
	mb.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
}

// A transition outside the target state arms the intent.
func TestHandleStateTransition_DepartureArmsIntent(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(false)

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()
	repo.On("Update", intent.OperationId, arm()).Return(nil, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateConfiguring,
		substateOff))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)
}

// An armed intent completes on the next rule match and publishes the
// operation-completed event that releases the manager lock.
func TestHandleStateTransition_ArmedIntentCompletes(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()
	repo.On("Update", intent.OperationId, closeAs(db.IntentCompleted)).
		Return(&intent, nil).Once()
	mb.On("PublishRequest",
		msgbus.PrepareRoute(testOrg, evt.EventRoutingKey[evt.EventIntentSatisfied]),
		mock.MatchedBy(func(event *epb.IntentSatisfiedEvent) bool {
			return event.Actor == operationMonitorActor && event.ObservedAt != nil
		})).Return(nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateConfiguring,
		pkg.SubstateOn))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}

func TestHandleStateTransition_PublishFailureCanRetry(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)
	publishErr := errors.New("message bus unavailable")

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Twice()
	mb.On("PublishRequest", mock.Anything, mock.Anything).
		Return(publishErr).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateOperational,
		pkg.SubstateOn))

	assertRetryable(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)

	mb.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("Update", intent.OperationId, closeAs(db.IntentCompleted)).
		Return(&intent, nil).Once()

	_, err = s.EventNotification(context.TODO(), transitionEvent(t, stateOperational,
		pkg.SubstateOn))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}

func TestHandleStateTransition_TerminalWriteFailureCanRetry(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)
	writeErr := errors.New("database unavailable")

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Twice()
	mb.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()
	repo.On("Update", intent.OperationId, closeAs(db.IntentCompleted)).
		Run(func(_ mock.Arguments) {
			mb.AssertNumberOfCalls(t, "PublishRequest", 1)
		}).Return(nil, writeErr).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateOperational,
		pkg.SubstateOn))
	assertRetryable(t, err)

	repo.On("Update", intent.OperationId, closeAs(db.IntentCompleted)).
		Return(&intent, nil).Once()
	_, err = s.EventNotification(context.TODO(), transitionEvent(t, stateOperational,
		pkg.SubstateOn))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}

func TestHandleStateTransition_ArmFailureCanRetry(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(false)
	writeErr := errors.New("database unavailable")

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Twice()
	repo.On("Update", intent.OperationId, arm()).Return(nil, writeErr).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateOperational,
		substateOff))
	assertRetryable(t, err)

	repo.On("Update", intent.OperationId, arm()).Return(nil, nil).Once()
	_, err = s.EventNotification(context.TODO(), transitionEvent(t, stateOperational,
		substateOff))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	mb.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
}

func assertRetryable(t *testing.T, err error) {
	t.Helper()
	assert.Equal(t, codes.Unavailable, status.Code(err))
}

func TestHandleStateTransition_LookupFailureIsRetryable(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return(nil, errors.New("database unavailable")).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateOperational,
		pkg.SubstateOn))

	assertRetryable(t, err)
	repo.AssertExpectations(t)
}

func TestHandleStateTransition_EmptySubstateDoesNotArm(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(false)

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateConfiguring, ""))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update", mock.Anything, arm())
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)
}

func TestHandleStateTransition_UppercaseNodeIdMatchesIntent(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(false)

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()
	repo.On("Update", intent.OperationId, arm()).Return(nil, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(),
		transitionEventFor(t, strings.ToUpper(testNodeId), stateConfiguring,
			substateOff))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestHandleAppsChanged_LookupFailureIsRetryable(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return(nil, errors.New("database unavailable")).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(),
		appsChangedEvent(t, map[string]string{"example": "1.2.4"}))

	assertRetryable(t, err)
	repo.AssertExpectations(t)
}

func rule(kv ...string) string {
	pairs := make([]string, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		pairs = append(pairs, kv[i]+"="+kv[i+1])
	}
	return strings.Join(pairs, ",")
}

func TestRuleObserved(t *testing.T) {
	stateRule := rule(ruleKeyState, stateConfiguring, pkg.RuleKeySubstate,
		pkg.SubstateOn)

	assert.True(t, ruleObserved(pkg.RuleSubstateOn,
		map[string]string{pkg.RuleKeySubstate: substateOff}))
	assert.True(t, ruleObserved(stateRule,
		map[string]string{ruleKeyState: stateFaulty, pkg.RuleKeySubstate: substateOff}))
	assert.False(t, ruleObserved(pkg.RuleSubstateOn,
		map[string]string{pkg.RuleKeySubstate: ""}))
	assert.False(t, ruleObserved(pkg.RuleSubstateOn,
		map[string]string{ruleKeyState: stateConfiguring}))
	assert.False(t, ruleObserved(stateRule,
		map[string]string{ruleKeyState: stateConfiguring, pkg.RuleKeySubstate: ""}))
}

func TestRuleMatches(t *testing.T) {
	tr := map[string]string{ruleKeyState: stateConfiguring,
		pkg.RuleKeySubstate: pkg.SubstateOn}

	assert.True(t, ruleMatches(pkg.RuleSubstateOn, tr))
	assert.True(t, ruleMatches(rule(ruleKeyState, stateConfiguring, pkg.RuleKeySubstate,
		pkg.SubstateOn), tr))
	assert.False(t, ruleMatches(rule(ruleKeyState, stateOperational), tr))
	assert.False(t, ruleMatches("", tr))
	assert.False(t, ruleMatches("garbage", tr))
}

func appsChangedEvent(t *testing.T, apps map[string]string) *epb.Event {
	t.Helper()
	ev := &epb.HealthAppsChangedEvent{NodeId: testNodeId, NodeType: "tower"}
	for name, version := range apps {
		ev.Apps = append(ev.Apps, &epb.HealthAppVersion{Name: name, Version: version})
	}
	msg, err := anypb.New(ev)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	route := msgbus.PrepareRoute(testOrg,
		"event.cloud.local.{{ .Org}}.node.health.apps.changed")
	return &epb.Event{RoutingKey: route, Msg: msg}
}

func updateIntent(rule string) db.Intent {
	i := watchingIntent(false)
	i.ActionType = "UpdateSoftware"
	i.CompletionRule = rule
	return i
}

func TestHandleAppsChanged_TargetVersionCompletesWithoutArming(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := updateIntent(appRulePrefix + "example=1.2.4")

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()
	repo.On("Update", intent.OperationId, closeAs(db.IntentCompleted)).
		Return(&intent, nil).Once()
	mb.On("PublishRequest",
		msgbus.PrepareRoute(testOrg, evt.EventRoutingKey[evt.EventIntentSatisfied]),
		mock.MatchedBy(func(event *epb.IntentSatisfiedEvent) bool {
			return event.Actor == operationMonitorActor &&
				event.FencingToken == intent.FencingToken
		})).Return(nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(),
		appsChangedEvent(t, map[string]string{"example": "1.2.4", "other": "0.1"}))

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update", mock.Anything, arm())
}

func TestHandleAppsChanged_OtherVersionKeepsWatching(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := updateIntent(appRulePrefix + "example=1.2.4")

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(),
		appsChangedEvent(t, map[string]string{"example": "1.2.3"}))

	assert.NoError(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)
	repo.AssertNotCalled(t, "Update", mock.Anything, arm())
	mb.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
}

func TestHandleAppsChanged_IgnoresStateRuleIntent(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(false)

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(),
		appsChangedEvent(t, map[string]string{"example": "1.2.4"}))

	assert.NoError(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, arm())
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)
}

func TestHandleStateTransition_IgnoresAppRuleIntent(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := updateIntent(appRulePrefix + "example=1.2.4")

	repo.On("List", watchingOn(nodeResourcePrefix+testNodeId)).
		Return([]db.Intent{intent}, nil).Once()

	s := newEventServer(repo, mb)
	_, err := s.EventNotification(context.TODO(), transitionEvent(t, stateConfiguring,
		substateOff))

	assert.NoError(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, arm())
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)
}

func resultEvent(t *testing.T, key evt.EventId, msg proto.Message) *epb.Event {
	t.Helper()
	body, err := anypb.New(msg)
	require.NoError(t, err)
	route := msgbus.PrepareRoute(testOrg, evt.EventRoutingKey[key])
	return &epb.Event{RoutingKey: route, Msg: body}
}

func TestHandleResult_ClosesIntent(t *testing.T) {
	opId := uuid.NewV4()
	for name, tc := range map[string]struct {
		event  *epb.Event
		status db.IntentStatus
	}{
		"Completed": {resultEvent(t, evt.EventOperationCompleted,
			&epb.OperationCompletedEvent{OperationId: opId.String()}), db.IntentCompleted},
		"Failed": {resultEvent(t, evt.EventOperationFailed,
			&epb.OperationFailedEvent{OperationId: opId.String()}), db.IntentFailed},
		"TimedOut": {resultEvent(t, evt.EventOperationTimedOut,
			&epb.OperationTimedOutEvent{OperationId: opId.String()}), db.IntentExpired},
		"Cancelled": {resultEvent(t, evt.EventOperationCancelled,
			&epb.OperationCancelledEvent{OperationId: opId.String()}), db.IntentCancelled},
	} {
		t.Run(name, func(t *testing.T) {
			repo := &mocks.IntentRepo{}
			mb := &mbmocks.MsgBusServiceClient{}
			closed := &db.Intent{OperationId: opId, Status: tc.status}
			repo.On("Update", opId, closeAs(tc.status)).Return(closed, nil).Once()

			resp, err := newEventServer(repo, mb).EventNotification(context.TODO(),
				tc.event)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			repo.AssertExpectations(t)
			mb.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
		})
	}
}

func TestHandleResult_Errors(t *testing.T) {
	opId := uuid.NewV4()
	failed := func() *epb.Event {
		return resultEvent(t, evt.EventOperationFailed,
			&epb.OperationFailedEvent{OperationId: opId.String()})
	}

	t.Run("NoIntentIsAcknowledged", func(t *testing.T) {
		repo := &mocks.IntentRepo{}
		repo.On("Update", opId, closeAs(db.IntentFailed)).
			Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := newEventServer(repo, &mbmocks.MsgBusServiceClient{}).
			EventNotification(context.TODO(), failed())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("StoreFailureIsRetryable", func(t *testing.T) {
		repo := &mocks.IntentRepo{}
		repo.On("Update", opId, closeAs(db.IntentFailed)).
			Return(nil, errors.New("db down")).Once()

		_, err := newEventServer(repo, &mbmocks.MsgBusServiceClient{}).
			EventNotification(context.TODO(), failed())

		assertRetryable(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("InvalidOperationIdIsRejected", func(t *testing.T) {
		repo := &mocks.IntentRepo{}

		_, err := newEventServer(repo, &mbmocks.MsgBusServiceClient{}).EventNotification(
			context.TODO(), resultEvent(t, evt.EventOperationCancelled,
				&epb.OperationCancelledEvent{OperationId: "nope"}))

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update", mock.Anything, closing)
	})

	t.Run("BadPayloadIsRejected", func(t *testing.T) {
		repo := &mocks.IntentRepo{}

		_, err := newEventServer(repo, &mbmocks.MsgBusServiceClient{}).EventNotification(
			context.TODO(), resultEvent(t, evt.EventOperationTimedOut,
				&epb.EventResponse{}))

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update", mock.Anything, closing)
	})
}
