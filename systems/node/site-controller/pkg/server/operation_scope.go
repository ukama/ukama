/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	log "github.com/sirupsen/logrus"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/node/site-controller/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SiteControllerServer) reserveSiteOperations(nodes []creg.NodeInfo, targets []string, action string) (map[string]*copr.OperationInfo, error) {
	if len(targets) == 0 {
		return nil, status.Error(codes.NotFound, "no target nodes on site")
	}
	if s.opManager == nil {
		return nil, status.Error(codes.Unavailable, "operation manager is not configured")
	}
	conflicts := copr.SiteConflictResources(nodes)
	keys := make([]string, 0, len(targets))
	seen := make(map[string]bool, len(targets))
	for _, id := range targets {
		if id == "" || seen[id] {
			return nil, status.Error(codes.InvalidArgument, "missing or duplicate site operation target")
		}
		seen[id] = true
		key := "node:" + id
		found := false
		for _, candidate := range conflicts {
			if candidate == key {
				found = true
				break
			}
		}
		if !found {
			return nil, status.Error(codes.InvalidArgument, "operation target is not on the site")
		}
		keys = append(keys, key)
	}
	resp, err := s.opManager.Start(copr.StartRequest{
		Type: action, System: "node", ResourceKey: keys[0],
		RequestedBy: pkg.ServiceName, LeaseSeconds: s.opLeaseSecs,
		ConflictResourceKeys: conflicts, AdditionalResourceKeys: keys[1:],
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || !resp.ConflictsChecked || len(resp.Operations) != len(targets) {
		return nil, status.Error(codes.Internal, "operation manager did not reserve all site targets")
	}
	reserved := make(map[string]*copr.OperationInfo, len(targets))
	for i, id := range targets {
		op := resp.Operations[i]
		if op == nil || op.Id == "" || op.FencingToken == 0 || op.ResourceKey != keys[i] ||
			op.Type != action || op.Status != copr.StatusPending {
			return nil, status.Error(codes.Internal, "operation manager returned an invalid site reservation")
		}
		reserved[id] = op
	}
	return reserved, nil
}

func (s *SiteControllerServer) releaseUnsent(reservations map[string]*copr.OperationInfo) {
	for _, op := range reservations {
		if _, err := s.opManager.ForceUnlock(op.Id, pkg.ServiceName, "site command was not dispatched"); err != nil {
			log.Errorf("failed to release unused site reservation %s: %v", op.Id, err)
		}
	}
}
