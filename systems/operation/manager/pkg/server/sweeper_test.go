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

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/operation/manager/mocks"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

func TestNewSweeper(t *testing.T) {
	repo := &mocks.OperationRepo{}
	s := NewSweeper(repo)

	assert.NotNil(t, s)
	assert.Equal(t, 100, s.batch)
}

func TestSweepOnce(t *testing.T) {
	t.Run("NothingExpired", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewSweeper(repo)

		repo.On("FindExpired", mock.AnythingOfType("time.Time"), 100).Return([]db.Operation{}, nil)

		s.sweepOnce()

		repo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Terminate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("FindExpiredError", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewSweeper(repo)

		repo.On("FindExpired", mock.AnythingOfType("time.Time"), 100).Return(nil, errors.New("db down"))

		s.sweepOnce()

		repo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Terminate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("TimesOutExpired", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewSweeper(repo)

		expired := []db.Operation{
			{Id: uuid.NewV4(), FencingToken: 1, ResourceKey: "node:a"},
			{Id: uuid.NewV4(), FencingToken: 2, ResourceKey: "node:b"},
		}
		repo.On("FindExpired", mock.AnythingOfType("time.Time"), 100).Return(expired, nil)
		repo.On("Terminate", expired[0].Id, uint64(1), db.OperationTimeout,
			mock.AnythingOfType("db.OperationAudit"), "lease expired").
			Return(&db.Operation{Id: expired[0].Id}, nil)
		repo.On("Terminate", expired[1].Id, uint64(2), db.OperationTimeout,
			mock.AnythingOfType("db.OperationAudit"), "lease expired").
			Return(&db.Operation{Id: expired[1].Id}, nil)

		s.sweepOnce()

		repo.AssertExpectations(t)
	})

	t.Run("TerminateErrorIsLoggedNotFatal", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewSweeper(repo)

		op := db.Operation{Id: uuid.NewV4(), FencingToken: 3, ResourceKey: "node:c"}
		repo.On("FindExpired", mock.AnythingOfType("time.Time"), 100).Return([]db.Operation{op}, nil)
		repo.On("Terminate", op.Id, uint64(3), db.OperationTimeout,
			mock.AnythingOfType("db.OperationAudit"), "lease expired").
			Return(nil, errors.New("token mismatch"))

		s.sweepOnce()

		repo.AssertExpectations(t)
	})
}

func TestSweeperRun(t *testing.T) {
	repo := &mocks.OperationRepo{}
	s := NewSweeper(repo)
	s.interval = 10 * time.Millisecond

	repo.On("FindExpired", mock.AnythingOfType("time.Time"), 100).Return([]db.Operation{}, nil).Maybe()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(35 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("sweeper Run did not stop on context cancel")
	}
}
