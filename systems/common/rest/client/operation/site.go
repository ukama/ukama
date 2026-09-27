/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package operation

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"google.golang.org/grpc/metadata"
)

const (
	reservationIDKey = "ukama-operation-id"
	reservationTokenKey = "ukama-operation-token"
)

// NodeConflictResources includes every site node, even when only one node
// will receive a command. Unassigned nodes retain their individual lock.
func NodeConflictResources(nodes creg.NodeClient, nodeID string) ([]string, error) {
	if nodes == nil {
		return nil, fmt.Errorf("node registry is not configured")
	}
	node, err := nodes.Get(nodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}
	if node.Site.SiteId == "" {
		return []string{"node:" + nodeID}, nil
	}
	site, err := nodes.GetNodesBySite(node.Site.SiteId)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, fmt.Errorf("site %s nodes unavailable", node.Site.SiteId)
	}
	keys := SiteConflictResources(site.Nodes)
	for _, key := range keys {
		if key == "node:"+nodeID {
			return keys, nil
		}
	}
	return nil, fmt.Errorf("node %s missing from site %s", nodeID, node.Site.SiteId)
}

func SiteConflictResources(nodes []creg.NodeInfo) []string {
	seen := make(map[string]bool, len(nodes))
	keys := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if node.Id != "" && !seen[node.Id] {
			seen[node.Id] = true
			keys = append(keys, "node:"+node.Id)
		}
	}
	sort.Strings(keys)
	return keys
}

// WithReservation passes an already-reserved child operation to controller.
// Controller validates and claims it before registering or dispatching it.
func WithReservation(ctx context.Context, op *OperationInfo) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	md.Set(reservationIDKey, op.Id)
	md.Set(reservationTokenKey, strconv.FormatUint(op.FencingToken, 10))
	return metadata.NewOutgoingContext(ctx, md)
}

func ReservationFromContext(ctx context.Context) (string, uint64, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	ids, tokens := md.Get(reservationIDKey), md.Get(reservationTokenKey)
	if len(ids) == 0 && len(tokens) == 0 {
		return "", 0, nil
	}
	if len(ids) != 1 || ids[0] == "" || len(tokens) != 1 {
		return "", 0, fmt.Errorf("invalid operation reservation metadata")
	}
	token, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil || token == 0 {
		return "", 0, fmt.Errorf("invalid operation reservation token")
	}
	return ids[0], token, nil
}
