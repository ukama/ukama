/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"google.golang.org/grpc"
)

type Controller struct {
	conn    *grpc.ClientConn
	client  pb.ControllerServiceClient
	timeout time.Duration
	host    string
}

func NewController(controllerHost string, timeout time.Duration) *Controller {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, controllerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewControllerServiceClient(conn)

	return &Controller{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    controllerHost,
	}
}

func NewControllerFromClient(mClient pb.ControllerServiceClient) *Controller {
	return &Controller{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *Controller) Close() {
	r.conn.Close()
}

func (r *Controller) RestartSite(siteName, networkId string) (*pb.RestartSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.RestartSite(ctx, &pb.RestartSiteRequest{SiteName: siteName, NetworkId: networkId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Controller) RestartNode(nodeId string) (*pb.RestartNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.RestartNode(ctx, &pb.RestartNodeRequest{NodeId: nodeId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Controller) RestartNodes(networkId string, nodeIds []string) (*pb.RestartNodesResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
    defer cancel()

    res, err := r.client.RestartNodes(ctx, &pb.RestartNodesRequest{NetworkId: networkId, NodeIds: nodeIds})
    if err != nil {
        return nil, err
    }

    return res, nil
}
