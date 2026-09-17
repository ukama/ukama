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
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

type Node struct {
	conn    *grpc.ClientConn
	client  pb.NodeServiceClient
	timeout time.Duration
	host    string
}

func NewNode(nodeHost string, timeout time.Duration) *Node {
	conn, err := grpc.NewClient(nodeHost, grpc.WithTransportCredentials(insecure.NewCredentials()),
		ugrpc.TracingDialOption())
	if err != nil {
		log.Fatalf("Failed to connect to Node service: %v", err)
	}
	client := pb.NewNodeServiceClient(conn)

	return &Node{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    nodeHost,
	}
}

func NewNodeFromClient(nodeClient pb.NodeServiceClient) *Node {
	return &Node{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  nodeClient,
	}
}

func (n *Node) Close() {
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Node Service connection: %v", err)
		}
	}
}

func (n *Node) AddNode(ctx context.Context, nodeId, name, state string) (*pb.AddNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.AddNode(ctx, &pb.AddNodeRequest{
		NodeId: nodeId,
		Name:   name,
	})
}

func (n *Node) GetNode(ctx context.Context, nodeId string) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.GetNode(ctx, &pb.GetNodeRequest{
		NodeId: nodeId,
	})
}

func (n *Node) GetNetworkNodes(ctx context.Context, networkId string) (*pb.GetByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.GetNodesForNetwork(ctx, &pb.GetByNetworkRequest{
		NetworkId: networkId,
	})
}

func (n *Node) GetSiteNodes(ctx context.Context, siteId string) (*pb.GetBySiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.GetNodesForSite(ctx, &pb.GetBySiteRequest{
		SiteId: siteId,
	})
}

func (n *Node) GetNodes(ctx context.Context) (*pb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.GetNodes(ctx, &pb.GetNodesRequest{})
}

func (n *Node) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("State: %v, Connectivity: %v", req.State, req.Connectivity)

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.List(ctx, req)
}

func (n *Node) GetNodesByState(ctx context.Context, connectivity, state string) (*pb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.GetNodesByState(ctx, &pb.GetNodesByStateRequest{
		Connectivity: cpb.NodeConnectivity(ukama.ParseNodeConnectivity(connectivity)),
		State:        cpb.NodeState(ukama.ParseNodeState(state)),
	})
}

func (n *Node) UpdateNodeState(ctx context.Context, nodeId string, state string) (*pb.UpdateNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.UpdateNodeState(ctx, &pb.UpdateNodeStateRequest{
		NodeId: nodeId,
		State:  state,
	})
}

func (n *Node) UpdateNode(ctx context.Context, nodeId string, name string, latitude string, longitude string) (*pb.UpdateNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.UpdateNode(ctx, &pb.UpdateNodeRequest{
		Name:      name,
		NodeId:    nodeId,
		Latitude:  latitude,
		Longitude: longitude,
	})
}

func (n *Node) DeleteNode(ctx context.Context, nodeId string) (*pb.DeleteNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.DeleteNode(ctx, &pb.DeleteNodeRequest{
		NodeId: nodeId,
	})
}

func (n *Node) AttachNodes(ctx context.Context, node, l, r string) (*pb.AttachNodesResponse, error) {
	var attachedNodes []string

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	if l != "" {
		attachedNodes = append(attachedNodes, strings.ToLower(l))
	}

	if r != "" {
		attachedNodes = append(attachedNodes, strings.ToLower(r))
	}

	return n.client.AttachNodes(ctx, &pb.AttachNodesRequest{
		NodeId:        strings.ToLower(node),
		AttachedNodes: attachedNodes,
	})
}

func (n *Node) DetachNode(ctx context.Context, nodeId string) (*pb.DetachNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.DetachNode(ctx, &pb.DetachNodeRequest{
		NodeId: nodeId,
	})
}

func (n *Node) AddNodeToSite(ctx context.Context, nodeId, networkId, siteId string) (*pb.AddNodeToSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.AddNodeToSite(ctx, &pb.AddNodeToSiteRequest{
		NodeId:    nodeId,
		NetworkId: networkId,
		SiteId:    siteId,
	})
}

func (n *Node) ReleaseNodeFromSite(ctx context.Context, nodeId string) (*pb.ReleaseNodeFromSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.ReleaseNodeFromSite(ctx, &pb.ReleaseNodeFromSiteRequest{
		NodeId: nodeId,
	})
}
