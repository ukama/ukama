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
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/mocks"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
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
		r.adds([]string{})

		resp, err := r.server().StartOperation(context.Background(),
			startRequest("node:abc"))

		assert.NoError(t, err)
		assert.Len(t, resp.Operations, 1)
		assert.Equal(t, "node:abc", resp.Operations[0].ResourceKey)
		assert.Equal(t, uint64(1), resp.Operations[0].FencingToken)
		assert.Empty(t, resp.Operations[0].GroupId)
		assert.Equal(t, []string{pkg.AuditLockAcquired}, r.events())
		assert.Equal(t, "controller", r.saved[0].Actor)
		assert.Equal(t, "node:abc", r.saved[0].ResourceKey)
		r.assertExpectations(t)
	})

	t.Run("GroupSortsTargetsAndLinksMembersToLeader", func(t *testing.T) {
		r := newTestRepos()
		r.findsGroup(nil)
		r.adds([]string{"node:c", "node:t"})
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
		assert.Equal(t, []string{pkg.AuditLockAcquired, pkg.AuditLockAcquired}, r.events())
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
		r.locks.On("Acquire", mock.Anything, mock.Anything, mock.Anything).
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
		r.locks.On("Acquire", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()
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
		assert.Equal(t, []string{pkg.AuditCompleted}, r.events())
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
		assert.Equal(t, []string{pkg.AuditRejected + pkg.AuditCompleted}, r.events())
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
		assert.Equal(t, []string{pkg.AuditFailed}, r.events())
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
		assert.Equal(t, []string{pkg.AuditIgnored + pkg.AuditFailed}, r.events())
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
		assert.Equal(t, []string{pkg.AuditForceUnlock}, r.events())
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
		assert.Equal(t, []string{pkg.AuditRunning}, r.events())
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
		assert.Equal(t, []string{pkg.AuditRejected + pkg.AuditRunning}, r.events())
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
		"ZeroUsesDefault": {0, defaultLeaseTTL},
		"MaxIsAccepted":   {3600, time.Hour},
	} {
		t.Run(name, func(t *testing.T) {
			r := newTestRepos()
			r.adds(mock.Anything)
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
		assert.Equal(t, []string{pkg.AuditDeleted}, r.events())
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

func TestTerminalPathsPublishResult(t *testing.T) {
	t.Run("CompleteOperation", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(7)
		r.updates(op)

		_, err := r.server().CompleteOperation(context.Background(),
			&pb.CompleteOperationRequest{Id: op.Id.String(), FencingToken: 7,
				Actor: "controller", Reason: "dispatched"})

		require.NoError(t, err)
		require.Len(t, r.published, 1)
		assert.Equal(t, db.OperationSuccess, r.published[0].Status)
		assert.Equal(t, "controller", r.results[0].Actor)
	})

	t.Run("FailOperation", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(7)
		r.updates(op)

		_, err := r.server().FailOperation(context.Background(),
			&pb.FailOperationRequest{Id: op.Id.String(), FencingToken: 7,
				Actor: "software", Reason: "publish failed"})

		require.NoError(t, err)
		require.Len(t, r.published, 1)
		assert.Equal(t, db.OperationFailed, r.published[0].Status)
		assert.Equal(t, "publish failed", r.published[0].Error)
	})

	t.Run("ForceUnlock", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(7)
		r.ops.On("Get", mock.Anything, op.Id).Return(&op, nil).Once()
		r.updates(op)

		_, err := r.server().ForceUnlock(context.Background(), &pb.ForceUnlockRequest{
			Id: op.Id.String(), Actor: "admin", Reason: "stuck"})

		require.NoError(t, err)
		require.Len(t, r.published, 1)
		assert.Equal(t, db.OperationCancelled, r.published[0].Status)
		assert.Equal(t, "admin", r.results[0].Actor)
	})

	t.Run("Sweeper", func(t *testing.T) {
		r := newTestRepos()
		op := runningOp(4)
		r.ops.On("List", mock.Anything, mock.Anything).
			Return([]*db.Operation{&op}, nil).Once()
		r.updates(op)

		NewSweeper(r.repos()).sweepOnce(context.Background())

		require.Len(t, r.published, 1)
		assert.Equal(t, db.OperationTimeout, r.published[0].Status)
		assert.Equal(t, operationManagerSweeperActor, r.results[0].Actor)
	})
}

const orgName = "ukama"

type testRepos struct {
	ops       *mocks.OperationRepo
	locks     *mocks.LockRepo
	audits    *mocks.AuditRepo
	saved     []db.OperationAudit
	released  []uuid.UUID
	published []db.Operation
	results   []db.OperationAudit
}

func (r *testRepos) Publish(op *db.Operation, audit db.OperationAudit) {
	r.published = append(r.published, *op)
	r.results = append(r.results, audit)
}

func newTestRepos() *testRepos {
	return &testRepos{
		ops: &mocks.OperationRepo{}, locks: &mocks.LockRepo{}, audits: &mocks.AuditRepo{},
	}
}

func (r *testRepos) repos() Repos {
	return Repos{Operations: r.ops, Locks: r.locks, Audits: r.audits, Results: r}
}

func (r *testRepos) server() *OperationServer {
	return NewOperationServer(orgName, "", r.repos())
}

func (r *testRepos) assertExpectations(t *testing.T) {
	t.Helper()
	r.ops.AssertExpectations(t)
	r.locks.AssertExpectations(t)
	r.audits.AssertExpectations(t)
}

func (r *testRepos) assertUntouched(t *testing.T) {
	t.Helper()
	assert.Empty(t, r.ops.Calls)
	assert.Empty(t, r.locks.Calls)
	assert.Empty(t, r.audits.Calls)
}

func (r *testRepos) recordAudits() *testRepos {
	record := func(args mock.Arguments) {
		r.saved = append(r.saved, *args.Get(1).(*db.OperationAudit))
	}
	r.audits.On("AddInTx", mock.Anything, mock.Anything).Run(record).Return(nil).Maybe()
	r.audits.On("Add", mock.Anything, mock.Anything).Run(record).Return(nil).Maybe()
	return r
}

func (r *testRepos) events() []string {
	out := make([]string, 0, len(r.saved))
	for _, a := range r.saved {
		out = append(out, a.Event)
	}
	return out
}

func (r *testRepos) calls(method string) []mock.Call {
	var out []mock.Call
	for _, c := range r.ops.Calls {
		if c.Method == method {
			out = append(out, c)
		}
	}
	return out
}

func (r *testRepos) adds(mustBeFree interface{}) {
	r.ops.On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(func(_ context.Context, ops []*db.Operation,
			nested func(*gorm.DB) error) error {
			for i, op := range ops {
				op.FencingToken = uint64(i + 1)
			}
			return nested(nil)
		}).Once()
	r.locks.On("Acquire", mock.Anything, mock.Anything, mustBeFree).Return(nil).Once()
	r.recordAudits()
}

func (r *testRepos) addRunsNested() {
	r.ops.On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(func(_ context.Context, _ []*db.Operation,
			nested func(*gorm.DB) error) error {
			return nested(nil)
		}).Once()
}

func (r *testRepos) addFails(err error) {
	r.ops.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(err).Once()
}

func leaderFilter(key string) interface{} {
	return mock.MatchedBy(func(f db.OperationFilter) bool {
		return f.IdempotencyKey == key && f.GroupId == nil
	})
}

func membersFilter(leader uuid.UUID) interface{} {
	return mock.MatchedBy(func(f db.OperationFilter) bool {
		return f.GroupId != nil && *f.GroupId == leader &&
			f.OrderBy == db.OrderByResourceKey && f.IdempotencyKey == ""
	})
}

func (r *testRepos) findsGroup(group []*db.Operation) {
	if len(group) == 0 {
		r.ops.On("List", mock.Anything, leaderFilter("key")).Return(nil, nil).Once()
		return
	}
	r.ops.On("List", mock.Anything, leaderFilter("key")).Return(group[:1], nil).Once()
	r.ops.On("List", mock.Anything, membersFilter(group[0].Id)).
		Return(group[1:], nil).Once()
}

func (r *testRepos) groupLookupFails(err error) {
	r.ops.On("List", mock.Anything, leaderFilter("key")).Return(nil, err).Once()
}

func (r *testRepos) updates(current db.Operation) {
	r.ops.On("Update", mock.Anything, current.Id, mock.Anything).
		Return(func(_ context.Context, _ uuid.UUID,
			nested func(*db.Operation, *gorm.DB) error) (*db.Operation, error) {
			op := current
			if err := nested(&op, nil); err != nil {
				return nil, err
			}
			return &op, nil
		}).Once()
	r.locks.On("Release", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r.released = append(r.released, args.Get(1).(*db.Operation).Id)
		}).Return(nil).Maybe()
	r.recordAudits()
}

func (r *testRepos) updateFails(id uuid.UUID, err error) {
	r.ops.On("Update", mock.Anything, id, mock.Anything).Return(nil, err).Once()
}

func (r *testRepos) deletes(current db.Operation) {
	r.ops.On("Delete", mock.Anything, current.Id, mock.Anything).
		Return(func(_ context.Context, _ uuid.UUID,
			nested func(*db.Operation, *gorm.DB) error) error {
			op := current
			return nested(&op, nil)
		}).Once()
	r.recordAudits()
}

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
		assert.Equal(t, []string{pkg.AuditRunning}, r.events())
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
		assert.Equal(t, []string{pkg.AuditIgnored + pkg.AuditRunning}, r.events())
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
			assert.Equal(t, []string{pkg.AuditRejected + pkg.AuditRunning}, r.events())
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
	completed := db.OperationAudit{Event: pkg.AuditCompleted, Actor: "controller",
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
		assert.Equal(t, []string{pkg.AuditCompleted}, r.events())
		assert.Equal(t, op.Id, r.saved[0].OperationId)
		assert.Equal(t, "dispatched", r.saved[0].Reason)
		assert.Equal(t, *out.TerminalAt, r.saved[0].At)
	})

	t.Run("FailOnPendingRecordsTheError", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationFailed,
			db.OperationAudit{Event: pkg.AuditFailed}, "boom")

		assert.NoError(t, err)
		assert.Equal(t, db.OperationFailed, out.Status)
		assert.Equal(t, "boom", out.Error)
		assert.Equal(t, []string{pkg.AuditFailed}, r.events())
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
		assert.Equal(t, []string{pkg.AuditRejected + pkg.AuditCompleted}, r.events())
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
		assert.Equal(t, []string{pkg.AuditIgnored + pkg.AuditCompleted}, r.events())
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
		assert.Equal(t, []string{pkg.AuditIgnored + pkg.AuditCompleted}, r.events())
	})

	t.Run("RejectsCompleteOnPending", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		r.updates(op)

		out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			completed, "")

		assert.True(t, errors.Is(err, store.ErrPrecondition), err)
		assert.Nil(t, out)
		assert.Equal(t, []string{pkg.AuditRejected + pkg.AuditCompleted}, r.events())
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
			db.OperationAudit{Event: pkg.AuditRunning}, "")

		assert.Error(t, err)
		assert.Nil(t, out)
		r.assertUntouched(t)
	})
}

func TestRefusalError(t *testing.T) {
	op := pendingOp()

	assert.Equal(t, pkg.AuditIgnored+"transition", ignore(&op, nil).Error())
	assert.Equal(t, "boom", reject(&op, errors.New("boom")).Error())
}

func TestTerminatePublishesResult(t *testing.T) {
	for _, st := range []db.OperationStatus{db.OperationSuccess, db.OperationFailed,
		db.OperationTimeout, db.OperationCancelled} {
		t.Run(st.String(), func(t *testing.T) {
			r := newTestRepos()
			op := pendingOp()
			op.Status = db.OperationRunning
			r.updates(op)

			out, err := r.repos().terminate(ctx, op.Id, op.FencingToken, st,
				db.OperationAudit{Event: "end", Actor: "controller"}, "reason")

			assert.NoError(t, err)
			assert.Len(t, r.published, 1)
			assert.Equal(t, st, r.published[0].Status)
			assert.Equal(t, "controller", r.results[0].Actor)
			assert.Equal(t, *out.TerminalAt, r.results[0].At)
		})
	}

	t.Run("RefusalPublishesNothing", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationRunning
		op.FencingToken = 9
		r.updates(op)

		_, err := r.repos().terminate(ctx, op.Id, 4, db.OperationSuccess,
			db.OperationAudit{Event: pkg.AuditCompleted}, "")

		assert.Error(t, err)
		assert.Empty(t, r.published)
	})

	t.Run("DuplicatePublishesNothing", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationSuccess
		r.updates(op)

		_, err := r.repos().terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			db.OperationAudit{Event: pkg.AuditCompleted}, "")

		assert.NoError(t, err)
		assert.Empty(t, r.published)
	})

	t.Run("NoPublisherIsFine", func(t *testing.T) {
		r := newTestRepos()
		op := pendingOp()
		op.Status = db.OperationRunning
		r.updates(op)
		repos := r.repos()
		repos.Results = nil

		_, err := repos.terminate(ctx, op.Id, op.FencingToken, db.OperationSuccess,
			db.OperationAudit{Event: pkg.AuditCompleted}, "")

		assert.NoError(t, err)
	})
}

func TestEndTime(t *testing.T) {
	created := time.Now().UTC().Add(-time.Hour)
	op := &db.Operation{CreatedAt: created}
	observed := created.Add(time.Minute)

	assert.Equal(t, observed, endTime(op, observed))
	for name, at := range map[string]time.Time{
		"Zero":          {},
		"Future":        time.Now().Add(time.Hour),
		"BeforeCreated": created.Add(-time.Minute),
	} {
		t.Run(name+"UsesNow", func(t *testing.T) {
			assert.WithinDuration(t, time.Now(), endTime(op, at), time.Second)
		})
	}
}

type fakeBus struct {
	mb.MsgBusServiceClient
	routes []string
	msgs   []protoreflect.ProtoMessage
	err    error
}

func (b *fakeBus) PublishRequest(route string, msg protoreflect.ProtoMessage) error {
	b.routes = append(b.routes, route)
	b.msgs = append(b.msgs, msg)
	return b.err
}

func TestResultPublisher(t *testing.T) {
	at := time.Now().UTC().Truncate(time.Second)
	op := func(st db.OperationStatus) *db.Operation {
		o := pendingOp()
		o.Status = st
		o.Error = pkg.ReasonLeaseExpired
		return &o
	}
	audit := db.OperationAudit{Actor: operationManagerSweeperActor, At: at}
	for name, tc := range map[string]struct {
		status db.OperationStatus
		route  string
		check  func(t *testing.T, msg protoreflect.ProtoMessage)
	}{
		"Completed": {db.OperationSuccess, completedRoute,
			func(t *testing.T, msg protoreflect.ProtoMessage) {
				m := msg.(*epb.OperationCompletedEvent)
				assert.Equal(t, at, m.CompletedAt.AsTime())
				assert.Equal(t, audit.Actor, m.Actor)
			}},
		"Failed": {db.OperationFailed, failedRoute,
			func(t *testing.T, msg protoreflect.ProtoMessage) {
				m := msg.(*epb.OperationFailedEvent)
				assert.Equal(t, pkg.ReasonLeaseExpired, m.Reason)
				assert.Equal(t, at, m.FailedAt.AsTime())
			}},
		"Timeout": {db.OperationTimeout, timeoutRoute,
			func(t *testing.T, msg protoreflect.ProtoMessage) {
				m := msg.(*epb.OperationTimedOutEvent)
				assert.Equal(t, pkg.ReasonLeaseExpired, m.Reason)
				assert.Equal(t, at, m.TimedOutAt.AsTime())
			}},
		"Cancelled": {db.OperationCancelled, cancelledRoute,
			func(t *testing.T, msg protoreflect.ProtoMessage) {
				m := msg.(*epb.OperationCancelledEvent)
				assert.Equal(t, pkg.ReasonLeaseExpired, m.Reason)
				assert.Equal(t, at, m.CancelledAt.AsTime())
			}},
	} {
		t.Run(name, func(t *testing.T) {
			bus := &fakeBus{}
			o := op(tc.status)

			NewResultPublisher(orgName, bus).Publish(o, audit)

			assert.Equal(t, []string{msgbus.PrepareRoute(orgName,
				tc.route)}, bus.routes)
			tc.check(t, bus.msgs[0])
		})
	}

	t.Run("PublishErrorIsLogged", func(t *testing.T) {
		bus := &fakeBus{err: errors.New("bus down")}

		assert.NotPanics(t, func() {
			NewResultPublisher(orgName, bus).Publish(op(db.OperationSuccess), audit)
		})
		assert.Len(t, bus.routes, 1)
	})

	t.Run("NonTerminalPublishesNothing", func(t *testing.T) {
		bus := &fakeBus{}

		NewResultPublisher(orgName, bus).Publish(op(db.OperationRunning), audit)

		assert.Empty(t, bus.routes)
	})
}

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
		assert.Equal(t, []string{pkg.AuditTimeout}, r.events())
		assert.Equal(t, pkg.ReasonLeaseExpired, r.saved[0].Reason)
		assert.Equal(t, operationManagerSweeperActor, r.saved[0].Actor)
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
