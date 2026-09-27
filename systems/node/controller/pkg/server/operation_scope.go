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
	"strings"

	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *ControllerServer) acquireOperation(ctx context.Context, actionType, resourceKey string) (*copr.OperationInfo, error) {
	id, token, err := copr.ReservationFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if id != "" {
		op, err := c.opManager.Get(id)
		if err != nil {
			return nil, err
		}
		if op == nil || op.Id != id || op.ResourceKey != resourceKey ||
			op.Type != actionType || op.System != "node" || op.RequestedBy != "site_controller" ||
			op.Status != copr.StatusPending || op.FencingToken != token {
			return nil, status.Error(codes.FailedPrecondition, "invalid or already claimed site operation reservation")
		}
		// Claim once under the manager's operation-row lock. A repeated or
		// concurrent delivery must not execute or cancel the first command.
		running, err := c.opManager.MarkRunning(id, token)
		if err != nil {
			return nil, err
		}
		if running == nil || running.Id != id || running.Status != copr.StatusRunning {
			return nil, status.Error(codes.Internal, "operation manager did not confirm reservation claim")
		}
		return running, nil
	}
	conflicts, err := copr.NodeConflictResources(c.nodeClient, strings.TrimPrefix(resourceKey, "node:"))
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "resolve site operation scope: %v", err)
	}
	resp, err := c.opManager.Start(copr.StartRequest{
		Type: actionType, System: "node", ResourceKey: resourceKey,
		RequestedBy: pkg.ServiceName, LeaseSeconds: c.opLeaseSecs,
		ConflictResourceKeys: conflicts,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Operation == nil {
		return nil, status.Error(codes.Internal, "operation manager returned no reservation")
	}
	return resp.Operation, nil
}
