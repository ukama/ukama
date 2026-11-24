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
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

var ctx = context.Background()

func pendingOp() db.Operation {
	return db.Operation{
		Id:             uuid.NewV4(),
		Type:           "RestartNode",
		System:         "node",
		Status:         db.OperationPending,
		FencingToken:   3,
		ResourceKey:    "node:abc",
		LeaseExpiresAt: time.Now().Add(time.Minute),
	}
}

func TestMarkRunningTransition(t *testing.T) {
	t.Run("TransitionsPendingToRunning", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		r.updates(op)

		out, err := r.repos().markRunning(ctx, op.Id, op.FencingToken, "controller")

		assert.NoError(t, err)
		assert.Equal(t, db.OperationRunning, out.Status)
		assert.NotNil(t, out.StartedAt)
		assert.Equal(t, []string{"running"}, r.events())
		assert.Equal(t, "controller", r.saved[0].Actor)
		assert.Equal(t, op.Id, r.saved[0].OperationId)
		assert.Equal(t, op.ResourceKey, r.saved[0].ResourceKey)
		r.audits.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("RepeatIsAuditedAndSucceeds", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationRunning
		r.updates(op)

		out, err := r.repos().markRunning(ctx, op.Id, op.FencingToken, "controller")

		assert.NoError(t, err)
		assert.Equal(t, db.OperationRunning, out.Status)
		assert.Equal(t, []string{"ignored_running"}, r.events())
		assert.Empty(t, r.saved[0].Reason)
		r.audits.AssertNotCalled(t, "AddInTx", mock.Anything, mock.Anything)
	})

	for name, change := range map[string]func(op *db.Operation) uint64{
		"RejectsStaleFencingToken": func(op *db.Operation) uint64 {
			return op.FencingToken - 1
		},
		"RejectsFinishedOperation": func(op *db.Operation) uint64 {
			op.Status = db.OperationTimeout
			return op.FencingToken
		},
		"RejectsExpiredLease": func(op *db.Operation) uint64 {
			op.LeaseExpiresAt = time.Now().Add(-time.Second)
			return op.FencingToken
		},
		"RejectsExpiredLeaseWhileRunning": func(op *db.Operation) uint64 {
			op.Status = db.OperationRunning
			op.LeaseExpiresAt = time.Now().Add(-time.Second)
			return op.FencingToken
		},
	} {
		t.Run(name, func(t *testing.T) {
			r := newTestRepos()
			op := pendingOp()
			token := change(&op)
			r.updates(op)

			out, err := r.repos().markRunning(ctx, op.Id, token, "controller")

			assert.True(t, errors.Is(err, store.ErrPrecondition), err)
			assert.Nil(t, out)
			assert.Equal(t, []string{"rejected_running"}, r.events())
			assert.Equal(t, err.Error(), r.saved[0].Reason)
			assert.Equal(t, "controller", r.saved[0].Actor)
			r.audits.AssertNotCalled(t, "AddInTx", mock.Anything, mock.Anything)
		})
	}

	t.Run("RepoErrorIsReturnedWithoutAudit", func(t *testing.T) {
		r := newTestRepos()
		id := uuid.NewV4()
		r.updateFails(id, fmt.Errorf("operation %s: %w", id, store.ErrNotFound))

		out, err := r.repos().markRunning(ctx, id, 3, "controller")

		assert.True(t, errors.Is(err, store.ErrNotFound), err)
		assert.Nil(t, out)
		r.audits.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("RefusalAuditFailureIsReturned", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationRunning
		r.ops.On("Update", mock.Anything, op.Id, mock.Anything).
			Return(nil, ignore(&op, nil)).Once()
		r.audits.On("Add", mock.Anything, mock.Anything).
			Return(errors.New("connection reset")).Once()

		out, err := r.repos().markRunning(ctx, op.Id, op.FencingToken, "controller")

		assert.EqualError(t, err, "connection reset")
		assert.Nil(t, out)
		r.assertExpectations(t)
	})
}

func TestTerminateTransition(t *testing.T) {
	completed := db.OperationAudit{Event: "completed", Actor: "controller",
		Reason: "dispatched"}

	t.Run("FinishesReleasesAndAudits", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationRunning
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			completed, "")

		assert.NoError(t, err)
		assert.Equal(t, db.OperationSuccess, out.Status)
		assert.NotNil(t, out.TerminalAt)
		assert.Empty(t, out.Error)
		assert.Equal(t, []uuid.UUID{op.Id}, r.released)
		assert.Equal(t, []string{"completed"}, r.events())
		assert.Equal(t, op.Id, r.saved[0].OperationId)
		assert.Equal(t, "dispatched", r.saved[0].Reason)
		assert.Equal(t, *out.TerminalAt, r.saved[0].At)
	})

	t.Run("FailOnPendingRecordsTheError", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationFailed,
			db.OperationAudit{Event: "failed"}, "boom")

		assert.NoError(t, err)
		assert.Equal(t, db.OperationFailed, out.Status)
		assert.Equal(t, "boom", out.Error)
		assert.Equal(t, []string{"failed"}, r.events())
	})

	t.Run("RejectsStaleFencingToken", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationRunning
		op.FencingToken = 9
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, 4, db.OperationSuccess, completed, "")

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, out)
		assert.Equal(t, []string{"rejected_completed"}, r.events())
		assert.Equal(t, "dispatched", r.saved[0].Reason)
		assert.Empty(t, r.released)
	})

	t.Run("DuplicateWithSameStatusIsAuditedAndSucceeds", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationSuccess
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			completed, "")

		assert.NoError(t, err)
		assert.Equal(t, db.OperationSuccess, out.Status)
		assert.Equal(t, []string{"ignored_completed"}, r.events())
		assert.Empty(t, r.released)
	})

	t.Run("DifferentStatusIsAuditedAndRejected", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationTimeout
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			completed, "")

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, out)
		assert.Equal(t, []string{"ignored_completed"}, r.events())
	})

	t.Run("RejectsCompleteOnPending", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			completed, "")

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, out)
		assert.Equal(t, []string{"rejected_completed"}, r.events())
	})

	t.Run("ReleaseFailureIsReturnedWithoutAudit", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationRunning
		r.ops.On("Update", mock.Anything, op.Id, mock.Anything).
			Return(func(_ context.Context, _ uuid.UUID,
				nested func(*db.Operation, *gorm.DB) error) (*db.Operation, error) {
				current := op
				return nil, nested(&current, nil)
			}).Once()
		r.locks.On("Release", mock.Anything, mock.Anything).
			Return(errors.New("connection reset")).Once()

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			completed, "")

		assert.EqualError(t, err, "connection reset")
		assert.Nil(t, out)
		r.audits.AssertNotCalled(t, "AddInTx", mock.Anything, mock.Anything)
		r.audits.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("RejectsNonTerminalStatus", func(t *testing.T) {
		r := newTestRepos()

		out, err := r.repos().terminate(ctx, uuid.NewV4(), 1, db.OperationRunning,
			db.OperationAudit{Event: "running"}, "")

		assert.Error(t, err)
		assert.Nil(t, out)
		r.assertUntouched(t)
	})
}

func TestRefusalError(t *testing.T) {
	op := pendingOp()

	assert.Equal(t, "ignored_transition", ignore(&op, nil).Error())
	assert.Equal(t, "boom", reject(&op, errors.New("boom")).Error())
}
