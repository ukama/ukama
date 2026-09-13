/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package server

import (
	"context"

	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecordConfiguration is internal bookkeeping for controller. It records
// correlation and cancellation before dispatch; it never sends node commands.
func (s *StateServer) RecordConfiguration(ctx context.Context, req *pb.RecordConfigurationRequest) (*pb.ConfigurationStatus, error) {
	id, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil || !routingPart(req.SiteId) || !routingPart(req.NetworkId) ||
		!routingPart(req.RequestId) || len(req.RequestId) > 95 {
		return nil, status.Error(codes.InvalidArgument, "invalid configuration identity")
	}
	if s.configurationEvents == nil || s.configurationEvents.lifecycleRepo == nil {
		return nil, status.Error(codes.Unavailable, "configuration tracking unavailable")
	}
	out := &pb.ConfigurationStatus{RequestId: req.RequestId}
	err = s.configurationEvents.lifecycleRepo.Update(ctx, id.String(), "configuration-recorded",
		func(record *db.LifecycleRecord, state *db.State) error {
			if record.Attempts == nil {
				record.Attempts = make(map[string]db.ProvisionAttempt)
			}
			attempt, exists := record.Attempts[req.RequestId]
			if exists && (attempt.Site != req.SiteId || attempt.Network != req.NetworkId) {
				return status.Error(codes.AlreadyExists, "request belongs to another assignment")
			}
			if !exists {
				previous := record.Attempts[record.RequestID]
				if record.RequestID != "" && !previous.Cleared {
					if req.Cancelled {
						// This controller could not have dispatched an unregistered request.
						record.Attempts[req.RequestId] = db.ProvisionAttempt{Site: req.SiteId, Network: req.NetworkId, Cancelled: true, Cleared: true}
						out.Cancelled, out.Cleared = true, true
						return nil
					}
					return status.Error(codes.FailedPrecondition, "node already assigned or cleanup pending")
				}
				attempt = db.ProvisionAttempt{Site: req.SiteId, Network: req.NetworkId}
				record.RequestID, record.Site, record.Network = req.RequestId, req.SiteId, req.NetworkId
				record.Completed = false
				record.ProvisionManaged = true
				record.AwaitingOperational = false
			}
			if req.Cancelled {
				attempt.Cancelled = true
				if record.RequestID == req.RequestId {
					record.AwaitingOperational = false
					record.ProvisionManaged = true
				}
			} else if attempt.Cancelled || record.RequestID != req.RequestId {
				return status.Error(codes.FailedPrecondition, "configuration request cancelled")
			}
			record.Attempts[req.RequestId] = attempt
			out.Completed = attempt.Completed && !attempt.Cancelled &&
				record.RequestID == req.RequestId && state.CurrentState == npb.NodeState_Operational
			out.Cancelled, out.Cleared = attempt.Cancelled, attempt.Cleared
			return nil
		})
	return out, err
}

func (s *StateServer) configurationStatus(ctx context.Context, nodeID, requestID string) (*pb.ConfigurationStatus, error) {
	if !routingPart(requestID) || len(requestID) > 95 {
		return nil, status.Error(codes.InvalidArgument, "invalid request ID")
	}
	if s.configurationEvents == nil || s.configurationEvents.lifecycleRepo == nil {
		return nil, status.Error(codes.Unavailable, "configuration tracking unavailable")
	}
	record, state, err := s.configurationEvents.lifecycleRepo.ReadConfiguration(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	attempt, found := record.Attempts[requestID]
	if !found {
		return nil, status.Error(codes.NotFound, "unknown configuration request")
	}
	return &pb.ConfigurationStatus{RequestId: requestID, Cancelled: attempt.Cancelled, Cleared: attempt.Cleared,
		Completed: attempt.Completed && !attempt.Cancelled && record.RequestID == requestID && state.CurrentState == npb.NodeState_Operational}, nil
}
