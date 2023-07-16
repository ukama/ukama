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

func (n *Node) GetOrgNodes(orgId string, free bool) (*pb.GetByOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetNodesForOrg(ctx, &pb.GetByOrgRequest{
		OrgId: orgId,
		Free:  free,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) GetAllNodes(free bool) (*pb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetNodes(ctx, &pb.GetNodesRequest{
		Free: free,
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
