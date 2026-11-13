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
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/mocks"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

const orgName = "ukama"

func TestStartOperation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("GetByIdempotencyKey", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
		repo.On("Start", mock.Anything, mock.AnythingOfType("*db.Operation"),
			mock.AnythingOfType("time.Duration")).
			Return(func(_ context.Context, op *db.Operation,
				_ time.Duration) *db.Operation {
				op.FencingToken = 1
				return op
			}, nil)

		resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
			RequestedBy: "controller",
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

		repo.On("GetByIdempotencyKey", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
		repo.On("Start", mock.Anything, mock.AnythingOfType("*db.Operation"),
			mock.AnythingOfType("time.Duration")).Return(nil, db.ErrLockConflict)

		resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
			RequestedBy: "controller",
			Type:        "RestartNode",
			System:      "node",
			ResourceKey: "node:abc",
		})

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("RepoFailureIsInternal", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("Start", mock.Anything, mock.AnythingOfType("*db.Operation"),
			mock.AnythingOfType("time.Duration")).Return(nil, fmt.Errorf("db down"))

		resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
			RequestedBy: "controller",
			Type:        "RestartNode",
			System:      "node",
			ResourceKey: "node:abc",
		})

		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, resp)
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
		repo.On("GetByIdempotencyKey", mock.Anything, "key-123").Return(existing, nil)

		resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
			RequestedBy:    "controller",
			Type:           "RestartNode",
			System:         "node",
			ResourceKey:    "node:abc",
			IdempotencyKey: "key-123",
		})

		assert.NoError(t, err)
		assert.Equal(t, existing.Id.String(), resp.Operation.Id)
		// Start must NOT be called when idempotency key matches
		repo.AssertNotCalled(t, "Start", mock.Anything, mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})
}

func TestStartOperationIdempotency(t *testing.T) {
	request := func() *pb.StartOperationRequest {
		return &pb.StartOperationRequest{
			RequestedBy: "controller", Type: "RestartNode", System: "node",
			ResourceKey: "node:abc", IdempotencyKey: "key",
		}
	}
	held := func(edit func(*db.Operation)) *db.Operation {
		op := &db.Operation{
			Id: uuid.NewV4(), Type: "RestartNode", System: "node",
			ResourceKey: "node:abc", FencingToken: 3,
		}
		if edit != nil {
			edit(op)
		}
		return op
	}
	startTaken := func(repo *mocks.OperationRepo) {
		repo.On("Start", mock.Anything, mock.AnythingOfType("*db.Operation"),
			mock.AnythingOfType("time.Duration")).
			Return(nil, db.ErrIdempotencyKeyTaken).Once()
	}

	for name, edit := range map[string]func(*db.Operation){
		"AnotherResource": func(op *db.Operation) { op.ResourceKey = "node:other" },
		"AnotherType":     func(op *db.Operation) { op.Type = "ToggleRadio" },
		"AnotherSystem":   func(op *db.Operation) { op.System = "operation" },
	} {
		t.Run("KeyUsedFor"+name+"IsInvalidArgument", func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)
			repo.On("GetByIdempotencyKey", mock.Anything, "key").
				Return(held(edit), nil)

			resp, err := s.StartOperation(context.Background(), request())

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			assert.Nil(t, resp)
			repo.AssertNotCalled(t, "Start", mock.Anything, mock.Anything,
				mock.Anything)
		})
	}

	t.Run("ConcurrentStartReturnsTheWinner", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)
		winner := held(nil)
		repo.On("GetByIdempotencyKey", mock.Anything, "key").Return(nil, nil).Once()
		startTaken(repo)
		repo.On("GetByIdempotencyKey", mock.Anything, "key").Return(winner, nil).Once()

		resp, err := s.StartOperation(context.Background(), request())

		assert.NoError(t, err)
		assert.Equal(t, winner.Id.String(), resp.Operation.Id)
		assert.Equal(t, uint64(3), resp.Operation.FencingToken)
		repo.AssertExpectations(t)
	})

	t.Run("ConcurrentStartForAnotherResourceIsInvalidArgument", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)
		repo.On("GetByIdempotencyKey", mock.Anything, "key").Return(nil, nil).Once()
		startTaken(repo)
		repo.On("GetByIdempotencyKey", mock.Anything, "key").
			Return(held(func(op *db.Operation) { op.ResourceKey = "node:other" }), nil).
			Once()

		resp, err := s.StartOperation(context.Background(), request())

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("ConcurrentWinnerNotFoundIsAborted", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)
		repo.On("GetByIdempotencyKey", mock.Anything, "key").Return(nil, nil).Twice()
		startTaken(repo)

		resp, err := s.StartOperation(context.Background(), request())

		assert.Equal(t, codes.Aborted, status.Code(err))
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("ConcurrentWinnerReadFailureIsInternal", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)
		repo.On("GetByIdempotencyKey", mock.Anything, "key").Return(nil, nil).Once()
		startTaken(repo)
		repo.On("GetByIdempotencyKey", mock.Anything, "key").
			Return(nil, fmt.Errorf("db down")).Once()

		resp, err := s.StartOperation(context.Background(), request())

		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})
}

func TestGetByResource(t *testing.T) {
	t.Run("ReturnsEmptyWhenFree", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("GetByResource", mock.Anything, "node:free").Return(nil, nil)

		resp, err := s.GetByResource(context.Background(), &pb.GetByResourceRequest{ResourceKey: "node:free"})

		assert.NoError(t, err)
		assert.Nil(t, resp.Operation)
		repo.AssertExpectations(t)
	})

	t.Run("ReturnsHolderWhenLocked", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		holder := &db.Operation{Id: uuid.NewV4(), ResourceKey: "node:abc", FencingToken: 1}
		repo.On("GetByResource", mock.Anything, "node:abc").Return(holder, nil)

		resp, err := s.GetByResource(context.Background(), &pb.GetByResourceRequest{ResourceKey: "node:abc"})

		assert.NoError(t, err)
		assert.Equal(t, holder.Id.String(), resp.Operation.Id)
		repo.AssertExpectations(t)
	})
}

func TestGetOperation(t *testing.T) {
	t.Run("ReturnsOperation", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)
		op := &db.Operation{Id: uuid.NewV4(), ResourceKey: "node:abc", FencingToken: 2}
		repo.On("Get", mock.Anything, op.Id).Return(op, nil).Once()

		resp, err := s.GetOperation(context.Background(),
			&pb.GetOperationRequest{Id: op.Id.String()})

		assert.NoError(t, err)
		assert.Equal(t, op.Id.String(), resp.Operation.Id)
		repo.AssertExpectations(t)
	})

	t.Run("MissingOperationIsNotFound", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)
		id := uuid.NewV4()
		repo.On("Get", mock.Anything, id).Return(nil, gorm.ErrRecordNotFound).Once()

		_, err := s.GetOperation(context.Background(),
			&pb.GetOperationRequest{Id: id.String()})

		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("InvalidId", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		_, err := s.GetOperation(context.Background(),
			&pb.GetOperationRequest{Id: "not-a-uuid"})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
	})
}

func TestCompleteOperation(t *testing.T) {
	id := uuid.NewV4()
	audit := mock.AnythingOfType("db.OperationAudit")
	req := &pb.CompleteOperationRequest{
		Id: id.String(), FencingToken: 7, Actor: "controller", Reason: "dispatched",
	}

	t.Run("UsesCallerToken", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("Terminate", mock.Anything, id, uint64(7), db.OperationSuccess, audit, "").
			Return(&db.Operation{Id: id, Status: db.OperationSuccess}, nil).Once()

		resp, err := s.CompleteOperation(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_SUCCESS, resp.Operation.Status)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("PreconditionMapsToFailedPrecondition", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("Terminate", mock.Anything, id, uint64(7), db.OperationSuccess, audit, "").
			Return(nil, fmt.Errorf("token mismatch: %w", store.ErrPrecondition)).Once()

		resp, err := s.CompleteOperation(context.Background(), req)

		assert.Nil(t, resp)
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("MissingOperationMapsToNotFound", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("Terminate", mock.Anything, id, uint64(7), db.OperationSuccess, audit, "").
			Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := s.CompleteOperation(context.Background(), req)

		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("InvalidId", func(t *testing.T) {
		s := NewOperationServer(orgName, "", &mocks.OperationRepo{}, nil)

		resp, err := s.CompleteOperation(context.Background(),
			&pb.CompleteOperationRequest{Id: "not-a-uuid"})

		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

func TestFailOperation(t *testing.T) {
	id := uuid.NewV4()
	audit := mock.AnythingOfType("db.OperationAudit")
	req := &pb.FailOperationRequest{
		Id: id.String(), FencingToken: 7, Actor: "controller", Reason: "publish failed",
	}

	t.Run("UsesCallerTokenAndRecordsReason", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("Terminate", mock.Anything, id, uint64(7), db.OperationFailed, audit,
			"publish failed").
			Return(&db.Operation{Id: id, Status: db.OperationFailed}, nil).Once()

		resp, err := s.FailOperation(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_FAILED, resp.Operation.Status)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("PreconditionMapsToFailedPrecondition", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		repo.On("Terminate", mock.Anything, id, uint64(7), db.OperationFailed, audit,
			"publish failed").
			Return(nil, fmt.Errorf("already success: %w", store.ErrPrecondition)).Once()

		resp, err := s.FailOperation(context.Background(), req)

		assert.Nil(t, resp)
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})
}

func TestForceUnlock(t *testing.T) {
	t.Run("CancelsAndReleases", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		id := uuid.NewV4()
		current := &db.Operation{Id: id, ResourceKey: "node:abc", FencingToken: 9, Status: db.OperationRunning}
		cancelled := &db.Operation{Id: id, ResourceKey: "node:abc", FencingToken: 9, Status: db.OperationCancelled}

		repo.On("Get", mock.Anything, id).Return(current, nil)
		repo.On("Terminate", mock.Anything, id, uint64(9), db.OperationCancelled,
			mock.AnythingOfType("db.OperationAudit"), mock.Anything).
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

	t.Run("FinishedOperationMapsToFailedPrecondition", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		id := uuid.NewV4()

		repo.On("Get", mock.Anything, id).Return(&db.Operation{Id: id, FencingToken: 9}, nil)
		repo.On("Terminate", mock.Anything, id, uint64(9), db.OperationCancelled,
			mock.AnythingOfType("db.OperationAudit"), mock.Anything).
			Return(nil, fmt.Errorf("already success: %w", store.ErrPrecondition))

		resp, err := s.ForceUnlock(context.Background(), &pb.ForceUnlockRequest{
			Id: id.String(), Actor: "owner@ukama.com", Reason: "stuck operation",
		})

		assert.Nil(t, resp)
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})
}

func TestMarkRunning(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)

		id := uuid.NewV4()
		running := &db.Operation{Id: id, ResourceKey: "node:abc", FencingToken: 3, Status: db.OperationRunning}
		repo.On("MarkRunning", mock.Anything, id, uint64(3), "controller").Return(running, nil)

		resp, err := s.MarkRunning(context.Background(), &pb.MarkRunningRequest{
			Actor:        "controller",
			Id:           id.String(),
			FencingToken: 3,
		})

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_RUNNING, resp.Operation.Status)
		repo.AssertExpectations(t)
	})
}

func TestRepoErrorCodes(t *testing.T) {
	id := uuid.NewV4()
	lockTimeout := fmt.Errorf("terminate: %w", db.ErrLockTimeout)

	for name, tc := range map[string]struct {
		err  error
		code codes.Code
	}{
		"LockTimeoutIsUnavailable":   {lockTimeout, codes.Unavailable},
		"DeadlineIsDeadlineExceeded": {context.DeadlineExceeded, codes.DeadlineExceeded},
		"CancelledIsCanceled":        {context.Canceled, codes.Canceled},
		"OtherErrorIsInternal":       {fmt.Errorf("db down"), codes.Internal},
	} {
		t.Run("Complete"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)
			repo.On("Terminate", mock.Anything, id, uint64(7), db.OperationSuccess,
				mock.Anything, "").Return(nil, tc.err).Once()

			_, err := s.CompleteOperation(context.Background(),
				&pb.CompleteOperationRequest{
					Id: id.String(), FencingToken: 7, Actor: "a", Reason: "r",
				})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("Start"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)
			repo.On("Start", mock.Anything, mock.Anything, mock.Anything).
				Return(nil, tc.err).Once()

			_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
				RequestedBy: "controller",
				Type:        "RestartNode", System: "node", ResourceKey: "node:abc",
			})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("StartIdempotencyLookup"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)
			repo.On("GetByIdempotencyKey", mock.Anything, "key").
				Return(nil, tc.err).Once()

			_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
				Type: "RestartNode", System: "node", ResourceKey: "node:abc",
				RequestedBy: "controller", IdempotencyKey: "key",
			})

			assert.Equal(t, tc.code, status.Code(err))
			repo.AssertNotCalled(t, "Start", mock.Anything, mock.Anything, mock.Anything)
		})

		t.Run("GetOperation"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)
			repo.On("Get", mock.Anything, id).Return(nil, tc.err).Once()

			_, err := s.GetOperation(context.Background(),
				&pb.GetOperationRequest{Id: id.String()})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("GetByResource"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)
			repo.On("GetByResource", mock.Anything, "node:abc").
				Return(nil, tc.err).Once()

			_, err := s.GetByResource(context.Background(),
				&pb.GetByResourceRequest{ResourceKey: "node:abc"})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("ForceUnlock"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)
			repo.On("Get", mock.Anything, id).Return(nil, tc.err).Once()

			_, err := s.ForceUnlock(context.Background(), &pb.ForceUnlockRequest{
				Id: id.String(), Actor: "a", Reason: "r",
			})

			assert.Equal(t, tc.code, status.Code(err))
			repo.AssertNotCalled(t, "Terminate", mock.Anything, mock.Anything,
				mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
	}

	t.Run("MarkRunningLockTimeoutIsUnavailable", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		s := NewOperationServer(orgName, "", repo, nil)
		repo.On("MarkRunning", mock.Anything, id, uint64(3), "controller").
			Return(nil, lockTimeout).Once()

		_, err := s.MarkRunning(context.Background(), &pb.MarkRunningRequest{
			Actor: "controller",
			Id:    id.String(), FencingToken: 3,
		})

		assert.Equal(t, codes.Unavailable, status.Code(err))
	})
}

var invalidResourceKeys = map[string]string{
	"Empty":          "",
	"Blank":          "   ",
	"LeadingSpace":   " node:abc",
	"InnerSpace":     "node abc",
	"Tab":            "node:\tabc",
	"Nul":            "node:\x00abc",
	"InvalidUtf8":    "node:\xffabc",
	"TooLong":        strings.Repeat("a", maxResourceKeyLen+1),
	"ControlNewline": "node:abc\n",
}

func TestValidateResourceKey(t *testing.T) {
	for _, key := range []string{
		"node:uk-sa2450-tnode-v0-4e86",
		"config:apply",
		"test:integration:" + uuid.NewV4().String(),
		strings.Repeat("a", maxResourceKeyLen),
	} {
		assert.NoError(t, validateResourceKey(key), key)
	}

	for name, key := range invalidResourceKeys {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, codes.InvalidArgument, status.Code(validateResourceKey(key)))
		})
	}
}

func TestInvalidResourceKeyIsRejected(t *testing.T) {
	for name, key := range invalidResourceKeys {
		t.Run("Start"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)

			_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
				RequestedBy: "controller",
				Type:        "RestartNode", System: "node", ResourceKey: key,
				IdempotencyKey: "key",
			})

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			repo.AssertNotCalled(t, "GetByIdempotencyKey", mock.Anything, mock.Anything)
			repo.AssertNotCalled(t, "Start", mock.Anything, mock.Anything, mock.Anything)
		})

		t.Run("GetByResource"+name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)

			_, err := s.GetByResource(context.Background(),
				&pb.GetByResourceRequest{ResourceKey: key})

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			repo.AssertNotCalled(t, "GetByResource", mock.Anything, mock.Anything)
		})
	}
}

func TestStartOperationLease(t *testing.T) {
	for name, tc := range map[string]struct {
		seconds uint32
		lease   time.Duration
	}{
		"ZeroUsesDefault": {0, pkg.DefaultLeaseTTL},
		"MaxIsAccepted":   {3600, time.Hour},
	} {
		t.Run(name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)

			repo.On("Start", mock.Anything, mock.AnythingOfType("*db.Operation"), tc.lease).
				Return(func(_ context.Context, op *db.Operation,
					_ time.Duration) *db.Operation {
					return op
				}, nil).Once()

			_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
				RequestedBy: "controller",
				Type:        "RestartNode", System: "node", ResourceKey: "node:abc",
				LeaseSeconds: tc.seconds,
			})

			assert.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}

	for name, seconds := range map[string]uint32{
		"AboveMaxIsRejected":  3601,
		"MaxUint32IsRejected": math.MaxUint32,
	} {
		t.Run(name, func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			s := NewOperationServer(orgName, "", repo, nil)

			_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
				RequestedBy: "controller",
				Type:        "RestartNode", System: "node", ResourceKey: "node:abc",
				IdempotencyKey: "key", LeaseSeconds: seconds,
			})

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			repo.AssertNotCalled(t, "GetByIdempotencyKey", mock.Anything, mock.Anything)
			repo.AssertNotCalled(t, "Start", mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

func TestBlankActorsAreRejected(t *testing.T) {
	id := uuid.NewV4().String()

	for _, blank := range []string{"", "   "} {
		calls := map[string]func(s *OperationServer) error{
			"StartRequestedBy": func(s *OperationServer) error {
				_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
					Type: "RestartNode", System: "node", ResourceKey: "node:abc",
					RequestedBy: blank,
				})
				return err
			},
			"MarkRunningActor": func(s *OperationServer) error {
				_, err := s.MarkRunning(context.Background(), &pb.MarkRunningRequest{
					Id: id, FencingToken: 1, Actor: blank,
				})
				return err
			},
			"CompleteActor": func(s *OperationServer) error {
				_, err := s.CompleteOperation(context.Background(),
					&pb.CompleteOperationRequest{
						Id: id, FencingToken: 1, Actor: blank, Reason: "r",
					})
				return err
			},
			"CompleteReason": func(s *OperationServer) error {
				_, err := s.CompleteOperation(context.Background(),
					&pb.CompleteOperationRequest{
						Id: id, FencingToken: 1, Actor: "a", Reason: blank,
					})
				return err
			},
			"FailActor": func(s *OperationServer) error {
				_, err := s.FailOperation(context.Background(), &pb.FailOperationRequest{
					Id: id, FencingToken: 1, Actor: blank, Reason: "r",
				})
				return err
			},
			"FailReason": func(s *OperationServer) error {
				_, err := s.FailOperation(context.Background(), &pb.FailOperationRequest{
					Id: id, FencingToken: 1, Actor: "a", Reason: blank,
				})
				return err
			},
			"ForceUnlockActor": func(s *OperationServer) error {
				_, err := s.ForceUnlock(context.Background(), &pb.ForceUnlockRequest{
					Id: id, Actor: blank, Reason: "r",
				})
				return err
			},
			"ForceUnlockReason": func(s *OperationServer) error {
				_, err := s.ForceUnlock(context.Background(), &pb.ForceUnlockRequest{
					Id: id, Actor: "a", Reason: blank,
				})
				return err
			},
		}

		for name, call := range calls {
			t.Run(fmt.Sprintf("%s%q", name, blank), func(t *testing.T) {
				repo := &mocks.OperationRepo{}

				err := call(NewOperationServer(orgName, "", repo, nil))

				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				repo.AssertExpectations(t)
			})
		}
	}
}
