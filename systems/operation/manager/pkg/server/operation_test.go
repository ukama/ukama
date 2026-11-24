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
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

func startRequest(keys ...string) *pb.StartOperationRequest {
	return &pb.StartOperationRequest{
		RequestedBy: "controller", Type: "RestartNode", System: "node",
		ResourceKeys: keys,
	}
}

func TestStartOperation(t *testing.T) {
	t.Run("SingleTargetIsLeaderWithoutGroup", func(t *testing.T) {
		r := newTestRepos()
		r.adds([]string{}, mock.AnythingOfType("time.Duration"))

		resp, err := r.server().StartOperation(context.Background(),
			startRequest("node:abc"))

		assert.NoError(t, err)
		assert.Len(t, resp.Operations, 1)
		assert.Equal(t, "node:abc", resp.Operations[0].ResourceKey)
		assert.Equal(t, uint64(1), resp.Operations[0].FencingToken)
		assert.Empty(t, resp.Operations[0].GroupId)
		assert.Equal(t, []string{"lock_acquired"}, r.events())
		assert.Equal(t, "controller", r.saved[0].Actor)
		assert.Equal(t, "node:abc", r.saved[0].ResourceKey)
		r.assertExpectations(t)
	})

	t.Run("GroupSortsTargetsAndLinksMembersToLeader", func(t *testing.T) {
		r := newTestRepos()
		r.findsGroup(nil)
		r.adds([]string{"node:c", "node:t"}, mock.AnythingOfType("time.Duration"))
		req := startRequest("node:t", "node:a", "node:t")
		req.MustBeFree = []string{"node:t", "node:c"}
		req.IdempotencyKey = "key"

		resp, err := r.server().StartOperation(context.Background(), req)

		assert.NoError(t, err)
		ops := r.calls("Add")[0].Arguments.Get(1).([]*db.Operation)
		assert.Len(t, ops, 2)
		assert.Equal(t, "node:a", ops[0].ResourceKey)
		assert.Nil(t, ops[0].GroupId)
		assert.Equal(t, "key", *ops[0].IdempotencyKey)
		assert.Equal(t, "node:t", ops[1].ResourceKey)
		assert.Equal(t, ops[0].Id, *ops[1].GroupId)
		assert.Nil(t, ops[1].IdempotencyKey)
		assert.Equal(t, ops[0].LeaseExpiresAt, ops[1].LeaseExpiresAt)
		assert.Equal(t, ops, r.locks.Calls[0].Arguments.Get(1))
		assert.Len(t, resp.Operations, 2)
		assert.Equal(t, ops[0].Id.String(), resp.Operations[1].GroupId)
		assert.Equal(t, []string{"lock_acquired", "lock_acquired"}, r.events())
		r.assertExpectations(t)
	})

	for name, req := range map[string]*pb.StartOperationRequest{
		"NoTargets":     startRequest(),
		"InvalidTarget": startRequest("node abc"),
		"InvalidMustBeFree": {RequestedBy: "controller", Type: "RestartNode",
			System: "node", ResourceKeys: []string{"node:a"},
			MustBeFree: []string{""}},
		"BlankRequestedBy": {Type: "RestartNode", System: "node",
			ResourceKeys: []string{"node:a"}, RequestedBy: " "},
	} {
		t.Run(name+"IsInvalidArgument", func(t *testing.T) {
			r := newTestRepos()

			_, err := r.server().StartOperation(context.Background(), req)

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			r.assertUntouched(t)
		})
	}

	t.Run("ConflictReturnsAlreadyExistsWithHolder", func(t *testing.T) {
		r := newTestRepos()
		holder := uuid.NewV4()
		r.addRunsNested()
		r.locks.On("Acquire", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).
			Return(&db.ConflictError{ResourceKey: "node:c", HolderId: holder,
				HolderType: "UpdateSoftware"}).Once()

		resp, err := r.server().StartOperation(context.Background(),
			startRequest("node:abc"))

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.Contains(t, status.Convert(err).Message(), holder.String())
		assert.Contains(t, status.Convert(err).Message(), "UpdateSoftware")
		assert.Nil(t, resp)
		r.audits.AssertNotCalled(t, "AddInTx", mock.Anything, mock.Anything)
		r.assertExpectations(t)
	})

	t.Run("AuditFailureIsInternal", func(t *testing.T) {
		r := newTestRepos()
		r.addRunsNested()
		r.locks.On("Acquire", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(nil).Once()
		r.audits.On("AddInTx", mock.Anything, mock.Anything).
			Return(fmt.Errorf("db down")).Once()

		resp, err := r.server().StartOperation(context.Background(),
			startRequest("node:abc"))

		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, resp)
		r.assertExpectations(t)
	})

	t.Run("RepoFailureIsInternal", func(t *testing.T) {
		r := newTestRepos()
		r.addFails(fmt.Errorf("db down"))

		resp, err := r.server().StartOperation(context.Background(),
			startRequest("node:abc"))

		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, resp)
		r.assertExpectations(t)
	})
}

func TestStartOperationIdempotency(t *testing.T) {
	request := func() *pb.StartOperationRequest {
		req := startRequest("node:t", "node:a")
		req.IdempotencyKey = "key"
		return req
	}
	group := func(edit func([]*db.Operation) []*db.Operation) []*db.Operation {
		leader := &db.Operation{Id: uuid.NewV4(), Type: "RestartNode",
			System: "node", ResourceKey: "node:a", FencingToken: 3}
		member := &db.Operation{Id: uuid.NewV4(), GroupId: &leader.Id,
			Type: "RestartNode", System: "node", ResourceKey: "node:t",
			FencingToken: 4}
		ops := []*db.Operation{leader, member}
		if edit != nil {
			ops = edit(ops)
		}
		return ops
	}

	t.Run("ReplayReturnsTheWholeGroup", func(t *testing.T) {
		r := newTestRepos()
		existing := group(nil)
		r.findsGroup(existing)

		resp, err := r.server().StartOperation(context.Background(), request())

		assert.NoError(t, err)
		assert.Len(t, resp.Operations, 2)
		assert.Equal(t, existing[0].Id.String(), resp.Operations[0].Id)
		assert.Equal(t, existing[0].Id.String(), resp.Operations[1].GroupId)
		assert.Empty(t, r.calls("Add"))
		r.assertExpectations(t)
	})

	for name, edit := range map[string]func([]*db.Operation) []*db.Operation{
		"AnotherResource": func(ops []*db.Operation) []*db.Operation {
			ops[1].ResourceKey = "node:other"
			return ops
		},
		"AnotherType": func(ops []*db.Operation) []*db.Operation {
			ops[0].Type = "ToggleRadio"
			return ops
		},
		"AnotherSystem": func(ops []*db.Operation) []*db.Operation {
			ops[1].System = "operation"
			return ops
		},
		"FewerTargets": func(ops []*db.Operation) []*db.Operation {
			return ops[:1]
		},
	} {
		t.Run("KeyUsedFor"+name+"IsInvalidArgument", func(t *testing.T) {
			r := newTestRepos()
			r.findsGroup(group(edit))

			resp, err := r.server().StartOperation(context.Background(), request())

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			assert.Nil(t, resp)
			assert.Empty(t, r.calls("Add"))
		})
	}

	for name, startErr := range map[string]error{
		"TakenKey":         db.ErrIdempotencyKeyTaken,
		"ConflictOnTarget": &db.ConflictError{ResourceKey: "node:a"},
	} {
		t.Run("Concurrent"+name+"ReturnsTheWinner", func(t *testing.T) {
			r := newTestRepos()
			winner := group(nil)
			r.findsGroup(nil)
			r.addFails(startErr)
			r.findsGroup(winner)

			resp, err := r.server().StartOperation(context.Background(), request())

			assert.NoError(t, err)
			assert.Equal(t, winner[0].Id.String(), resp.Operations[0].Id)
			assert.Equal(t, uint64(4), resp.Operations[1].FencingToken)
			r.assertExpectations(t)
		})
	}

	t.Run("ConflictWithoutGroupIsAlreadyExists", func(t *testing.T) {
		r := newTestRepos()
		r.findsGroup(nil)
		r.addFails(&db.ConflictError{ResourceKey: "node:a"})
		r.findsGroup(nil)

		resp, err := r.server().StartOperation(context.Background(), request())

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.Nil(t, resp)
		r.assertExpectations(t)
	})

	t.Run("ConcurrentStartForAnotherResourceIsInvalidArgument", func(t *testing.T) {
		r := newTestRepos()
		r.findsGroup(nil)
		r.addFails(db.ErrIdempotencyKeyTaken)
		r.findsGroup(group(func(ops []*db.Operation) []*db.Operation {
			ops[0].ResourceKey = "node:other"
			return ops
		}))

		resp, err := r.server().StartOperation(context.Background(), request())

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Nil(t, resp)
		r.assertExpectations(t)
	})

	t.Run("ConcurrentWinnerNotFoundIsAborted", func(t *testing.T) {
		r := newTestRepos()
		r.findsGroup(nil)
		r.addFails(db.ErrIdempotencyKeyTaken)
		r.findsGroup(nil)

		resp, err := r.server().StartOperation(context.Background(), request())

		assert.Equal(t, codes.Aborted, status.Code(err))
		assert.Nil(t, resp)
		r.assertExpectations(t)
	})

	t.Run("ConcurrentWinnerReadFailureIsInternal", func(t *testing.T) {
		r := newTestRepos()
		r.findsGroup(nil)
		r.addFails(db.ErrIdempotencyKeyTaken)
		r.groupLookupFails(fmt.Errorf("db down"))

		resp, err := r.server().StartOperation(context.Background(), request())

		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, resp)
		r.assertExpectations(t)
	})

	t.Run("MemberReadFailureIsInternal", func(t *testing.T) {
		r := newTestRepos()
		leader := group(nil)[0]
		r.ops.On("List", mock.Anything, leaderFilter("key")).
			Return([]*db.Operation{leader}, nil).Once()
		r.ops.On("List", mock.Anything, membersFilter(leader.Id)).
			Return(nil, fmt.Errorf("db down")).Once()

		resp, err := r.server().StartOperation(context.Background(), request())

		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, resp)
		assert.Empty(t, r.calls("Add"))
		r.assertExpectations(t)
	})
}

func TestGetByResource(t *testing.T) {
	t.Run("ReturnsEmptyWhenFree", func(t *testing.T) {
		r := newTestRepos()
		r.locks.On("Holder", mock.Anything, "node:free").Return(nil, nil).Once()

		resp, err := r.server().GetByResource(context.Background(),
			&pb.GetByResourceRequest{ResourceKey: "node:free"})

		assert.NoError(t, err)
		assert.Nil(t, resp.Operation)
		r.assertExpectations(t)
	})

	t.Run("ReturnsHolderWhenLocked", func(t *testing.T) {
		r := newTestRepos()
		holder := &db.Operation{Id: uuid.NewV4(), ResourceKey: "node:abc",
			FencingToken: 1}
		r.locks.On("Holder", mock.Anything, "node:abc").
			Return(&db.ResourceLock{ResourceKey: "node:abc", OperationId: holder.Id},
				nil).Once()
		r.ops.On("Get", mock.Anything, holder.Id).Return(holder, nil).Once()

		resp, err := r.server().GetByResource(context.Background(),
			&pb.GetByResourceRequest{ResourceKey: "node:abc"})

		assert.NoError(t, err)
		assert.Equal(t, holder.Id.String(), resp.Operation.Id)
		r.assertExpectations(t)
	})
}

func TestGetOperation(t *testing.T) {
	t.Run("ReturnsOperation", func(t *testing.T) {
		r := newTestRepos()
		op := &db.Operation{Id: uuid.NewV4(), ResourceKey: "node:abc", FencingToken: 2}
		r.ops.On("Get", mock.Anything, op.Id).Return(op, nil).Once()

		resp, err := r.server().GetOperation(context.Background(),
			&pb.GetOperationRequest{Id: op.Id.String()})

		assert.NoError(t, err)
		assert.Equal(t, op.Id.String(), resp.Operation.Id)
		r.assertExpectations(t)
	})

	t.Run("MissingOperationIsNotFound", func(t *testing.T) {
		r := newTestRepos()
		id := uuid.NewV4()
		r.ops.On("Get", mock.Anything, id).Return(nil, gorm.ErrRecordNotFound).Once()

		_, err := r.server().GetOperation(context.Background(),
			&pb.GetOperationRequest{Id: id.String()})

		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("InvalidId", func(t *testing.T) {
		r := newTestRepos()

		_, err := r.server().GetOperation(context.Background(),
			&pb.GetOperationRequest{Id: "not-a-uuid"})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		r.assertUntouched(t)
	})
}

func runningOp(token uint64) db.Operation {
	op := pendingOp()
	op.Status = db.OperationRunning
	op.FencingToken = token
	return op
}

func TestCompleteOperation(t *testing.T) {
	request := func(id uuid.UUID) *pb.CompleteOperationRequest {
		return &pb.CompleteOperationRequest{
			Id: id.String(), FencingToken: 7, Actor: "controller", Reason: "dispatched",
		}
	}

	t.Run("UsesCallerToken", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(7)
		r.updates(op)

		resp, err := r.server().CompleteOperation(context.Background(), request(op.Id))

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_SUCCESS, resp.Operation.Status)
		assert.Equal(t, []uuid.UUID{op.Id}, r.released)
		assert.Equal(t, []string{"completed"}, r.events())
		assert.Equal(t, "controller", r.saved[0].Actor)
		assert.Equal(t, "dispatched", r.saved[0].Reason)
		assert.Empty(t, r.calls("Get"))
	})

	t.Run("PreconditionMapsToFailedPrecondition", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(9)
		r.updates(op)

		resp, err := r.server().CompleteOperation(context.Background(), request(op.Id))

		assert.Nil(t, resp)
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assert.Equal(t, []string{"rejected_completed"}, r.events())
	})

	t.Run("MissingOperationMapsToNotFound", func(t *testing.T) {
		r := newTestRepos()
		id := uuid.NewV4()
		r.updateFails(id, fmt.Errorf("operation %s: %w", id, store.ErrNotFound))

		resp, err := r.server().CompleteOperation(context.Background(), request(id))

		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("InvalidId", func(t *testing.T) {
		r := newTestRepos()

		resp, err := r.server().CompleteOperation(context.Background(),
			&pb.CompleteOperationRequest{Id: "not-a-uuid", Actor: "a", Reason: "r"})

		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		r.assertUntouched(t)
	})
}

func TestFailOperation(t *testing.T) {
	request := func(id uuid.UUID) *pb.FailOperationRequest {
		return &pb.FailOperationRequest{
			Id: id.String(), FencingToken: 7, Actor: "controller",
			Reason: "publish failed",
		}
	}

	t.Run("UsesCallerTokenAndRecordsReason", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(7)
		r.updates(op)

		resp, err := r.server().FailOperation(context.Background(), request(op.Id))

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_FAILED, resp.Operation.Status)
		assert.Equal(t, "publish failed", resp.Operation.Error)
		assert.Equal(t, []string{"failed"}, r.events())
		assert.Equal(t, "publish failed", r.saved[0].Reason)
		assert.Empty(t, r.calls("Get"))
	})

	t.Run("PreconditionMapsToFailedPrecondition", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(7)
		op.Status = db.OperationSuccess
		r.updates(op)

		resp, err := r.server().FailOperation(context.Background(), request(op.Id))

		assert.Nil(t, resp)
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assert.Equal(t, []string{"ignored_failed"}, r.events())
	})

	t.Run("InvalidId", func(t *testing.T) {
		r := newTestRepos()

		resp, err := r.server().FailOperation(context.Background(),
			&pb.FailOperationRequest{Id: "not-a-uuid", Actor: "a", Reason: "r"})

		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		r.assertUntouched(t)
	})
}

func TestForceUnlock(t *testing.T) {
	request := func(id uuid.UUID) *pb.ForceUnlockRequest {
		return &pb.ForceUnlockRequest{
			Id: id.String(), Actor: "owner@ukama.com", Reason: "stuck operation",
		}
	}

	t.Run("CancelsAndReleases", func(t *testing.T) {
		r := newTestRepos()
		current := runningOp(9)
		r.ops.On("Get", mock.Anything, current.Id).Return(&current, nil).Once()
		r.updates(current)

		resp, err := r.server().ForceUnlock(context.Background(), request(current.Id))

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_CANCELLED, resp.Operation.Status)
		assert.Equal(t, "force unlocked: stuck operation", resp.Operation.Error)
		assert.Equal(t, []uuid.UUID{current.Id}, r.released)
		assert.Equal(t, []string{"force_unlock"}, r.events())
		assert.Equal(t, "owner@ukama.com", r.saved[0].Actor)
		r.assertExpectations(t)
	})

	t.Run("FinishedOperationMapsToFailedPrecondition", func(t *testing.T) {
		r := newTestRepos()
		current := runningOp(9)
		current.Status = db.OperationSuccess
		r.ops.On("Get", mock.Anything, current.Id).Return(&current, nil).Once()
		r.updates(current)

		resp, err := r.server().ForceUnlock(context.Background(), request(current.Id))

		assert.Nil(t, resp)
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assert.Equal(t, []string{"ignored_force_unlock"}, r.events())
	})

	t.Run("MissingOperationIsNotFound", func(t *testing.T) {
		r := newTestRepos()
		id := uuid.NewV4()
		r.ops.On("Get", mock.Anything, id).Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := r.server().ForceUnlock(context.Background(), request(id))

		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Empty(t, r.calls("Update"))
	})

	t.Run("InvalidId", func(t *testing.T) {
		r := newTestRepos()

		resp, err := r.server().ForceUnlock(context.Background(),
			&pb.ForceUnlockRequest{Id: "not-a-uuid", Actor: "a", Reason: "r"})

		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		r.assertUntouched(t)
	})
}

func TestMarkRunning(t *testing.T) {
	request := func(id uuid.UUID) *pb.MarkRunningRequest {
		return &pb.MarkRunningRequest{Actor: "controller", Id: id.String(),
			FencingToken: 3}
	}

	t.Run("Success", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		r.updates(op)

		resp, err := r.server().MarkRunning(context.Background(), request(op.Id))

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_RUNNING, resp.Operation.Status)
		assert.Equal(t, []string{"running"}, r.events())
		r.assertExpectations(t)
	})

	t.Run("UnknownIsNotFound", func(t *testing.T) {
		r := newTestRepos()
		id := uuid.NewV4()
		r.updateFails(id, fmt.Errorf("operation: %w", store.ErrNotFound))

		_, err := r.server().MarkRunning(context.Background(), request(id))

		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("RejectedIsFailedPrecondition", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.LeaseExpiresAt = time.Now().Add(-time.Second)
		r.updates(op)

		_, err := r.server().MarkRunning(context.Background(), request(op.Id))

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assert.Equal(t, []string{"rejected_running"}, r.events())
	})

	t.Run("ZeroTokenIsInvalid", func(t *testing.T) {
		r := newTestRepos()

		_, err := r.server().MarkRunning(context.Background(), &pb.MarkRunningRequest{
			Actor: "controller", Id: uuid.NewV4().String(),
		})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		r.assertUntouched(t)
	})

	t.Run("InvalidId", func(t *testing.T) {
		r := newTestRepos()

		_, err := r.server().MarkRunning(context.Background(), &pb.MarkRunningRequest{
			Actor: "controller", Id: "not-a-uuid", FencingToken: 3,
		})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		r.assertUntouched(t)
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
			r := newTestRepos()
			r.updateFails(id, tc.err)

			_, err := r.server().CompleteOperation(context.Background(),
				&pb.CompleteOperationRequest{
					Id: id.String(), FencingToken: 7, Actor: "a", Reason: "r",
				})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("Start"+name, func(t *testing.T) {
			r := newTestRepos()
			r.addFails(tc.err)

			_, err := r.server().StartOperation(context.Background(),
				startRequest("node:abc"))

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("StartIdempotencyLookup"+name, func(t *testing.T) {
			r := newTestRepos()
			r.groupLookupFails(tc.err)
			req := startRequest("node:abc")
			req.IdempotencyKey = "key"

			_, err := r.server().StartOperation(context.Background(), req)

			assert.Equal(t, tc.code, status.Code(err))
			assert.Empty(t, r.calls("Add"))
		})

		t.Run("GetOperation"+name, func(t *testing.T) {
			r := newTestRepos()
			r.ops.On("Get", mock.Anything, id).Return(nil, tc.err).Once()

			_, err := r.server().GetOperation(context.Background(),
				&pb.GetOperationRequest{Id: id.String()})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("GetByResourceLock"+name, func(t *testing.T) {
			r := newTestRepos()
			r.locks.On("Holder", mock.Anything, "node:abc").Return(nil, tc.err).Once()

			_, err := r.server().GetByResource(context.Background(),
				&pb.GetByResourceRequest{ResourceKey: "node:abc"})

			assert.Equal(t, tc.code, status.Code(err))
			assert.Empty(t, r.calls("Get"))
		})

		t.Run("GetByResourceHolder"+name, func(t *testing.T) {
			r := newTestRepos()
			r.locks.On("Holder", mock.Anything, "node:abc").
				Return(&db.ResourceLock{ResourceKey: "node:abc", OperationId: id}, nil).
				Once()
			r.ops.On("Get", mock.Anything, id).Return(nil, tc.err).Once()

			_, err := r.server().GetByResource(context.Background(),
				&pb.GetByResourceRequest{ResourceKey: "node:abc"})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("MarkRunning"+name, func(t *testing.T) {
			r := newTestRepos()
			r.updateFails(id, tc.err)

			_, err := r.server().MarkRunning(context.Background(), &pb.MarkRunningRequest{
				Actor: "controller", Id: id.String(), FencingToken: 3,
			})

			assert.Equal(t, tc.code, status.Code(err))
		})

		t.Run("ForceUnlock"+name, func(t *testing.T) {
			r := newTestRepos()
			r.ops.On("Get", mock.Anything, id).Return(nil, tc.err).Once()

			_, err := r.server().ForceUnlock(context.Background(), &pb.ForceUnlockRequest{
				Id: id.String(), Actor: "a", Reason: "r",
			})

			assert.Equal(t, tc.code, status.Code(err))
			assert.Empty(t, r.calls("Update"))
		})
	}
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
			r := newTestRepos()
			req := startRequest(key)
			req.IdempotencyKey = "key"

			_, err := r.server().StartOperation(context.Background(), req)

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			r.assertUntouched(t)
		})

		t.Run("GetByResource"+name, func(t *testing.T) {
			r := newTestRepos()

			_, err := r.server().GetByResource(context.Background(),
				&pb.GetByResourceRequest{ResourceKey: key})

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			r.assertUntouched(t)
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
			r := newTestRepos()
			r.adds(mock.Anything, tc.lease)
			req := startRequest("node:abc")
			req.LeaseSeconds = tc.seconds

			_, err := r.server().StartOperation(context.Background(), req)

			assert.NoError(t, err)
			ops := r.calls("Add")[0].Arguments.Get(1).([]*db.Operation)
			assert.WithinDuration(t, time.Now().Add(tc.lease), ops[0].LeaseExpiresAt,
				5*time.Second)
			r.assertExpectations(t)
		})
	}

	for name, seconds := range map[string]uint32{
		"AboveMaxIsRejected":  3601,
		"MaxUint32IsRejected": math.MaxUint32,
	} {
		t.Run(name, func(t *testing.T) {
			r := newTestRepos()
			req := startRequest("node:abc")
			req.IdempotencyKey = "key"
			req.LeaseSeconds = seconds

			_, err := r.server().StartOperation(context.Background(), req)

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			r.assertUntouched(t)
		})
	}
}

func TestBlankActorsAreRejected(t *testing.T) {
	id := uuid.NewV4().String()

	for _, blank := range []string{"", "   "} {
		calls := map[string]func(s *OperationServer) error{
			"StartRequestedBy": func(s *OperationServer) error {
				_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
					Type: "RestartNode", System: "node",
					ResourceKeys: []string{"node:abc"},
					RequestedBy:  blank,
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
				r := newTestRepos()

				err := call(r.server())

				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				r.assertUntouched(t)
			})
		}
	}
}

func TestDeleteOperation(t *testing.T) {
	request := func(id uuid.UUID) *pb.DeleteOperationRequest {
		return &pb.DeleteOperationRequest{Id: id.String(), Actor: "admin"}
	}

	t.Run("SoftDeletesFinishedOperationWithAudit", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(2)
		op.Status = db.OperationSuccess
		r.deletes(op)

		_, err := r.server().DeleteOperation(context.Background(), request(op.Id))

		assert.NoError(t, err)
		assert.Equal(t, []string{"deleted"}, r.events())
		assert.Equal(t, "admin", r.saved[0].Actor)
		assert.Equal(t, op.Id, r.saved[0].OperationId)
		assert.Equal(t, op.ResourceKey, r.saved[0].ResourceKey)
		r.assertExpectations(t)
	})

	for _, st := range []db.OperationStatus{db.OperationPending, db.OperationRunning} {
		t.Run("Rejects"+st.String(), func(t *testing.T) {
			r := newTestRepos()
			op := runningOp(2)
			op.Status = st
			r.deletes(op)

			_, err := r.server().DeleteOperation(context.Background(), request(op.Id))

			assert.Equal(t, codes.FailedPrecondition, status.Code(err))
			assert.Empty(t, r.saved)
		})
	}

	for name, req := range map[string]*pb.DeleteOperationRequest{
		"InvalidId":  {Id: "nope", Actor: "admin"},
		"BlankActor": {Id: uuid.NewV4().String(), Actor: "  "},
	} {
		t.Run(name+"IsInvalidArgument", func(t *testing.T) {
			r := newTestRepos()

			_, err := r.server().DeleteOperation(context.Background(), req)

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			r.assertUntouched(t)
		})
	}

	for name, tc := range map[string]struct {
		err  error
		code codes.Code
	}{
		"UnknownIsNotFound": {
			fmt.Errorf("operation: %w", store.ErrNotFound), codes.NotFound,
		},
		"LockTimeoutIsUnavailable": {
			fmt.Errorf("delete: %w", db.ErrLockTimeout), codes.Unavailable,
		},
		"OtherErrorIsInternal": {fmt.Errorf("db down"), codes.Internal},
	} {
		t.Run(name, func(t *testing.T) {
			r := newTestRepos()
			id := uuid.NewV4()
			r.ops.On("Delete", mock.Anything, id, mock.Anything).Return(tc.err).Once()

			_, err := r.server().DeleteOperation(context.Background(), request(id))

			assert.Equal(t, tc.code, status.Code(err))
		})
	}
}
