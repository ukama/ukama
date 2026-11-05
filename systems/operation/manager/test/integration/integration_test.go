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
	"testing"
	"time"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	"github.com/tj/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/uuid"
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
		Type:        "IntegrationTest",
		System:      "operation",
		ResourceKey: "test:integration:" + uuid.NewV4().String(),
		RequestedBy: actor,
	})
	assert.NoError(t, err)
	assert.Equal(t, pb.OperationStatus_PENDING, resp.Operation.Status)
	assert.True(t, resp.Operation.FencingToken > 0)
	assertLockHeldBy(t, ctx, c, resp.Operation)

	return resp.Operation
}

func startRunningOperation(t *testing.T, ctx context.Context,
	c pb.OperationManagerServiceClient) *pb.Operation {
	t.Helper()

	op := startOperation(t, ctx, c)
	resp, err := c.MarkRunning(ctx, &pb.MarkRunningRequest{
		Id: op.Id, FencingToken: op.FencingToken,
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
				Type:        "IntegrationTest",
				System:      "operation",
				ResourceKey: key,
				RequestedBy: actor,
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
			ResourceKey:  "test:integration:" + uuid.NewV4().String(),
			RequestedBy:  actor,
			LeaseSeconds: seconds,
		})
	}

	t.Run("MaxIsAccepted", func(t *testing.T) {
		resp, err := start(3600)

		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now().Add(time.Hour),
			resp.Operation.LeaseExpiresAt.AsTime(), time.Minute)
	})

	t.Run("AboveMaxIsRejected", func(t *testing.T) {
		_, err := start(3601)

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}
