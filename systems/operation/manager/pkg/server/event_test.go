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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	completedRoute = "event.cloud.global.{{ .Org}}.operation.manager.operation.completed"
	failedRoute    = "event.cloud.global.{{ .Org}}.operation.manager.operation.failed"
	timeoutRoute   = "event.cloud.global.{{ .Org}}.operation.manager.operation.timeout"
	cancelledRoute = "event.cloud.global.{{ .Org}}.operation.manager.operation.cancelled"
	satisfiedRoute = "event.cloud.global.{{ .Org}}.node.operationmonitor.intent.satisfied"
	expiredRoute   = "event.cloud.global.{{ .Org}}.node.operationmonitor.intent.expired"
)

const monitorDeadlineReason = "deadline exceeded"

func newOperationEvent(t *testing.T, route string, msg proto.Message) *epb.Event {
	t.Helper()

	body, err := anypb.New(msg)
	require.NoError(t, err)

	return &epb.Event{
		RoutingKey: msgbus.PrepareRoute(orgName, route),
		Msg:        body,
	}
}

func TestEventNotification_terminalEvents(t *testing.T) {
	observed := time.Now().UTC().Add(-time.Minute).Truncate(time.Second)
	events := []struct {
		name   string
		status db.OperationStatus
		audit  string
		reason string
		event  func(t *testing.T, opId uuid.UUID) *epb.Event
	}{
		{"Satisfied", db.OperationSuccess, pkg.AuditCompleted, "",
			func(t *testing.T, opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, satisfiedRoute,
					&epb.IntentSatisfiedEvent{OperationId: opId.String(), FencingToken: 5,
						ObservedAt: timestamppb.New(observed)})
			}},
		{"Expired", db.OperationTimeout, pkg.AuditTimeout, monitorDeadlineReason,
			func(t *testing.T, opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, expiredRoute,
					&epb.IntentExpiredEvent{OperationId: opId.String(), FencingToken: 5,
						Reason: monitorDeadlineReason, ExpiredAt: timestamppb.New(observed)})
			}},
	}

	for _, e := range events {
		t.Run(e.name+"Applied", func(t *testing.T) {
			r := newTestRepos()
			op := runningOp(5)
			op.CreatedAt = observed.Add(-time.Hour)
			r.updates(op)

			resp, err := NewEventServer(orgName, r.repos()).
				EventNotification(context.Background(), e.event(t, op.Id))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, []uuid.UUID{op.Id}, r.released)
			assert.Equal(t, []string{e.audit}, r.events())
			assert.Equal(t, e.reason, r.saved[0].Reason)
			require.Len(t, r.published, 1)
			assert.Equal(t, e.status, r.published[0].Status)
			assert.Equal(t, observed, *r.published[0].TerminalAt)
			assert.Equal(t, observed, r.results[0].At)
			r.assertExpectations(t)
		})

		t.Run(e.name+"LateIsAcknowledged", func(t *testing.T) {
			r := newTestRepos()
			op := runningOp(5)
			op.Status = db.OperationCancelled
			r.updates(op)

			resp, err := NewEventServer(orgName, r.repos()).
				EventNotification(context.Background(), e.event(t, op.Id))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, []string{"ignored_" + e.audit}, r.events())
			assert.Empty(t, r.published)
			r.assertExpectations(t)
		})

		t.Run(e.name+"UnknownOperationIsAcknowledged", func(t *testing.T) {
			r := newTestRepos()
			opId := uuid.NewV4()
			r.updateFails(opId, fmt.Errorf("operation %s: %w", opId, store.ErrNotFound))

			resp, err := NewEventServer(orgName, r.repos()).
				EventNotification(context.Background(), e.event(t, opId))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Empty(t, r.published)
			r.assertExpectations(t)
		})

		t.Run(e.name+"RepoFailureIsUnavailable", func(t *testing.T) {
			r := newTestRepos()
			opId := uuid.NewV4()
			r.updateFails(opId, assert.AnError)

			resp, err := NewEventServer(orgName, r.repos()).
				EventNotification(context.Background(), e.event(t, opId))

			assert.Equal(t, codes.Unavailable, status.Code(err))
			assert.Nil(t, resp)
			assert.Empty(t, r.published)
			r.assertExpectations(t)
		})
	}
}

func TestEventNotification_actor(t *testing.T) {
	for name, tc := range map[string]struct {
		event func(opId uuid.UUID) *epb.Event
		actor string
	}{
		"SatisfiedCarriesActor": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, satisfiedRoute,
					&epb.IntentSatisfiedEvent{OperationId: opId.String(),
						FencingToken: 5, Actor: operationMonitorActor})
			},
			operationMonitorActor,
		},
		"ExpiredCarriesActor": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, expiredRoute,
					&epb.IntentExpiredEvent{OperationId: opId.String(), FencingToken: 5,
						Reason: monitorDeadlineReason, Actor: operationMonitorSweeperActor})
			},
			operationMonitorSweeperActor,
		},
		"MissingActorFallsBack": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, satisfiedRoute,
					&epb.IntentSatisfiedEvent{OperationId: opId.String(),
						FencingToken: 5})
			},
			operationMonitorActor,
		},
		"BlankActorFallsBack": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, expiredRoute,
					&epb.IntentExpiredEvent{OperationId: opId.String(), FencingToken: 5,
						Reason: monitorDeadlineReason, Actor: "  "})
			},
			operationMonitorActor,
		},
	} {
		t.Run(name, func(t *testing.T) {
			r := newTestRepos()
			op := runningOp(5)
			r.updates(op)

			_, err := NewEventServer(orgName, r.repos()).
				EventNotification(context.Background(), tc.event(op.Id))

			assert.NoError(t, err)
			assert.Len(t, r.saved, 1)
			assert.Equal(t, tc.actor, r.saved[0].Actor)
			assert.Equal(t, tc.actor, r.results[0].Actor)
			r.assertExpectations(t)
		})
	}
}

func TestEventNotification_badPayload(t *testing.T) {
	for name, route := range map[string]string{
		"Satisfied": satisfiedRoute,
		"Expired":   expiredRoute,
	} {
		t.Run(name, func(t *testing.T) {
			r := newTestRepos()
			body, err := anypb.New(&epb.EventResponse{})
			require.NoError(t, err)

			_, err = NewEventServer(orgName, r.repos()).EventNotification(
				context.Background(), &epb.Event{
					RoutingKey: msgbus.PrepareRoute(orgName, route),
					Msg:        body,
				})

			assert.Error(t, err)
			r.assertUntouched(t)
		})
	}
}

func TestEventNotification_invalidInput(t *testing.T) {
	t.Run("NilEvent", func(t *testing.T) {
		r := newTestRepos()

		_, err := NewEventServer(orgName, r.repos()).
			EventNotification(context.Background(), nil)

		assert.Error(t, err)
		r.assertUntouched(t)
	})

	t.Run("InvalidOperationId", func(t *testing.T) {
		r := newTestRepos()

		_, err := NewEventServer(orgName, r.repos()).EventNotification(
			context.Background(), newOperationEvent(t, satisfiedRoute,
				&epb.IntentSatisfiedEvent{OperationId: "nope", FencingToken: 5}))

		assert.Error(t, err)
		r.assertUntouched(t)
	})

	t.Run("UnknownRoutingKeyIsAcknowledged", func(t *testing.T) {
		r := newTestRepos()

		resp, err := NewEventServer(orgName, r.repos()).EventNotification(
			context.Background(), &epb.Event{RoutingKey: "event.cloud.local.other"})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		r.assertUntouched(t)
	})
}

func TestSweeperActor(t *testing.T) {
	r := newTestRepos()
	op := runningOp(4)
	r.ops.On("List", mock.Anything, mock.MatchedBy(func(f db.OperationFilter) bool {
		return f.Limit == 100
	})).Return([]*db.Operation{&op}, nil).Once()
	r.updates(op)

	NewSweeper(r.repos()).sweepOnce(context.Background())

	assert.Equal(t, []string{pkg.AuditTimeout}, r.events())
	assert.Equal(t, operationManagerSweeperActor, r.saved[0].Actor)
	r.assertExpectations(t)
}
