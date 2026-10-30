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

	"github.com/ukama/ukama/systems/registry/node/pkg"

	log "github.com/sirupsen/logrus"
	metric "github.com/ukama/ukama/systems/common/metrics"
)

func (n *NodeServer) pushNodeMeterics(args ...string) {
	nodesCount, onlineCount, offlineCount, err := n.nodeRepo.GetNodeCount()
	if err != nil {
		log.Errorf("Error while getting node count %s", err.Error())

		return
	}

	log.Infof("Updating metrics for node NodeCount %d Online %d Offline %d", nodesCount,
		onlineCount, offlineCount)

	counts := map[string]int64{
		pkg.NumberOfNodes:        nodesCount,
		pkg.NumberOfOnlineNodes:  onlineCount,
		pkg.NumberOfOfflineNodes: offlineCount,
	}

	for _, arg := range args {
		count, ok := counts[arg]
		if !ok {
			continue
		}

		err = metric.CollectAndPushSystemMetrics(n.pushGateway, pkg.NodeMetric, arg,
			float64(count), nil, pkg.SystemName+"-"+pkg.ServiceName)
		if err != nil {
			log.Errorf("Error while pushing node metric %s to pushgateway %s", arg,
				err.Error())
		}
	}
}

func (n *NodeServer) PushMetrics() {
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfOnlineNodes, pkg.NumberOfOfflineNodes)
}

// requestMetricsPush asks the background pusher for fresh node counts without waiting on
// the pushgateway. A request already pending covers this one.
func (n *NodeServer) requestMetricsPush() {
	select {
	case n.metricsPush <- struct{}{}:
	default:
	}
}

// StartMetricsPusher pushes node counts after writes, off the request path, until ctx is
// done.
func (n *NodeServer) StartMetricsPusher(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-n.metricsPush:
				n.PushMetrics()
			}
		}
	}()
}
