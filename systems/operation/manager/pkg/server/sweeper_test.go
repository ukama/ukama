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

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/mocks"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

func newTestSweeper(repo db.OperationRepo, interval time.Duration) *Sweeper {
	return &Sweeper{repo: repo, interval: interval, batch: 10}
}

func timeoutAudit() interface{} {
	return mock.MatchedBy(func(a db.OperationAudit) bool {
		return a.Event == "timeout" && a.Reason == "lease expired" &&
			a.Actor == pkg.OperationManagerSweeperActor
	})
}

func TestSweeperSweepOnce(t *testing.T) {
	t.Run("TimesOutEveryExpiredOperation", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		first := db.Operation{Id: uuid.NewV4(), FencingToken: 3}
		second := db.Operation{Id: uuid.NewV4(), FencingToken: 4}

		repo.On("FindExpired", mock.Anything, mock.Anything, 10).
			Return([]db.Operation{first, second}, nil).Once()
		repo.On("Terminate", mock.Anything, first.Id, first.FencingToken,
			db.OperationTimeout, timeoutAudit(), "lease expired").
			Return(nil, fmt.Errorf("db down")).Once()
		repo.On("Terminate", mock.Anything, second.Id, second.FencingToken,
			db.OperationTimeout, timeoutAudit(), "lease expired").
			Return(&second, nil).Once()

		newTestSweeper(repo, time.Minute).sweepOnce(context.Background())

		repo.AssertExpectations(t)
	})

	t.Run("BoundsEachSweepByTheInterval", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		repo.On("FindExpired", mock.MatchedBy(func(ctx context.Context) bool {
			deadline, ok := ctx.Deadline()
			return ok && time.Until(deadline) <= time.Minute
		}), mock.Anything, 10).Return(nil, nil).Once()

		newTestSweeper(repo, time.Minute).sweepOnce(context.Background())

		repo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Terminate", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("FindExpiredErrorStopsTheSweep", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		repo.On("FindExpired", mock.Anything, mock.Anything, 10).
			Return(nil, fmt.Errorf("db down")).Once()

		newTestSweeper(repo, time.Minute).sweepOnce(context.Background())

		repo.AssertNotCalled(t, "Terminate", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestSweeperRun(t *testing.T) {
	repo := &mocks.OperationRepo{}
	ctx, cancel := context.WithCancel(context.Background())
	swept := make(chan struct{})

	repo.On("FindExpired", mock.Anything, mock.Anything, 10).
		Run(func(mock.Arguments) {
			cancel()
			close(swept)
		}).Return(nil, nil).Once()

	done := make(chan struct{})
	go func() {
		newTestSweeper(repo, 10*time.Millisecond).Run(ctx)
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
	repo.AssertExpectations(t)
}
