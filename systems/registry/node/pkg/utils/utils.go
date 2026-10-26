/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package utils

import (
	"errors"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

// IsValidNodeConnectivity accepts the proto enum names in any case and rejects numbers
// and unknown strings.
func IsValidNodeConnectivity(value string) bool {
	for name := range cpb.NodeConnectivity_value {
		if strings.EqualFold(name, value) {
			return true
		}
	}

	return false
}

func InvalidNodeIDError(nodeId string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeId,
		err.Error())
}

func ProcessNodeDuplErrors(err error, nodeId string) error {
	var pge *pgconn.PgError

	if errors.As(err, &pge) && pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION {
		return status.Errorf(codes.AlreadyExists, "node with node id %s already exist", nodeId)
	}

	return grpc.StoreErrorToGrpc(err, "node")
}

func DbNodesToPbNodes(nodes []db.Node) []*pb.Node {
	pbNodes := []*pb.Node{}

	for _, n := range nodes {
		pbNodes = append(pbNodes, DbNodeToPbNode(&n))
	}

	return pbNodes
}

func DbNodeToPbNode(dbn *db.Node) *pb.Node {
	n := &pb.Node{
		Id: dbn.Id,
		Status: &pb.NodeStatus{
			Connectivity: cpb.NodeConnectivity(dbn.Status.Connectivity),
			State:        cpb.NodeState(dbn.Status.State),
		},
		Type:      dbn.Type.String(),
		Name:      dbn.Name,
		Latitude:  dbn.Latitude,
		Longitude: dbn.Longitude,
	}

	if dbn.Site.NodeId != "" {
		n.Site = &pb.Site{}
		n.Site.NodeId = dbn.Site.NodeId
		n.Site.SiteId = dbn.Site.SiteId.String()
		n.Site.NetworkId = dbn.Site.NetworkId.String()
		n.Site.AddedAt = timestamppb.New(dbn.Site.CreatedAt)
	}

	if len(dbn.Attached) > 0 {
		n.Attached = make([]*pb.Node, 0)
	}

	for _, nd := range dbn.Attached {
		n.Attached = append(n.Attached, DbNodeToPbNode(nd))
	}

	return n
}
