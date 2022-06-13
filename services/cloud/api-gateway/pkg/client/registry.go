package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/services/common/errors"
	"github.com/ukama/ukama/services/common/rest"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	pbnode "github.com/ukama/ukama/services/cloud/node/pb/gen"
	pborg "github.com/ukama/ukama/services/cloud/org/pb/gen"
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
	conn, err := grpc.Dial(networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewNetworkServiceClient(conn)

	orgConn, err := grpc.Dial(networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	orgClient := pborg.NewOrgServiceClient(orgConn)

	nodeConn, err := grpc.Dial(nodeHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (r *Registry) AddOrUpdate(orgName string, nodeId string, name string) (node *pbnode.Node, isCreated bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	getResp, err := r.nodeClient.GetNode(ctx, &pbnode.GetNodeRequest{NodeId: nodeId})
	if err != nil && status.Code(err) == codes.NotFound {
		ar, err := r.nodeClient.AddNode(ctx, &pbnode.AddNodeRequest{Node: &pbnode.Node{NodeId: nodeId, Name: name}})
		if err != nil {
			return nil, false, err
		}

		_, err = r.client.AddNode(ctx, &pb.AddNodeRequest{
			OrgName: orgName,
			Node: &pb.Node{
				NodeId: nodeId,
				Name:   name,
			},
			Network: DefaultNetworkName,
		})

		if err != nil {
			return nil, false, errors.Wrap(err, "failed to add node to network")
		}

		return ar.GetNode(), true, nil
	} else if err != nil {
		return nil, false, err
	}

	_, err = r.nodeClient.UpdateNode(ctx, &pbnode.UpdateNodeRequest{NodeId: nodeId, Name: name})
	if err != nil {
		return nil, false, err
	}
	getResp.Node.Name = name

	return getResp.Node, false, nil
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
