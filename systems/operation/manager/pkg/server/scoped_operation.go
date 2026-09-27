/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"errors"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *OperationServer) startScopedOperation(req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	if req.Type == "" || req.System == "" || req.IdempotencyKey != "" ||
		len(req.ConflictResourceKeys) == 0 || len(req.ConflictResourceKeys) > 256 || len(req.AdditionalResourceKeys) > 255 {
		return nil, status.Error(codes.InvalidArgument, "invalid scoped operation request; idempotency keys are not supported for grouped admission")
	}
	conflicts := make(map[string]bool, len(req.ConflictResourceKeys))
	for _, key := range req.ConflictResourceKeys {
		if key == "" {
			return nil, status.Error(codes.InvalidArgument, "empty conflict resource key")
		}
		conflicts[key] = true
	}
	lease := time.Duration(req.LeaseSeconds) * time.Second
	if lease == 0 {
		lease = pkg.DefaultLeaseTTL
	}
	keys := append([]string{req.ResourceKey}, req.AdditionalResourceKeys...)
	ops := make([]*db.Operation, 0, len(keys))
	seen := make(map[string]bool, len(keys))
	for _, key := range keys {
		if key == "" || seen[key] || !conflicts[key] {
			return nil, status.Error(codes.InvalidArgument, "operation targets must be unique members of the conflict set")
		}
		seen[key] = true
		ops = append(ops, &db.Operation{
			Id: uuid.NewV4(), Type: req.Type, System: req.System,
			Status: db.OperationPending, RequestedBy: req.RequestedBy,
			ResourceKey: key, LeaseExpiresAt: time.Now().UTC().Add(lease),
		})
	}
	ops, holder, err := s.repo.StartBatch(ops, req.ConflictResourceKeys, lease)
	if errors.Is(err, db.ErrLockConflict) {
		return &pb.StartOperationResponse{ConflictingOperation: toPb(holder)},
			status.Error(codes.AlreadyExists, "site resource is locked by an active operation")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "start site operation: %v", err)
	}
	resp := &pb.StartOperationResponse{Operation: toPb(ops[0]), ConflictsChecked: true}
	for _, op := range ops {
		resp.Operations = append(resp.Operations, toPb(op))
	}
	return resp, nil
}
