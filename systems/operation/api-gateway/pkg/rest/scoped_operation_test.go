/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestScopedStartPassesConflictSetAndReservations(t *testing.T) {
	manager := &fakeManager{startResp: &pb.StartOperationResponse{
		ConflictsChecked: true, Operations: []*pb.Operation{{ResourceKey: "node:t"}, {ResourceKey: "node:a"}},
	}}
	r := &Router{clients: &Clients{Manager: manager}}
	resp, err := r.postStartHandler(&gin.Context{}, &StartOperationRequest{
		Type: "RestartNode", System: "node", ResourceKey: "node:t",
		ConflictResourceKeys: []string{"node:t", "node:a", "node:c"}, AdditionalResourceKeys: []string{"node:a"},
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"node:t", "node:a", "node:c"}, manager.lastStartReq.ConflictResourceKeys)
	assert.Equal(t, []string{"node:a"}, manager.lastStartReq.AdditionalResourceKeys)
	assert.True(t, resp.ConflictsChecked)
	assert.Len(t, resp.Operations, 2)
}

func TestForceUnlockAcceptsServiceActor(t *testing.T) {
	manager := &fakeManager{}
	r := &Router{clients: &Clients{Manager: manager}}
	_, err := r.postForceUnlockHandler(&gin.Context{}, &ForceUnlockRequest{
		Id: "op", Actor: "site-controller", Reason: "not dispatched",
	})
	require.NoError(t, err)
	assert.Equal(t, "site-controller", manager.lastForceUnlockActor)
	_, err = r.postForceUnlockHandler(&gin.Context{}, &ForceUnlockRequest{Id: "op", Reason: "not dispatched"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
