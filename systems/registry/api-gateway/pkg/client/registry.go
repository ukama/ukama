package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/systems/common/errors"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	pbnode "github.com/ukama/ukama/systems/registry/node/pb/gen"
	pborg "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"google.golang.org/grpc"
)

const DefaultNetworkName = "default"

type Registry struct {
	conn       *grpc.ClientConn
	orgConn    *grpc.ClientConn
	nodeConn   *grpc.ClientConn
	client     pb.NetworkServiceClient
	orgClient  pborg.OrgServiceClient
	nodeClient pbnode.NodeServiceClient
	timeout    time.Duration
	host       string
}

func NewRegistry(networkHost string, orgHost string, nodeHost string, timeout time.Duration) *Registry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewNetworkServiceClient(conn)

	orgConn, err := grpc.DialContext(ctx, orgHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	orgClient := pborg.NewOrgServiceClient(orgConn)

	nodeConn, err := grpc.DialContext(ctx, nodeHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	nodeClient := pbnode.NewNodeServiceClient(nodeConn)

	return &Registry{
		conn:       conn,
		client:     client,
		orgConn:    orgConn,
		orgClient:  orgClient,
		nodeConn:   nodeConn,
		nodeClient: nodeClient,
		timeout:    timeout,
		host:       networkHost,
	}
}

func NewRegistryFromClient(networkClient pb.NetworkServiceClient, orgClient pborg.OrgServiceClient, nodeClient pbnode.NodeServiceClient) *Registry {
	return &Registry{
		host:       "localhost",
		timeout:    1 * time.Second,
		conn:       nil,
		client:     networkClient,
		orgClient:  orgClient,
		nodeClient: nodeClient,
	}
}

func (r *Registry) Close() {
	r.conn.Close()
	r.orgConn.Close()
	r.nodeConn.Close()
}

func (r *Registry) GetOrg(orgName string) (*pborg.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.Get(ctx, &pborg.GetRequest{Name: orgName})
	if err != nil {
		return nil, err
	}
	return res.Org, nil
}

func (r *Registry) AddOrg(orgName string, owner string) (*pborg.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	organization := &pborg.Organization{Name: orgName, Owner: owner}
	res, err := r.orgClient.Add(ctx, &pborg.AddRequest{Org: organization})

	if err != nil {
		return nil, err
	}

	return res.Org, nil
}

func (r *Registry) Add(orgName string, nodeId string, name string, attachedNodes ...string) (node *pbnode.Node, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	addedNd, err := r.nodeClient.AddNode(ctx, &pbnode.AddNodeRequest{Node: &pbnode.Node{NodeId: nodeId, Name: name}})
	if err != nil {
		return nil, err
	}

	_, err = r.client.AddNode(ctx, &pb.AddNodeRequest{
		OrgName: orgName,
		Node: &pb.Node{
			NodeId: nodeId,
			Name:   addedNd.GetNode().Name,
		},
		Network: DefaultNetworkName,
	})

	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()
		logrus.Info("Deleting node from node service because it was not added to networks")
		_, derr := r.nodeClient.Delete(ctx, &pbnode.DeleteRequest{
			NodeId: nodeId,
		})
		if derr != nil {
			logrus.Errorf("Error deleting the node after failure. Error: %v", derr)
		}
		return nil, errors.Wrap(err, "failed to add node to network")
	}

	// consider making it async when we have UI notification infrastructure
	if len(attachedNodes) != 0 {
		if addedNd.Node.GetType() != pbnode.NodeType_TOWER {
			return nil, rest.HttpError{
				HttpCode: http.StatusBadRequest,
				Message:  fmt.Sprintf("Failed to attach nodes to node type %s", addedNd.Node.GetType())}
		}

		_, err = r.nodeClient.AttachNodes(ctx, &pbnode.AttachNodesRequest{
			ParentNodeId:    nodeId,
			AttachedNodeIds: attachedNodes,
		})

		if err != nil {
			return nil, errors.Wrap(err, "failed to attach node(s)")
		}
	}

	getResp, err := r.nodeClient.GetNode(ctx, &pbnode.GetNodeRequest{NodeId: nodeId})
	if err != nil {
		return nil, err
	}

	return getResp.Node, nil
}

func (r *Registry) UpdateNode(orgName string, nodeId string, name string, attachedNodes ...string) (node *pbnode.Node, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	nd, err := r.nodeClient.UpdateNode(ctx, &pbnode.UpdateNodeRequest{NodeId: nodeId, Name: name})
	if err != nil {
		return nil, err
	}

	if len(attachedNodes) > 0 {
		if nd.Node.GetType() != pbnode.NodeType_TOWER {
			return nil, rest.HttpError{
				HttpCode: http.StatusBadRequest,
				Message:  fmt.Sprintf("Cannot attach nodes to node type %s", nd.Node.GetType())}
		}

		_, err = r.nodeClient.AttachNodes(ctx, &pbnode.AttachNodesRequest{
			ParentNodeId:    nodeId,
			AttachedNodeIds: attachedNodes,
		})
		if err != nil {
			return nil, errors.Wrap(err, "Failed to attach nodes")
		}
	}

	getResp, err := r.nodeClient.GetNode(ctx, &pbnode.GetNodeRequest{NodeId: nodeId})
	if err != nil {
		return nil, err
	}

	return getResp.Node, nil
}

// GetOrg returns list of nodes
func (r *Registry) GetNodes(orgName string) (*pb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if len(orgName) == 0 {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest, Message: "Organization name is required"}
	}

	res, err := r.client.GetNodes(ctx, &pb.GetNodesRequest{OrgName: orgName})
	if err != nil {
		return nil, err
	}
	if res.Nodes == nil {
		// to keep 'nodes' as empty array in json response
		return &pb.GetNodesResponse{Nodes: []*pb.Node{}, OrgName: orgName}, nil
	}

	return res, nil
}

func (r *Registry) GetNode(nodeId string) (*pbnode.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.nodeClient.GetNode(ctx, &pbnode.GetNodeRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Registry) AttachNode(towerNodeId string, amplNodeId ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.nodeClient.AttachNodes(ctx, &pbnode.AttachNodesRequest{
		ParentNodeId:    towerNodeId,
		AttachedNodeIds: amplNodeId,
	})
	if err != nil {
		logrus.Errorf("Failed to attach node: %v", err)
	}
}

func (r *Registry) IsAuthorized(userId string, org string) (bool, error) {
	orgResp, err := r.GetOrg(org)
	if err != nil {
		if gErr, ok := err.(rest.HttpError); ok {
			if gErr.HttpCode != http.StatusNotFound {
				return false, nil
			}

			return false, gErr
		} else {
			return false, fmt.Errorf(err.Error())
		}
	}
	if orgResp.Owner == userId {
		return true, nil
	}
	return false, nil
}

func (r *Registry) DeleteNode(nodeId string) (*pb.DeleteNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.DeleteNode(ctx, &pb.DeleteNodeRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Registry) DetachNode(nodeId string, attachedId string) (*pbnode.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.nodeClient.DetachNode(ctx, &pbnode.DetachNodeRequest{
		DetachedNodeId: attachedId,
	})
	if err != nil {
		logrus.Errorf("Error detaching node %s. Error: %s", nodeId, err.Error())
		return nil, err
	}

	resp, err := r.nodeClient.GetNode(ctx, &pbnode.GetNodeRequest{
		NodeId: nodeId,
	})

	if err != nil {
		logrus.Warnf("Error getting node %s. Error %s", nodeId, err.Error())
		return nil, nil
	}
	return resp.Node, nil

}
