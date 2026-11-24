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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

func newOperationEvent(t *testing.T, key evt.EventId, msg proto.Message) *epb.Event {
	t.Helper()

	body, err := anypb.New(msg)
	require.NoError(t, err)

	return &epb.Event{
		RoutingKey: msgbus.PrepareRoute(orgName, evt.EventRoutingKey[key]),
		Msg:        body,
	}
}

func TestEventNotification_terminalEvents(t *testing.T) {
	events := []struct {
		name   string
		status db.OperationStatus
		audit  string
		event  func(t *testing.T, opId uuid.UUID) *epb.Event
	}{
		{"Completed", db.OperationSuccess, "completed",
			func(t *testing.T, opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, evt.EventOperationCompleted,
					&epb.OperationCompletedEvent{OperationId: opId.String(),
						FencingToken: 5})
			}},
		{"Failed", db.OperationFailed, "failed",
			func(t *testing.T, opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, evt.EventOperationFailed,
					&epb.OperationFailedEvent{OperationId: opId.String(), FencingToken: 5,
						Reason: "deadline exceeded"})
			}},
	}

	for _, e := range events {
		t.Run(e.name+"Applied", func(t *testing.T) {
			r := newTestRepos()
			op := runningOp(5)
			r.updates(op)

			resp, err := NewEventServer(orgName, r.repos()).
				EventNotification(context.Background(), e.event(t, op.Id))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, []uuid.UUID{op.Id}, r.released)
			assert.Equal(t, []string{e.audit}, r.events())
			r.assertExpectations(t)
		})

		t.Run(e.name+"LateIsAcknowledged", func(t *testing.T) {
			r := newTestRepos()
			op := runningOp(5)
			op.Status = db.OperationTimeout
			r.updates(op)

			resp, err := NewEventServer(orgName, r.repos()).
				EventNotification(context.Background(), e.event(t, op.Id))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, []string{"ignored_" + e.audit}, r.events())
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
			r.assertExpectations(t)
		})
	}
}

func TestEventNotification_actor(t *testing.T) {
	for name, tc := range map[string]struct {
		event func(opId uuid.UUID) *epb.Event
		actor string
	}{
		"CompletedCarriesActor": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, evt.EventOperationCompleted,
					&epb.OperationCompletedEvent{OperationId: opId.String(),
						FencingToken: 5, Actor: "operation-monitor"})
			},
			"operation-monitor",
		},
		"FailedCarriesActor": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, evt.EventOperationFailed,
					&epb.OperationFailedEvent{OperationId: opId.String(), FencingToken: 5,
						Reason: "deadline exceeded", Actor: "operation-monitor-sweeper"})
			},
			"operation-monitor-sweeper",
		},
		"MissingActorFallsBack": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, evt.EventOperationCompleted,
					&epb.OperationCompletedEvent{OperationId: opId.String(),
						FencingToken: 5})
			},
			pkg.OperationMonitorActor,
		},
		"BlankActorFallsBack": {
			func(opId uuid.UUID) *epb.Event {
				return newOperationEvent(t, evt.EventOperationFailed,
					&epb.OperationFailedEvent{OperationId: opId.String(), FencingToken: 5,
						Reason: "deadline exceeded", Actor: "  "})
			},
			pkg.OperationMonitorActor,
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
			r.assertExpectations(t)
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
			context.Background(), newOperationEvent(t, evt.EventOperationCompleted,
				&epb.OperationCompletedEvent{OperationId: "nope", FencingToken: 5}))

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

	assert.Equal(t, []string{"timeout"}, r.events())
	assert.Equal(t, pkg.OperationManagerSweeperActor, r.saved[0].Actor)
	r.assertExpectations(t)
}
