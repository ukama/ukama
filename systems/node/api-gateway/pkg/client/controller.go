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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
)

type Controller struct {
	conn    *grpc.ClientConn
	client  pb.ControllerServiceClient
	timeout time.Duration
	host    string
}

func NewController(controllerHost string, timeout time.Duration) *Controller {
	conn, err := grpc.NewClient(controllerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Controller service: %v", err)
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

func (c *Controller) Close() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close connection from Controller Service: %v", err)
		}
	}
}

func (c *Controller) RestartSite(siteId, networkId string) (*pb.RestartSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.RestartSite(ctx, &pb.RestartSiteRequest{SiteId: siteId, NetworkId: networkId})
}

func (c *Controller) PingNode(req *pb.PingNodeRequest) (*pb.PingNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.PingNode(ctx, req)
}

func (c *Controller) RestartNode(nodeId string) (*pb.RestartNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.RestartNode(ctx, &pb.RestartNodeRequest{NodeId: nodeId})
}

func (c *Controller) RestartNodes(networkId string, nodeIds []string) (*pb.RestartNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.RestartNodes(ctx, &pb.RestartNodesRequest{NetworkId: networkId, NodeIds: nodeIds})
}

func (c *Controller) ToggleInternetSwitch(status bool, port int32, siteId string) (*pb.ToggleInternetSwitchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.ToggleInternetSwitch(ctx, &pb.ToggleInternetSwitchRequest{Status: status, SiteId: siteId, Port: port})
}

func (c *Controller) ToggleRf(nodeId string, status bool) (*pb.ToggleRfSwitchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.ToggleRfSwitch(ctx, &pb.ToggleRfSwitchRequest{NodeId: nodeId, Status: status})
}
