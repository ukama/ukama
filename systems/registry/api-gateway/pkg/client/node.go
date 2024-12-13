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
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/credentials/insecure"

	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"google.golang.org/grpc"
)

type Node struct {
	conn    *grpc.ClientConn
	client  pb.NodeServiceClient
	timeout time.Duration
	host    string
}

func NewNode(nodeHost string, timeout time.Duration) *Node {

	conn, err := grpc.NewClient(nodeHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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
	n.conn.Close()
}

func (n *Node) AddNode(nodeId, name, state string, latitude, longitude float64) (*pb.AddNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.AddNode(ctx, &pb.AddNodeRequest{
		NodeId:    nodeId,
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) GetNode(nodeId string) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetNode(ctx, &pb.GetNodeRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) GetNetworkNodes(networkId string) (*pb.GetByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetNodesForNetwork(ctx, &pb.GetByNetworkRequest{
		NetworkId: networkId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) GetSiteNodes(siteId string) (*pb.GetBySiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetNodesForSite(ctx, &pb.GetBySiteRequest{
		SiteId: siteId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) GetNodes() (*pb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetNodes(ctx, &pb.GetNodesRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) GetNodesByState(connectivity, state string) (*pb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetNodesByState(ctx, &pb.GetNodesByStateRequest{
		Connectivity: cpb.NodeConnectivity(ukama.ParseNodeConnectivity(connectivity)),
		State:        cpb.NodeState(ukama.ParseNodeState(state)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) UpdateNodeState(nodeId string, state string) (*pb.UpdateNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.UpdateNodeState(ctx, &pb.UpdateNodeStateRequest{
		NodeId: nodeId,
		State:  state,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) UpdateNode(nodeId string, name string, latitude, longitude float64) (*pb.UpdateNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.UpdateNode(ctx, &pb.UpdateNodeRequest{
		Name:      name,
		NodeId:    nodeId,
		Latitude:  latitude,
		Longitude: longitude,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) DeleteNode(nodeId string) (*pb.DeleteNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.DeleteNode(ctx, &pb.DeleteNodeRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) AttachNodes(node, l, r string) (*pb.AttachNodesResponse, error) {
	var attachedNodes []string

	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	if l != "" {
		attachedNodes = append(attachedNodes, strings.ToLower(l))
	}

	if r != "" {
		attachedNodes = append(attachedNodes, strings.ToLower(r))
	}

	res, err := n.client.AttachNodes(ctx, &pb.AttachNodesRequest{
		NodeId:        strings.ToLower(node),
		AttachedNodes: attachedNodes,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) DetachNode(nodeId string) (*pb.DetachNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.DetachNode(ctx, &pb.DetachNodeRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) AddNodeToSite(nodeId, networkId, siteId string) (*pb.AddNodeToSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.AddNodeToSite(ctx, &pb.AddNodeToSiteRequest{
		NodeId:    nodeId,
		NetworkId: networkId,
		SiteId:    siteId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) ReleaseNodeFromSite(nodeId string) (*pb.ReleaseNodeFromSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.ReleaseNodeFromSite(ctx, &pb.ReleaseNodeFromSiteRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
