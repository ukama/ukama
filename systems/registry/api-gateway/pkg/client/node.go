package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
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
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, nodeHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
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

func (n *Node) AttachNodes(node, l, r string) (*pb.AttachNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.AttachNodes(ctx, &pb.AttachNodesRequest{
		NodeId:        node,
		AttachedNodes: []string{l, r},
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

func (n *Node) UpdateNode(nodeId string, name string) (*pb.UpdateNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.UpdateNode(ctx, &pb.UpdateNodeRequest{
		NodeId: nodeId,
		Name:   name,
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

func (n *Node) GetAllNodes() (*pb.GetAllNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetAllNodes(ctx, &pb.GetAllNodesRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) GetFreeNodes() (*pb.GetFreeNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetFreeNodes(ctx, &pb.GetFreeNodesRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) AddNode(nodeId, name, orgId, state string) (*pb.AddNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.AddNode(ctx, &pb.AddNodeRequest{
		NodeId: nodeId,
		Name:   name,
		OrgId:  orgId,
		State:  state,
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

func (n *Node) AddNodeToNetwork(nodeId, networkId, siteId string) (*pb.AddNodeToNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.AddNodeToNetwork(ctx, &pb.AddNodeToNetworkRequest{
		NodeId:    nodeId,
		NetworkId: networkId,
		SiteId:    siteId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) ReleaseNodeFromNetwork(nodeId string) (*pb.ReleaseNodeFromNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.ReleaseNodeFromNetwork(ctx, &pb.ReleaseNodeFromNetworkRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
