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

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

func newTestSweeper(r *testRepos, interval time.Duration) *Sweeper {
	return &Sweeper{repos: r.repos(), interval: interval, batch: 10}
}

func expiredFilter() interface{} {
	return mock.MatchedBy(func(f db.OperationFilter) bool {
		return len(f.Statuses) == 2 && f.Statuses[0] == db.OperationPending &&
			f.Statuses[1] == db.OperationRunning && f.LeaseExpiresBefore != nil &&
			!f.LeaseExpiresBefore.After(time.Now()) &&
			f.OrderBy == db.OrderByLeaseExpiry && f.Limit == 10
	})
}

func TestSweeperSweepOnce(t *testing.T) {
	t.Run("TimesOutEveryExpiredOperationOldestFirst", func(t *testing.T) {
		r := newTestRepos()
		first := runningOp(3)
		second := runningOp(4)

		r.ops.On("List", mock.Anything, expiredFilter()).
			Return([]*db.Operation{&first, &second}, nil).Once()
		r.updateFails(first.Id, fmt.Errorf("db down"))
		r.updates(second)

		newTestSweeper(r, time.Minute).sweepOnce(context.Background())

		assert.Equal(t, []uuid.UUID{second.Id}, r.released)
		assert.Equal(t, []string{"timeout"}, r.events())
		assert.Equal(t, "lease expired", r.saved[0].Reason)
		assert.Equal(t, pkg.OperationManagerSweeperActor, r.saved[0].Actor)
		r.assertExpectations(t)
	})

	t.Run("BoundsEachSweepByTheInterval", func(t *testing.T) {
		r := newTestRepos()
		r.ops.On("List", mock.MatchedBy(func(ctx context.Context) bool {
			deadline, ok := ctx.Deadline()
			return ok && time.Until(deadline) <= time.Minute
		}), expiredFilter()).Return(nil, nil).Once()

		newTestSweeper(r, time.Minute).sweepOnce(context.Background())

		r.assertExpectations(t)
		assert.Empty(t, r.calls("Update"))
	})

	t.Run("ListErrorStopsTheSweep", func(t *testing.T) {
		r := newTestRepos()
		r.ops.On("List", mock.Anything, expiredFilter()).
			Return(nil, fmt.Errorf("db down")).Once()

		newTestSweeper(r, time.Minute).sweepOnce(context.Background())

		assert.Empty(t, r.calls("Update"))
	})
}

func TestSweeperRun(t *testing.T) {
	r := newTestRepos()
	ctx, cancel := context.WithCancel(context.Background())
	swept := make(chan struct{})

	r.ops.On("List", mock.Anything, expiredFilter()).
		Run(func(mock.Arguments) {
			cancel()
			close(swept)
		}).Return(nil, nil).Once()

	done := make(chan struct{})
	go func() {
		newTestSweeper(r, 10*time.Millisecond).Run(ctx)
		close(done)
	}()

	select {
	case <-swept:
	case <-time.After(5 * time.Second):
		t.Fatal("sweeper did not sweep")
	}
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("sweeper did not stop after cancel")
	}
	r.assertExpectations(t)
}
