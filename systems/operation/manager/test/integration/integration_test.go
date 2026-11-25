//go:build integration
// +build integration

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package integration

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tj/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/uuid"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

const (
	actor  = "integration-test"
	reason = "integration test"
)

var tConfig *TestConfig

type TestConfig struct {
	ServiceHost string `default:"localhost:9090"`
}

func init() {
	tConfig = &TestConfig{}

	if err := rconf.NewConfReader("integration").Read(tConfig); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
}

func newClient(t *testing.T) pb.OperationManagerServiceClient {
	t.Helper()

	conn, err := grpc.NewClient(tConfig.ServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, conn.Close())
	})

	return pb.NewOperationManagerServiceClient(conn)
}

func newContext(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	return ctx
}

func startOperation(t *testing.T, ctx context.Context,
	c pb.OperationManagerServiceClient) *pb.Operation {
	t.Helper()

	resp, err := c.StartOperation(ctx, &pb.StartOperationRequest{
		Type:         "IntegrationTest",
		System:       "operation",
		ResourceKeys: []string{"test:integration:" + uuid.NewV4().String()},
		RequestedBy:  actor,
	})
	assert.NoError(t, err)
	assert.Equal(t, pb.OperationStatus_PENDING, resp.Operations[0].Status)
	assert.True(t, resp.Operations[0].FencingToken > 0)
	assertLockHeldBy(t, ctx, c, resp.Operations[0])

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cur, err := c.GetOperation(ctx,
			&pb.GetOperationRequest{Id: resp.Operations[0].Id})
		if err != nil || !unfinished(cur.Operation.Status) {
			return
		}
		_, _ = c.ForceUnlock(ctx, &pb.ForceUnlockRequest{
			Id: resp.Operations[0].Id, Actor: actor, Reason: reason,
		})
	})

	return resp.Operations[0]
}

func unfinished(st pb.OperationStatus) bool {
	return st == pb.OperationStatus_PENDING || st == pb.OperationStatus_RUNNING
}

func startRunningOperation(t *testing.T, ctx context.Context,
	c pb.OperationManagerServiceClient) *pb.Operation {
	t.Helper()

	op := startOperation(t, ctx, c)
	resp, err := c.MarkRunning(ctx, &pb.MarkRunningRequest{
		Id: op.Id, FencingToken: op.FencingToken, Actor: actor,
	})
	assert.NoError(t, err)
	assert.Equal(t, pb.OperationStatus_RUNNING, resp.Operation.Status)

	return resp.Operation
}

func assertLockHeldBy(t *testing.T, ctx context.Context,
	c pb.OperationManagerServiceClient, op *pb.Operation) {
	t.Helper()

	resp, err := c.GetByResource(ctx, &pb.GetByResourceRequest{ResourceKey: op.ResourceKey})
	assert.NoError(t, err)
	assert.NotNil(t, resp.Operation)
	assert.Equal(t, op.Id, resp.Operation.Id)
}

func assertLockFree(t *testing.T, ctx context.Context,
	c pb.OperationManagerServiceClient, op *pb.Operation) {
	t.Helper()

	resp, err := c.GetByResource(ctx, &pb.GetByResourceRequest{ResourceKey: op.ResourceKey})
	assert.NoError(t, err)
	assert.Nil(t, resp.Operation)
}

func completeOperation(ctx context.Context, c pb.OperationManagerServiceClient,
	op *pb.Operation, token uint64) (*pb.CompleteOperationResponse, error) {
	return c.CompleteOperation(ctx, &pb.CompleteOperationRequest{
		Id: op.Id, FencingToken: token, Actor: actor, Reason: reason,
	})
}

func failOperation(ctx context.Context, c pb.OperationManagerServiceClient,
	op *pb.Operation, token uint64) (*pb.FailOperationResponse, error) {
	return c.FailOperation(ctx, &pb.FailOperationRequest{
		Id: op.Id, FencingToken: token, Actor: actor, Reason: reason,
	})
}

func TestComplete(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	t.Run("WrongTokenKeepsTheLock", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		_, err := completeOperation(ctx, c, op, op.FencingToken+1)

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assertLockHeldBy(t, ctx, c, op)
	})

	t.Run("RightTokenReleasesTheLock", func(t *testing.T) {
		op := startRunningOperation(t, ctx, c)

		resp, err := completeOperation(ctx, c, op, op.FencingToken)

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_SUCCESS, resp.Operation.Status)
		assertLockFree(t, ctx, c, op)
	})

	t.Run("OnPendingIsRejected", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		_, err := completeOperation(ctx, c, op, op.FencingToken)

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assertLockHeldBy(t, ctx, c, op)
	})

	t.Run("DuplicateIsIdempotent", func(t *testing.T) {
		op := startRunningOperation(t, ctx, c)

		_, err := completeOperation(ctx, c, op, op.FencingToken)
		assert.NoError(t, err)

		resp, err := completeOperation(ctx, c, op, op.FencingToken)

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_SUCCESS, resp.Operation.Status)
	})

	t.Run("AfterFailureIsRejected", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		_, err := failOperation(ctx, c, op, op.FencingToken)
		assert.NoError(t, err)

		_, err = completeOperation(ctx, c, op, op.FencingToken)

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))

		got, err := c.GetOperation(ctx, &pb.GetOperationRequest{Id: op.Id})
		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_FAILED, got.Operation.Status)
	})

	t.Run("UnknownOperation", func(t *testing.T) {
		_, err := completeOperation(ctx, c, &pb.Operation{Id: uuid.NewV4().String()}, 1)

		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}

func TestFail(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	t.Run("WrongTokenKeepsTheLock", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		_, err := failOperation(ctx, c, op, op.FencingToken+1)

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assertLockHeldBy(t, ctx, c, op)
	})

	t.Run("RightTokenReleasesTheLock", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		resp, err := failOperation(ctx, c, op, op.FencingToken)

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_FAILED, resp.Operation.Status)
		assert.Equal(t, reason, resp.Operation.Error)
		assertLockFree(t, ctx, c, op)
	})

	t.Run("AfterCompletionIsRejected", func(t *testing.T) {
		op := startRunningOperation(t, ctx, c)

		_, err := completeOperation(ctx, c, op, op.FencingToken)
		assert.NoError(t, err)

		_, err = failOperation(ctx, c, op, op.FencingToken)

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})
}

func TestForceUnlock(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	t.Run("ReleasesTheLockWithoutAToken", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		resp, err := c.ForceUnlock(ctx, &pb.ForceUnlockRequest{
			Id: op.Id, Actor: actor, Reason: reason,
		})

		assert.NoError(t, err)
		assert.Equal(t, pb.OperationStatus_CANCELLED, resp.Operation.Status)
		assertLockFree(t, ctx, c, op)
	})

	t.Run("FinishedOperationIsRejected", func(t *testing.T) {
		op := startRunningOperation(t, ctx, c)

		_, err := completeOperation(ctx, c, op, op.FencingToken)
		assert.NoError(t, err)

		_, err = c.ForceUnlock(ctx, &pb.ForceUnlockRequest{
			Id: op.Id, Actor: actor, Reason: reason,
		})

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})
}

func TestInvalidResourceKey(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	for name, key := range map[string]string{
		"Blank":   "   ",
		"Nul":     "test:integration:\x00",
		"TooLong": "test:integration:" + strings.Repeat("a", 255),
	} {
		t.Run("Start"+name, func(t *testing.T) {
			_, err := c.StartOperation(ctx, &pb.StartOperationRequest{
				Type:         "IntegrationTest",
				System:       "operation",
				ResourceKeys: []string{key},
				RequestedBy:  actor,
			})

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		})

		t.Run("GetByResource"+name, func(t *testing.T) {
			_, err := c.GetByResource(ctx, &pb.GetByResourceRequest{ResourceKey: key})

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}

func TestStartLease(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	start := func(seconds uint32) (*pb.StartOperationResponse, error) {
		return c.StartOperation(ctx, &pb.StartOperationRequest{
			Type:         "IntegrationTest",
			System:       "operation",
			ResourceKeys: []string{"test:integration:" + uuid.NewV4().String()},
			RequestedBy:  actor,
			LeaseSeconds: seconds,
		})
	}

	t.Run("MaxIsAccepted", func(t *testing.T) {
		resp, err := start(3600)

		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now().Add(time.Hour),
			resp.Operations[0].LeaseExpiresAt.AsTime(), time.Minute)

		_, err = c.ForceUnlock(ctx, &pb.ForceUnlockRequest{
			Id: resp.Operations[0].Id, Actor: actor, Reason: reason,
		})
		assert.NoError(t, err)
		assertLockFree(t, ctx, c, resp.Operations[0])
	})

	t.Run("AboveMaxIsRejected", func(t *testing.T) {
		_, err := start(3601)

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

func TestStartConflict(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	t.Run("SecondStartIsAlreadyExists", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		resp, err := c.StartOperation(ctx, &pb.StartOperationRequest{
			Type:         "IntegrationTest",
			System:       "operation",
			ResourceKeys: []string{op.ResourceKey},
			RequestedBy:  actor,
		})

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.Nil(t, resp)
		assertLockHeldBy(t, ctx, c, op)
	})
}

func startWithKey(t *testing.T, ctx context.Context,
	c pb.OperationManagerServiceClient, resourceKey,
	key string) (*pb.StartOperationResponse, error) {
	t.Helper()

	resp, err := c.StartOperation(ctx, &pb.StartOperationRequest{
		Type:           "IntegrationTest",
		System:         "operation",
		ResourceKeys:   []string{resourceKey},
		RequestedBy:    actor,
		IdempotencyKey: key,
	})
	if err == nil {
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, _ = c.ForceUnlock(ctx, &pb.ForceUnlockRequest{
				Id: resp.Operations[0].Id, Actor: actor, Reason: reason,
			})
		})
	}

	return resp, err
}

func TestStartIdempotency(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)
	resource := func() string { return "test:integration:" + uuid.NewV4().String() }

	t.Run("RetryReturnsTheSameOperation", func(t *testing.T) {
		key, res := uuid.NewV4().String(), resource()
		first, err := startWithKey(t, ctx, c, res, key)
		assert.NoError(t, err)

		retry, err := startWithKey(t, ctx, c, res, key)

		assert.NoError(t, err)
		assert.Equal(t, first.Operations[0].Id, retry.Operations[0].Id)
		assert.Equal(t, first.Operations[0].FencingToken,
			retry.Operations[0].FencingToken)
		assertLockHeldBy(t, ctx, c, first.Operations[0])
	})

	t.Run("RetryAfterFinishReturnsTheFinishedOperation", func(t *testing.T) {
		key, res := uuid.NewV4().String(), resource()
		first, err := startWithKey(t, ctx, c, res, key)
		assert.NoError(t, err)
		_, err = c.MarkRunning(ctx, &pb.MarkRunningRequest{
			Id: first.Operations[0].Id, FencingToken: first.Operations[0].FencingToken,
			Actor: actor,
		})
		assert.NoError(t, err)
		_, err = completeOperation(ctx, c, first.Operations[0],
			first.Operations[0].FencingToken)
		assert.NoError(t, err)

		retry, err := startWithKey(t, ctx, c, res, key)

		assert.NoError(t, err)
		assert.Equal(t, first.Operations[0].Id, retry.Operations[0].Id)
		assert.Equal(t, pb.OperationStatus_SUCCESS, retry.Operations[0].Status)
		assertLockFree(t, ctx, c, first.Operations[0])
	})

	t.Run("KeyOnAnotherResourceIsInvalidArgument", func(t *testing.T) {
		key := uuid.NewV4().String()
		first, err := startWithKey(t, ctx, c, resource(), key)
		assert.NoError(t, err)
		other := resource()

		resp, err := startWithKey(t, ctx, c, other, key)

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Nil(t, resp)
		assertLockHeldBy(t, ctx, c, first.Operations[0])
		assertLockFree(t, ctx, c, &pb.Operation{ResourceKey: other})
	})

	t.Run("ConcurrentStartsReturnOneOperation", func(t *testing.T) {
		const callers = 10
		key, res := uuid.NewV4().String(), resource()
		ids := make([]string, callers)
		errs := make([]error, callers)

		var wg sync.WaitGroup
		for i := 0; i < callers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				resp, err := startWithKey(t, ctx, c, res, key)
				errs[i] = err
				if err == nil {
					ids[i] = resp.Operations[0].Id
				}
			}(i)
		}
		wg.Wait()

		for i := 0; i < callers; i++ {
			assert.NoError(t, errs[i])
			assert.Equal(t, ids[0], ids[i])
		}
		assertLockHeldBy(t, ctx, c, &pb.Operation{Id: ids[0], ResourceKey: res})
	})
}

func TestBlankActors(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	t.Run("StartWithoutRequestedBy", func(t *testing.T) {
		_, err := c.StartOperation(ctx, &pb.StartOperationRequest{
			Type:         "IntegrationTest",
			System:       "operation",
			ResourceKeys: []string{"test:integration:" + uuid.NewV4().String()},
			RequestedBy:  "   ",
		})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("MarkRunningWithoutActor", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		_, err := c.MarkRunning(ctx, &pb.MarkRunningRequest{
			Id: op.Id, FencingToken: op.FencingToken,
		})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("CompleteWithBlankActor", func(t *testing.T) {
		op := startRunningOperation(t, ctx, c)

		_, err := c.CompleteOperation(ctx, &pb.CompleteOperationRequest{
			Id: op.Id, FencingToken: op.FencingToken, Actor: "  ", Reason: reason,
		})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assertLockHeldBy(t, ctx, c, op)
	})
}

func startGroup(t *testing.T, ctx context.Context, c pb.OperationManagerServiceClient,
	targets, mustBeFree []string, key string) (*pb.StartOperationResponse, error) {
	t.Helper()

	resp, err := c.StartOperation(ctx, &pb.StartOperationRequest{
		Type:           "IntegrationTest",
		System:         "operation",
		ResourceKeys:   targets,
		MustBeFree:     mustBeFree,
		RequestedBy:    actor,
		IdempotencyKey: key,
	})
	if err == nil {
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for _, op := range resp.Operations {
				_, _ = c.ForceUnlock(ctx, &pb.ForceUnlockRequest{
					Id: op.Id, Actor: actor, Reason: reason,
				})
			}
		})
	}

	return resp, err
}

func TestStartGroup(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)
	site := func() (string, string, string) {
		id := uuid.NewV4().String()
		return "test:integration:" + id + ":a", "test:integration:" + id + ":c",
			"test:integration:" + id + ":t"
	}

	t.Run("LocksEveryTargetAndLinksMembersToLeader", func(t *testing.T) {
		amp, ctl, tower := site()

		resp, err := startGroup(t, ctx, c, []string{tower, amp},
			[]string{amp, ctl, tower}, "")

		assert.NoError(t, err)
		assert.Len(t, resp.Operations, 2)
		leader, member := resp.Operations[0], resp.Operations[1]
		assert.Equal(t, amp, leader.ResourceKey)
		assert.Empty(t, leader.GroupId)
		assert.Equal(t, tower, member.ResourceKey)
		assert.Equal(t, leader.Id, member.GroupId)
		assertLockHeldBy(t, ctx, c, leader)
		assertLockHeldBy(t, ctx, c, member)
		assertLockFree(t, ctx, c, &pb.Operation{ResourceKey: ctl})
	})

	t.Run("BusySiblingRejectsTheWholeGroup", func(t *testing.T) {
		amp, ctl, tower := site()
		holder, err := startGroup(t, ctx, c, []string{ctl}, nil, "")
		assert.NoError(t, err)

		resp, err := startGroup(t, ctx, c, []string{tower, amp},
			[]string{amp, ctl, tower}, "")

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.Contains(t, status.Convert(err).Message(), holder.Operations[0].Id)
		assert.Nil(t, resp)
		assertLockFree(t, ctx, c, &pb.Operation{ResourceKey: tower})
		assertLockFree(t, ctx, c, &pb.Operation{ResourceKey: amp})
	})

	t.Run("FinishedSiblingDoesNotBlock", func(t *testing.T) {
		amp, ctl, tower := site()
		holder, err := startGroup(t, ctx, c, []string{ctl}, nil, "")
		assert.NoError(t, err)
		_, err = c.FailOperation(ctx, &pb.FailOperationRequest{
			Id: holder.Operations[0].Id, FencingToken: holder.Operations[0].FencingToken,
			Actor: actor, Reason: reason,
		})
		assert.NoError(t, err)

		_, err = startGroup(t, ctx, c, []string{tower, amp},
			[]string{amp, ctl, tower}, "")

		assert.NoError(t, err)
	})

	t.Run("RetryWithKeyReturnsTheWholeGroup", func(t *testing.T) {
		amp, ctl, tower := site()
		key := uuid.NewV4().String()
		first, err := startGroup(t, ctx, c, []string{tower, amp},
			[]string{amp, ctl, tower}, key)
		assert.NoError(t, err)

		retry, err := startGroup(t, ctx, c, []string{amp, tower},
			[]string{amp, ctl, tower}, key)

		assert.NoError(t, err)
		assert.Len(t, retry.Operations, 2)
		assert.Equal(t, first.Operations[0].Id, retry.Operations[0].Id)
		assert.Equal(t, first.Operations[1].Id, retry.Operations[1].Id)
	})

	t.Run("OverlappingConcurrentGroupsNeverBothWin", func(t *testing.T) {
		const rounds = 20
		for i := 0; i < rounds; i++ {
			amp, ctl, tower := site()
			all := []string{amp, ctl, tower}
			errs := make([]error, 2)

			var wg sync.WaitGroup
			for j, target := range []string{tower, amp} {
				wg.Add(1)
				go func(j int, target string) {
					defer wg.Done()
					_, errs[j] = startGroup(t, ctx, c, []string{target}, all, "")
				}(j, target)
			}
			wg.Wait()

			won := 0
			for _, err := range errs {
				if err == nil {
					won++
					continue
				}
				assert.Equal(t, codes.AlreadyExists, status.Code(err))
			}
			assert.Equal(t, 1, won)
		}
	})
}

func TestDeleteOperation(t *testing.T) {
	c := newClient(t)
	ctx := newContext(t)

	t.Run("ActiveOperationIsRejected", func(t *testing.T) {
		op := startOperation(t, ctx, c)

		_, err := c.DeleteOperation(ctx,
			&pb.DeleteOperationRequest{Id: op.Id, Actor: actor})

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assertLockHeldBy(t, ctx, c, op)
	})

	t.Run("FinishedOperationIsDeleted", func(t *testing.T) {
		op := startOperation(t, ctx, c)
		_, err := c.FailOperation(ctx, &pb.FailOperationRequest{
			Id: op.Id, FencingToken: op.FencingToken, Actor: actor, Reason: reason,
		})
		assert.NoError(t, err)

		_, err = c.DeleteOperation(ctx,
			&pb.DeleteOperationRequest{Id: op.Id, Actor: actor})

		assert.NoError(t, err)
		_, err = c.GetOperation(ctx, &pb.GetOperationRequest{Id: op.Id})
		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("UnknownOperationIsNotFound", func(t *testing.T) {
		_, err := c.DeleteOperation(ctx, &pb.DeleteOperationRequest{
			Id: uuid.NewV4().String(), Actor: actor,
		})

		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}
