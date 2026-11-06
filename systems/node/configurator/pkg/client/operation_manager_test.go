/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

type fakeManagerClient struct {
	pb.OperationManagerServiceClient
	deadline bool
	start    *pb.StartOperationRequest
	running  *pb.MarkRunningRequest
}

func (f *fakeManagerClient) StartOperation(ctx context.Context,
	req *pb.StartOperationRequest, _ ...grpc.CallOption) (*pb.StartOperationResponse,
	error) {
	_, f.deadline = ctx.Deadline()
	f.start = req
	return &pb.StartOperationResponse{Operation: &pb.Operation{Id: "op"}}, nil
}

func (f *fakeManagerClient) MarkRunning(ctx context.Context,
	req *pb.MarkRunningRequest, _ ...grpc.CallOption) (*pb.MarkRunningResponse,
	error) {
	_, f.deadline = ctx.Deadline()
	f.running = req
	return &pb.MarkRunningResponse{Operation: &pb.Operation{Id: req.Id}}, nil
}

func TestOperationManager(t *testing.T) {
	t.Run("StartForwardsWithDeadline", func(t *testing.T) {
		fake := &fakeManagerClient{}
		m := &operationManager{client: fake, timeout: time.Second}

		resp, err := m.Start(&pb.StartOperationRequest{RequestedBy: "configurator"})

		assert.NoError(t, err)
		assert.Equal(t, "op", resp.Operation.Id)
		assert.Equal(t, "configurator", fake.start.RequestedBy)
		assert.True(t, fake.deadline)
	})

	t.Run("MarkRunningSendsActor", func(t *testing.T) {
		fake := &fakeManagerClient{}
		m := &operationManager{client: fake, timeout: time.Second}

		_, err := m.MarkRunning("op", 3, "configurator")

		assert.NoError(t, err)
		assert.Equal(t, "op", fake.running.Id)
		assert.Equal(t, uint64(3), fake.running.FencingToken)
		assert.Equal(t, "configurator", fake.running.Actor)
		assert.True(t, fake.deadline)
	})

	t.Run("CloseWithoutConnection", func(t *testing.T) {
		assert.NotPanics(t, (&operationManager{}).Close)
	})
}
