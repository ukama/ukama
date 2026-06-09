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
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/mocks"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

const orgName = "ukama"

func TestStartOperation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("GetByIdempotencyKey", mock.Anything).Return(nil, nil).Maybe()
		repo.On("Start", mock.AnythingOfType("*db.Operation"), mock.AnythingOfType("time.Duration")).
			Return(func(op *db.Operation, _ time.Duration) *db.Operation {
				op.FencingToken = 1
				return op
			}, nil)

		resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
			Type:        "RestartNode",
			System:      "node",
			ResourceKey: "node:abc",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp.Operation)
		assert.Equal(t, "node:abc", resp.Operation.ResourceKey)
		repo.AssertExpectations(t)
	})

	t.Run("ConflictReturnsAlreadyExists", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		holder := &db.Operation{
			Id:           uuid.NewV4(),
			Type:         "RestartNode",
			System:       "node",
			Status:       db.OperationRunning,
			FencingToken: 5,
			ResourceKey:  "node:abc",
		}

		repo.On("GetByIdempotencyKey", mock.Anything).Return(nil, nil).Maybe()
		repo.On("Start", mock.AnythingOfType("*db.Operation"), mock.AnythingOfType("time.Duration")).
			Return(holder, db.ErrLockConflict)

		resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
			Type:        "RestartNode",
			System:      "node",
			ResourceKey: "node:abc",
		})

		assert.Error(t, err)
		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.NotNil(t, resp.ConflictingOperation)
		assert.Equal(t, holder.Id.String(), resp.ConflictingOperation.Id)
		repo.AssertExpectations(t)
	})

	t.Run("IdempotencyShortCircuits", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		existing := &db.Operation{
			Id:           uuid.NewV4(),
			Type:         "RestartNode",
			System:       "node",
			FencingToken: 2,
			ResourceKey:  "node:abc",
		}
		repo.On("GetByIdempotencyKey", "key-123").Return(existing, nil)

		resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
			Type:           "RestartNode",
			System:         "node",
			ResourceKey:    "node:abc",
			IdempotencyKey: "key-123",
		})

		assert.NoError(t, err)
		assert.Equal(t, existing.Id.String(), resp.Operation.Id)
		// Start must NOT be called when idempotency key matches
		repo.AssertNotCalled(t, "Start", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})
}

func TestGetByResource(t *testing.T) {
	t.Run("ReturnsEmptyWhenFree", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("GetByResource", "node:free").Return(nil, nil)

		resp, err := s.GetByResource(context.Background(), &pb.GetByResourceRequest{ResourceKey: "node:free"})

		assert.NoError(t, err)
		assert.Nil(t, resp.Operation)
		repo.AssertExpectations(t)
	})

	t.Run("ReturnsHolderWhenLocked", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		holder := &db.Operation{Id: uuid.NewV4(), ResourceKey: "node:abc", FencingToken: 1}
		repo.On("GetByResource", "node:abc").Return(holder, nil)

		resp, err := s.GetByResource(context.Background(), &pb.GetByResourceRequest{ResourceKey: "node:abc"})

		assert.NoError(t, err)
		assert.Equal(t, holder.Id.String(), resp.Operation.Id)
		repo.AssertExpectations(t)
	})
}

func TestForceUnlock(t *testing.T) {
	t.Run("CancelsAndReleases", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		id := uuid.NewV4()
		current := &db.Operation{Id: id, ResourceKey: "node:abc", FencingToken: 9, Status: db.OperationRunning}
		cancelled := &db.Operation{Id: id, ResourceKey: "node:abc", FencingToken: 9, Status: db.OperationCancelled}

		repo.On("Get", id).Return(current, nil)
		repo.On("Terminate", id, uint64(9), db.OperationCancelled, mock.AnythingOfType("db.OperationAudit"), mock.Anything).
			Return(cancelled, nil)

		resp, err := s.ForceUnlock(context.Background(), &pb.ForceUnlockRequest{
			Id:     id.String(),
			Actor:  "owner@ukama.com",
			Reason: "stuck operation",
		})

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_CANCELLED, resp.Operation.Status)
		repo.AssertExpectations(t)
	})
}

func TestMarkRunning(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		id := uuid.NewV4()
		running := &db.Operation{Id: id, ResourceKey: "node:abc", FencingToken: 3, Status: db.OperationRunning}
		repo.On("MarkRunning", id, uint64(3)).Return(running, nil)

		resp, err := s.MarkRunning(context.Background(), &pb.MarkRunningRequest{
			Id:           id.String(),
			FencingToken: 3,
		})

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_RUNNING, resp.Operation.Status)
		repo.AssertExpectations(t)
	})
}
