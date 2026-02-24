/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"context"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"

	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	providers "github.com/ukama/ukama/systems/messaging/node-feeder/pkg/provider"
)

type NodeIpResolver interface {
	Resolve(nodeId ukama.NodeID) (string, error)
}

type nodeIpResolver struct {
	nnsClient     providers.NnsClientProvider
	timeoutSecond int
}

func NewNodeIpResolver(netHost string, timeoutSecond int) (*nodeIpResolver, error) {
	return &nodeIpResolver{timeoutSecond: timeoutSecond, nnsClient: providers.NewNnsClientProvider(netHost)}, nil
}

func (r *nodeIpResolver) Resolve(nodeId ukama.NodeID) (string, error) {
	logrus.Infof("Resolving node %v", nodeId)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()

	svc, err := r.nnsClient.GetClient()
	if err != nil {
		logrus.Errorf("Error getting NNS client: %v", err)
		return "", err
	}

	res, err := svc.GetNode(ctx, &pb.GetNodeRequest{NodeId: nodeId.String()})
	if err != nil {
		logrus.Errorf("Error resolving node %v: %v", nodeId, err)
		return "", err
	}

	logrus.Infof("Resolved node %v to %v:%v", nodeId, res.NodeIp, res.NodePort)
	return res.NodeIp + ":" + strconv.Itoa(int(res.NodePort)), nil
}
