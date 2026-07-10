/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package multipl

import (
	"fmt"
	"strings"

	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	rc "github.com/ukama/ukama/systems/common/rest/client/registry"
)

type requestMultiplier struct {
	nodeClient rc.NodeClient
	queue      QueuePublisher
}

func NewRequestMultiplier(registryClient string, queue QueuePublisher) pkg.RequestMultiplier {
	return &requestMultiplier{
		nodeClient: rc.NewNodeClient(registryClient),
		queue:      queue,
	}
}

// Process fans a wildcard-target request (org.network.site.node, where any of
// network/site/node may be "*") out into one concrete request per matching
// node, republished on the node-feeder retry queue. Republished targets are
// always fully-qualified 4-segment targets ending in a concrete node id, so
// they are executed directly on the next consumption and can never loop back
// into the multiplier.
func (r *requestMultiplier) Process(req *epb.NodeFeederMessage) error {
	//target = org.network.site.node
	segments := strings.Split(req.Target, ".")
	if len(segments) != 4 {
		return fmt.Errorf("invalid format of target: %s", req.Target)
	}

	orgName := segments[0]
	networkName := segments[1]
	siteName := segments[2]
	nodeId := segments[3]

	nodeResp, err := r.nodeClient.GetAll()
	if err != nil {
		return err
	}

	counter := 0

	for _, n := range nodeResp.Nodes {
		/* Figure a better way : This is generating a multiple nested request */
		nResp, err := r.nodeClient.Get(n.Id)
		if err != nil {
			return err
		}

		if networkName != "*" && nResp.Site.NetworkId != networkName {
			continue
		}

		if siteName != "*" && nResp.Site.SiteId != siteName {
			continue
		}

		if nodeId != "*" && !strings.EqualFold(n.Id, nodeId) {
			continue
		}

		err = r.PublishToNode(req, orgName, networkName, siteName, n.Id)
		if err != nil {
			log.Errorf("Failed to publish message to queue: %s", err)
			return fmt.Errorf("failed to publish message to queue")
		}

		counter++
	}

	log.Infof("Created %d node requests for target %s", counter, req.Target)
	return nil
}

// PublishToNode republishes the request for a single concrete node. The target
// keeps the original org/network/site segments (only the node segment is
// resolved), which is sufficient for the executor: it validates and uses only
// the node id segment.
func (r *requestMultiplier) PublishToNode(req *epb.NodeFeederMessage, orgName string, networkName string, siteName string, nodeId string) error {
	err := r.queue.Publish(&epb.NodeFeederMessage{
		Target:     orgName + "." + networkName + "." + siteName + "." + nodeId,
		HttpMethod: req.HttpMethod,
		Path:       req.Path,
		Msg:        req.Msg,
		NodeId:     nodeId,
	})

	if err != nil {
		log.Errorf("Failed to publish message to queue: %s", err)
		return fmt.Errorf("failed to publish message to queue")
	}
	return nil
}
