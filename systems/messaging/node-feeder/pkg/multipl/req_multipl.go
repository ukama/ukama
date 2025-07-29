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
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
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

func (r *requestMultiplier) Process(req *cpb.NodeFeederMessage) error {
	// "org.nodeId"
	counter := 0
	//target = orgId.networkId.siteId.nodeId
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

	if nodeId != "*" {
		err := r.PublishToFilteredNodes(req, nodeResp.Nodes, orgName, networkName, siteName, nodeId)
		if err != nil {
			log.Errorf("Failed to publish message to queue: %s", err)
		}
	} else {
		err := r.PublishToNode(req, orgName, nodeId)
		if err != nil {
			log.Errorf("Failed to publish message to queue: %s", err)
			return fmt.Errorf("failed to publish message to queue")
		} else {
			log.Infof("Created %d requests for node id %s", counter, nodeId)
		}
	}

	log.Infof("Pulished requests %+v.", req)
	return nil
}

func (r *requestMultiplier) PublishToNode(req *cpb.NodeFeederMessage, orgName string, node string) error {
	err := r.queue.Publish(&cpb.NodeFeederMessage{
		Target:     orgName + "." + node,
		HTTPMethod: req.HTTPMethod,
		Path:       req.Path,
		Msg:        req.Msg,
	})

	if err != nil {
		log.Errorf("Failed to publish message to queue: %s", err)
		return fmt.Errorf("failed to publish message to queue")
	}
	return nil
}

func (r *requestMultiplier) PublishToFilteredNodes(req *cpb.NodeFeederMessage, nodeResp []*rc.NodeInfo, orgName string, networkId string, siteId string, nodeId string) error {
	counter := 0

	for _, n := range nodeResp {
		/* Figure a better way : This is generating a multiple nested request */
		nResp, err := r.nodeClient.Get(n.Id)
		if err != nil {
			return err
		}

		if networkId == "*" || nResp.Site.NetworkId == networkId {
			if siteId == "*" || nResp.Site.SiteId == siteId {
				if nodeId == "*" {
					err := r.PublishToNode(req, orgName, n.Id)
					if err != nil {
						log.Errorf("Failed to publish message to queue: %s", err)
						return fmt.Errorf("failed to publish message to queue")
					} else {
						counter++
					}
				}
			}
		}
	}

	log.Infof("Created %d requests", counter)
	return nil
}
