/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	cpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	cmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	spb "github.com/ukama/ukama/systems/node/state/pb/gen"
	smocks "github.com/ukama/ukama/systems/node/state/pb/gen/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConfigurationClientsPreserveIdentityAndContext(t *testing.T) {
	for _, cancelled := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "cancelled"}[cancelled], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if cancelled {
				cancel()
			}
			match := mock.MatchedBy(func(got context.Context) bool {
				deadline, ok := got.Deadline()
				return ok && time.Until(deadline) <= time.Second && (got.Err() != nil) == cancelled
			})
			var rpcErr error
			if cancelled {
				rpcErr = status.Error(codes.Canceled, "caller cancelled")
			}
			c := &cmocks.ControllerServiceClient{}
			controller := NewControllerFromClient(c)
			req := &cpb.ConfigNodeRequest{NodeId: "node", RequestId: "attempt", SiteId: "site", NetworkId: "network"}
			c.On("ConfigNode", match, req).Return(&cpb.ConfigNodeResponse{Status: "DISPATCHED"}, rpcErr).Once()
			response, err := controller.ConfigNode(ctx, req)
			require.Equal(t, rpcErr, err)
			require.Equal(t, "DISPATCHED", response.Status)
			deletion := &cpb.DeleteNodeConfigRequest{NodeId: "node", RequestId: "attempt", SiteId: "site", NetworkId: "network"}
			c.On("DeleteNodeConfig", match, deletion).Return(&cpb.DeleteNodeConfigResponse{Status: "DISPATCHED"}, rpcErr).Once()
			removed, err := controller.DeleteNodeConfig(ctx, deletion)
			require.Equal(t, rpcErr, err)
			require.Equal(t, "DISPATCHED", removed.Status)
			c.AssertExpectations(t)
			s := &smocks.StateServiceClient{}
			state := NewStateFromClient(s)
			query := &spb.GetLatestStateRequest{NodeId: "node", RequestId: "attempt"}
			expected := &spb.GetLatestStateResponse{Configuration: &spb.ConfigurationStatus{RequestId: "attempt", Completed: true}}
			s.On("GetLatestState", match, query).Return(expected, rpcErr).Once()
			latest, err := state.GetLatestState(ctx, query)
			require.Equal(t, rpcErr, err)
			require.Same(t, expected, latest)
			s.AssertExpectations(t)
		})
	}
}
