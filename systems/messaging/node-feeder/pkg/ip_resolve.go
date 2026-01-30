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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"

	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
)

type NodeIpResolver interface {
	Resolve(nodeId ukama.NodeID) (string, error)
}

type nodeIpResolver struct {
	nnsClient     pb.NnsClient
	timeoutSecond int
}

func NewNodeIpResolver(netHost string, timeoutSecond int) (*nodeIpResolver, error) {
	conn, err := grpc.NewClient(netHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Errorf("Could not connect to network service: %v", err)
		return nil, err
	}

	return &nodeIpResolver{timeoutSecond: timeoutSecond, nnsClient: pb.NewNnsClient(conn)}, nil
}

func (r *nodeIpResolver) Resolve(nodeId ukama.NodeID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()
	res, err := r.nnsClient.GetNode(ctx, &pb.GetNodeRequest{NodeId: nodeId.String()})
	if err != nil {
		return "", err
	}
	return res.NodeIp + ":" + strconv.Itoa(int(res.NodePort)), nil
}
