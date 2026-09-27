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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/mocks"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestScopedStartProtocolRoundTrip(t *testing.T) {
	req := &pb.StartOperationRequest{Type: "RestartNode", System: "node", ResourceKey: "node:t",
		ConflictResourceKeys: []string{"node:t", "node:a", "node:c"}, AdditionalResourceKeys: []string{"node:a"}}
	encoded, err := proto.Marshal(req)
	require.NoError(t, err)
	decoded := &pb.StartOperationRequest{}
	require.NoError(t, proto.Unmarshal(encoded, decoded))
	require.True(t, proto.Equal(req, decoded))
	resp := &pb.StartOperationResponse{ConflictsChecked: true,
		Operations: []*pb.Operation{{ResourceKey: "node:t"}, {ResourceKey: "node:a"}}}
	encoded, err = proto.Marshal(resp)
	require.NoError(t, err)
	got := &pb.StartOperationResponse{}
	require.NoError(t, proto.Unmarshal(encoded, got))
	require.True(t, proto.Equal(resp, got))
}

func TestScopedStartReturnsAllReservedTargets(t *testing.T) {
	repo := &mocks.OperationRepo{}
	ops := []*db.Operation{
		{Id: uuid.NewV4(), ResourceKey: "node:t"},
		{Id: uuid.NewV4(), ResourceKey: "node:a"},
	}
	repo.On("StartBatch", mock.MatchedBy(func(in []*db.Operation) bool {
		return len(in) == 2 && in[0].ResourceKey == "node:t" && in[1].ResourceKey == "node:a"
	}), []string{"node:c", "node:t", "node:a"}, mock.Anything).Return(ops, nil, nil).Once()
	s := NewOperationServer(orgName, "", repo, nil)
	resp, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
		Type: "RestartNode", System: "node", ResourceKey: "node:t",
		ConflictResourceKeys: []string{"node:c", "node:t", "node:a"}, AdditionalResourceKeys: []string{"node:a"},
	})
	require.NoError(t, err)
	assert.True(t, resp.ConflictsChecked)
	assert.Len(t, resp.Operations, 2)
	assert.Equal(t, ops[0].Id.String(), resp.Operation.Id)
	repo.AssertExpectations(t)
}

func TestScopedStartRejectsControllerConflict(t *testing.T) {
	repo := &mocks.OperationRepo{}
	holder := &db.Operation{Id: uuid.NewV4(), ResourceKey: "node:c", Type: "UpdateSoftware"}
	repo.On("StartBatch", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, holder, db.ErrLockConflict).Once()
	s := NewOperationServer(orgName, "", repo, nil)
	_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
		Type: "ToggleService", System: "node", ResourceKey: "node:t",
		ConflictResourceKeys: []string{"node:c", "node:t"},
	})
	require.Equal(t, codes.AlreadyExists, status.Code(err))
	repo.AssertExpectations(t)
}

func TestScopedStartRejectsTargetOutsideConflictSet(t *testing.T) {
	repo := &mocks.OperationRepo{}
	s := NewOperationServer(orgName, "", repo, nil)
	_, err := s.StartOperation(context.Background(), &pb.StartOperationRequest{
		Type: "RestartNode", System: "node", ResourceKey: "node:t",
		ConflictResourceKeys: []string{"node:c"},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Empty(t, repo.Calls)
}
