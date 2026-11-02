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

	"github.com/tj/assert"
	"google.golang.org/grpc"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

type fakeServiceClient struct {
	pb.OperationManagerServiceClient
	completeReq *pb.CompleteOperationRequest
	failReq     *pb.FailOperationRequest
}

func (f *fakeServiceClient) CompleteOperation(_ context.Context,
	in *pb.CompleteOperationRequest,
	_ ...grpc.CallOption) (*pb.CompleteOperationResponse, error) {
	f.completeReq = in

	return &pb.CompleteOperationResponse{Operation: &pb.Operation{Id: in.Id}}, nil
}

func (f *fakeServiceClient) FailOperation(_ context.Context, in *pb.FailOperationRequest,
	_ ...grpc.CallOption) (*pb.FailOperationResponse, error) {
	f.failReq = in

	return &pb.FailOperationResponse{Operation: &pb.Operation{Id: in.Id}}, nil
}

func TestManager_Complete(t *testing.T) {
	fake := &fakeServiceClient{}

	resp, err := NewManagerFromClient(fake).Complete("op-1", 7, "controller", "dispatched")

	assert.NoError(t, err)
	assert.Equal(t, "op-1", resp.Operation.Id)
	assert.Equal(t, "op-1", fake.completeReq.Id)
	assert.Equal(t, uint64(7), fake.completeReq.FencingToken)
	assert.Equal(t, "controller", fake.completeReq.Actor)
	assert.Equal(t, "dispatched", fake.completeReq.Reason)
}

func TestManager_Fail(t *testing.T) {
	fake := &fakeServiceClient{}

	resp, err := NewManagerFromClient(fake).Fail("op-1", 7, "controller", "publish failed")

	assert.NoError(t, err)
	assert.Equal(t, "op-1", resp.Operation.Id)
	assert.Equal(t, "op-1", fake.failReq.Id)
	assert.Equal(t, uint64(7), fake.failReq.FencingToken)
	assert.Equal(t, "controller", fake.failReq.Actor)
	assert.Equal(t, "publish failed", fake.failReq.Reason)
}
