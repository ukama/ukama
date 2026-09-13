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
	"fmt"

	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	spb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConfigurationState interface {
	RecordConfiguration(context.Context, *spb.RecordConfigurationRequest) (*spb.ConfigurationStatus, error)
}

func (c *ControllerServer) SetConfigurationState(state ConfigurationState) {
	c.configurationState = state
}

func (c *ControllerServer) ConfigNode(ctx context.Context, req *pb.ConfigNodeRequest) (*pb.ConfigNodeResponse, error) {
	node, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid node ID")
	}
	err = c.dispatchConfiguration(ctx, node.String(), req.RequestId, req.SiteId, req.NetworkId, false)
	if err != nil {
		return nil, err
	}
	// This is a dispatch acknowledgement. Registry waits for lifecycle evidence.
	return &pb.ConfigNodeResponse{ResourceKey: nodeKey(node.String()), Status: "DISPATCHED"}, nil
}

func (c *ControllerServer) DeleteNodeConfig(ctx context.Context, req *pb.DeleteNodeConfigRequest) (*pb.DeleteNodeConfigResponse, error) {
	node, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid node ID")
	}
	err = c.dispatchConfiguration(ctx, node.String(), req.RequestId, req.SiteId, req.NetworkId, true)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteNodeConfigResponse{ResourceKey: nodeKey(node.String()), Status: "DISPATCHED"}, nil
}

func (c *ControllerServer) dispatchConfiguration(ctx context.Context, nodeID, requestID, siteID, networkID string, cancelled bool) error {
	if c.configurationState == nil {
		return status.Error(codes.Unavailable, "node-state unavailable")
	}
	// A durable record must exist before a command can leave controller. This
	// also stops assignment retries before DELETE, including offline cleanup.
	receipt, err := c.configurationState.RecordConfiguration(ctx, &spb.RecordConfigurationRequest{
		NodeId: nodeID, RequestId: requestID, SiteId: siteID, NetworkId: networkID, Cancelled: cancelled,
	})
	if err != nil {
		return err
	}
	if receipt == nil {
		return status.Error(codes.Unavailable, "missing configuration receipt")
	}
	if cancelled && receipt.Cleared || !cancelled && receipt.Completed {
		return nil
	}
	body := map[string]string{"requestId": requestID}
	action := "DELETE_CONFIG"
	if !cancelled {
		body["mode"] = "NOCONFIG"
		action = "CONFIG"
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return status.FromContextError(err).Err()
	}
	if c.msgbus == nil {
		return status.Error(codes.Unavailable, "message bus unavailable")
	}
	if err := c.publishMessage(fmt.Sprintf("%s...%s", c.orgName, nodeID), actions[action].method, actions[action].path, nodeID, data); err != nil {
		return status.Errorf(codes.Unavailable, "configuration dispatch: %v", err)
	}
	return nil
}
