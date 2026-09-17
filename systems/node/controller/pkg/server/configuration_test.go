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
			bus.On("PublishRequestWithContext", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				require.True(t, recorded)
				message := args.Get(2).(*epb.NodeFeederMessage)
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

func TestConfigurationDispatchFailuresAndReplay(t *testing.T) {
	for _, test := range []struct {
		name             string
		cancel           bool
		receipt          *spb.ConfigurationStatus
		recordErr        error
		cancelledContext bool
		noState          bool
		noBus            bool
		publishErr       error
		wantPublish      bool
		wantErr          bool
	}{
		{name: "missing state", noState: true, wantErr: true},
		{name: "missing receipt", wantErr: true},
		{name: "record failure", recordErr: errors.New("write failed"), wantErr: true},
		{name: "completed replay", receipt: &spb.ConfigurationStatus{Completed: true}},
		{name: "cleared replay", cancel: true, receipt: &spb.ConfigurationStatus{Cleared: true}},
		{name: "cancelled context", receipt: &spb.ConfigurationStatus{}, cancelledContext: true, wantErr: true},
		{name: "missing bus", receipt: &spb.ConfigurationStatus{}, noBus: true, wantErr: true},
		{name: "dispatch failure", receipt: &spb.ConfigurationStatus{}, publishErr: errors.New("broker down"), wantPublish: true, wantErr: true},
		{name: "cancel dispatch failure", cancel: true, receipt: &spb.ConfigurationStatus{}, publishErr: errors.New("broker down"), wantPublish: true, wantErr: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			bus := &mbmocks.MsgBusServiceClient{}
			server := &ControllerServer{orgName: "org", msgbus: bus}
			if test.noBus {
				server.msgbus = nil
			}
			if !test.noState {
				server.SetConfigurationState(configurationRecorder(func(context.Context, *spb.RecordConfigurationRequest) (*spb.ConfigurationStatus, error) {
					return test.receipt, test.recordErr
				}))
			}
			if test.wantPublish {
				bus.On("PublishRequestWithContext", mock.Anything, mock.Anything, mock.Anything).Return(test.publishErr).Once()
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if test.cancelledContext {
				cancel()
			}
			node := "uk-983794-hnode-78-7830"
			var err error
			if test.cancel {
				_, err = server.DeleteNodeConfig(ctx, &pb.DeleteNodeConfigRequest{NodeId: node, RequestId: "attempt", SiteId: "site", NetworkId: "network"})
			} else {
				_, err = server.ConfigNode(ctx, &pb.ConfigNodeRequest{NodeId: node, RequestId: "attempt", SiteId: "site", NetworkId: "network"})
			}
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			bus.AssertExpectations(t)
			if !test.wantPublish {
				bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestConfigurationRejectsInvalidNode(t *testing.T) {
	server := &ControllerServer{}
	_, err := server.ConfigNode(context.Background(), &pb.ConfigNodeRequest{NodeId: "bad"})
	require.Error(t, err)
	_, err = server.DeleteNodeConfig(context.Background(), &pb.DeleteNodeConfigRequest{NodeId: "bad"})
	require.Error(t, err)
}
