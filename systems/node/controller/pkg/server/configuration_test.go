/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package server

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	spb "github.com/ukama/ukama/systems/node/state/pb/gen"
)

type configurationRecorder func(context.Context, *spb.RecordConfigurationRequest) (*spb.ConfigurationStatus, error)

func (f configurationRecorder) RecordConfiguration(ctx context.Context, req *spb.RecordConfigurationRequest) (*spb.ConfigurationStatus, error) {
	return f(ctx, req)
}

func TestConfigurationDispatchAfterRecording(t *testing.T) {
	for _, cancelled := range []bool{false, true} {
		t.Run(map[bool]string{false: "configure", true: "cancel"}[cancelled], func(t *testing.T) {
			recorded := false
			bus := &mbmocks.MsgBusServiceClient{}
			bus.On("PublishRequest", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				require.True(t, recorded)
				message := args.Get(1).(*epb.NodeFeederMessage)
				require.Equal(t, "configd/v1/config", message.Path)
				require.Equal(t, map[bool]string{false: "POST", true: "DELETE"}[cancelled], message.HttpMethod)
				var body map[string]string
				require.NoError(t, json.Unmarshal(message.Msg, &body))
				require.Equal(t, "attempt-1", body["requestId"])
				if !cancelled {
					require.Equal(t, "NOCONFIG", body["mode"])
				}
			}).Return(nil).Once()
			server := &ControllerServer{orgName: "org", msgbus: bus}
			server.SetConfigurationState(configurationRecorder(func(ctx context.Context, req *spb.RecordConfigurationRequest) (*spb.ConfigurationStatus, error) {
				require.Equal(t, cancelled, req.Cancelled)
				require.Equal(t, "site-1", req.SiteId)
				recorded = true
				return &spb.ConfigurationStatus{RequestId: req.RequestId}, nil
			}))
			node := "uk-983794-hnode-78-7830"
			if cancelled {
				out, err := server.DeleteNodeConfig(context.Background(), &pb.DeleteNodeConfigRequest{NodeId: node, RequestId: "attempt-1", SiteId: "site-1", NetworkId: "net-1"})
				require.NoError(t, err)
				require.Equal(t, "DISPATCHED", out.Status)
			} else {
				out, err := server.ConfigNode(context.Background(), &pb.ConfigNodeRequest{NodeId: node, RequestId: "attempt-1", SiteId: "site-1", NetworkId: "net-1"})
				require.NoError(t, err)
				require.Equal(t, "DISPATCHED", out.Status)
			}
			bus.AssertExpectations(t)
		})
	}
}

func TestConfigurationRecordFailurePreventsDispatch(t *testing.T) {
	bus := &mbmocks.MsgBusServiceClient{}
	server := &ControllerServer{msgbus: bus}
	server.SetConfigurationState(configurationRecorder(func(context.Context, *spb.RecordConfigurationRequest) (*spb.ConfigurationStatus, error) {
		return nil, errors.New("database unavailable")
	}))
	err := server.dispatchConfiguration(context.Background(), "node", "attempt", "site", "network", false)
	require.Error(t, err)
	bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
}
