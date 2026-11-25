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
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/operation-monitor/mocks"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
)

func newMonitorServer(repo db.IntentRepo) *MonitorServer {
	return NewMonitorServer(testOrg, "org-id", repo, &mbmocks.MsgBusServiceClient{})
}

func registerRequest() *pb.RegisterIntentRequest {
	return &pb.RegisterIntentRequest{
		OperationId:  uuid.NewV4().String(),
		ResourceKey:  nodeResourcePrefix + testNodeId,
		ActionType:   "RestartNode",
		FencingToken: 1,
	}
}

func TestRegisterIntent_Duplicate(t *testing.T) {
	repo := &mocks.IntentRepo{}
	dup := &pgconn.PgError{Code: sql.PGERROR_CODE_UNIQUE_VIOLATION}
	repo.On("Add", mock.Anything).Return(dup).Once()

	_, err := newMonitorServer(repo).RegisterIntent(context.TODO(), registerRequest())

	assert.Equal(t, codes.AlreadyExists, status.Code(err))
	repo.AssertExpectations(t)
}

func TestRegisterIntent_StoreFailure(t *testing.T) {
	repo := &mocks.IntentRepo{}
	repo.On("Add", mock.Anything).Return(errors.New("database unavailable")).Once()

	_, err := newMonitorServer(repo).RegisterIntent(context.TODO(), registerRequest())

	assert.Equal(t, codes.Internal, status.Code(err))
	repo.AssertExpectations(t)
}

func TestRegisterIntent_Success(t *testing.T) {
	repo := &mocks.IntentRepo{}
	repo.On("Add", mock.MatchedBy(func(i *db.Intent) bool {
		return i.CompletionRule == pkg.RuleSubstateOn && i.Status == db.IntentWatching
	})).Return(nil).Once()

	resp, err := newMonitorServer(repo).RegisterIntent(context.TODO(), registerRequest())

	assert.NoError(t, err)
	assert.Equal(t, pkg.RuleSubstateOn, resp.Intent.CompletionRule)
	repo.AssertExpectations(t)
}

func TestCancelIntent_Unknown(t *testing.T) {
	repo := &mocks.IntentRepo{}
	opId := uuid.NewV4()
	repo.On("Update", opId, closeAs(db.IntentCancelled)).
		Return(nil, gorm.ErrRecordNotFound).Once()

	_, err := newMonitorServer(repo).CancelIntent(context.TODO(),
		&pb.CancelIntentRequest{OperationId: opId.String()})

	assert.Equal(t, codes.NotFound, status.Code(err))
	repo.AssertExpectations(t)
}

func TestCancelIntent_StoreFailure(t *testing.T) {
	repo := &mocks.IntentRepo{}
	opId := uuid.NewV4()
	repo.On("Update", opId, closeAs(db.IntentCancelled)).
		Return(nil, errors.New("database unavailable")).Once()

	_, err := newMonitorServer(repo).CancelIntent(context.TODO(),
		&pb.CancelIntentRequest{OperationId: opId.String()})

	assert.Equal(t, codes.Internal, status.Code(err))
	repo.AssertExpectations(t)
}

func TestCancelIntent_Success(t *testing.T) {
	repo := &mocks.IntentRepo{}
	opId := uuid.NewV4()
	repo.On("Update", opId, closeAs(db.IntentCancelled)).
		Return(&db.Intent{OperationId: opId}, nil).Once()

	_, err := newMonitorServer(repo).CancelIntent(context.TODO(),
		&pb.CancelIntentRequest{OperationId: opId.String()})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSweeper_PublishFailureCanRetry(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)
	intent.Deadline = time.Now().Add(-time.Minute)
	s := NewSweeper(NewMonitorServer(testOrg, "org-id", repo, mb))

	repo.On("List", expiredFilter(100)).
		Return([]db.Intent{intent}, nil).Twice()
	mb.On("PublishRequest", mock.Anything, mock.Anything).
		Return(errors.New("message bus unavailable")).Once()

	s.sweepOnce()
	repo.AssertNotCalled(t, "Update", mock.Anything, closing)

	mb.On("PublishRequest",
		msgbus.PrepareRoute(testOrg, expiredRoute),
		mock.MatchedBy(func(event *epb.IntentExpiredEvent) bool {
			return event.OperationId == intent.OperationId.String() &&
				event.FencingToken == intent.FencingToken &&
				event.Reason == reasonDeadlineExceeded && event.ExpiredAt != nil &&
				event.Actor == operationMonitorSweeperActor
		})).Return(nil).Once()
	repo.On("Update", intent.OperationId, closeAs(db.IntentExpired)).
		Return(&intent, nil).Once()

	s.sweepOnce()
	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}

func TestSweeper_TerminalWriteFailureCanRetry(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)
	intent.Deadline = time.Now().Add(-time.Minute)
	s := NewSweeper(NewMonitorServer(testOrg, "org-id", repo, mb))

	repo.On("List", expiredFilter(100)).
		Return([]db.Intent{intent}, nil).Twice()
	mb.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()
	repo.On("Update", intent.OperationId, closeAs(db.IntentExpired)).
		Run(func(_ mock.Arguments) {
			mb.AssertNumberOfCalls(t, "PublishRequest", 1)
		}).Return(nil, errors.New("database unavailable")).Once()

	s.sweepOnce()
	mb.AssertNumberOfCalls(t, "PublishRequest", 1)

	repo.On("Update", intent.OperationId, closeAs(db.IntentExpired)).
		Return(&intent, nil).Once()
	s.sweepOnce()

	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}

var closing = mock.MatchedBy(func(u db.IntentUpdate) bool { return u.Status != nil })

func expiredFilter(limit int) any {
	return mock.MatchedBy(func(f db.IntentFilter) bool {
		return f.Status != nil && *f.Status == db.IntentWatching &&
			!f.DeadlineBefore.IsZero() && f.Limit == limit && f.ResourceKey == ""
	})
}

func TestIntentQueryHelpers(t *testing.T) {
	now := time.Now()

	w := watchingOn("node:abc")
	assert.Equal(t, "node:abc", w.ResourceKey)
	assert.Equal(t, db.IntentWatching, *w.Status)
	assert.True(t, w.DeadlineBefore.IsZero())

	e := expiredBefore(now, 5)
	assert.Equal(t, db.IntentWatching, *e.Status)
	assert.Equal(t, now, e.DeadlineBefore)
	assert.Equal(t, 5, e.Limit)

	assert.Equal(t, db.IntentExpired, *closeAs(db.IntentExpired).Status)
	assert.Nil(t, closeAs(db.IntentExpired).Armed)
	assert.True(t, *arm().Armed)
	assert.Nil(t, arm().Status)
}
