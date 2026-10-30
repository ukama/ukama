/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

// assignNodesToSite assigns the nodes in one transaction and publishes node.assign for
// each new assignment.
func (n *NodeServer) assignNodesToSite(siteId, networkId uuid.UUID,
	nodeIds []string) error {
	assigned, err := n.siteRepo.AssignNodes(siteId, networkId, nodeIds)
	if err != nil {
		return grpc.StoreErrorToGrpc(err, "node")
	}

	if n.msgbus == nil {
		return nil
	}

	route := n.baseRoutingKey.SetAction("assign").SetObject("node").MustBuild()

	for _, site := range assigned {
		evt := &epb.EventRegistryNodeAssign{
			NodeId:  site.NodeId,
			Type:    ukama.NodeID(site.NodeId).GetNodeType(),
			Site:    site.SiteId.String(),
			Network: site.NetworkId.String(),
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route,
				err.Error())
		}
	}

	return nil
}

// publishNodeReleased publishes node.release for each released assignment.
func (n *NodeServer) publishNodeReleased(released ...db.Site) {
	if n.msgbus == nil {
		return
	}

	route := n.baseRoutingKey.SetAction("release").SetObject("node").MustBuild()

	for _, site := range released {
		evt := &epb.NodeReleasedEvent{
			NodeId:  site.NodeId,
			Type:    ukama.NodeID(site.NodeId).GetNodeType(),
			Site:    site.SiteId.String(),
			Network: site.NetworkId.String(),
		}

		err := n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route,
				err.Error())
		}
	}
}
