package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/services/common/rest"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	"google.golang.org/grpc"
)

type Registry struct {
	conn    *grpc.ClientConn
	client  pb.RegistryServiceClient
	timeout int
	host    string
}

func NewRegistry(host string, timeout int) *Registry {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewRegistryServiceClient(conn)

	return &Registry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewRegistryFromClient(networkClient pb.RegistryServiceClient) *Registry {
	return &Registry{
		host:    "localhost",
		timeout: 1,
		conn:    nil,
		client:  networkClient,
	}
}

func (r *Registry) Close() {
	r.conn.Close()
}

func (r *Registry) GetOrg(orgName string) (*pb.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	res, err := r.client.GetOrg(ctx, &pb.GetOrgRequest{Name: orgName})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Registry) AddOrUpdate(orgName string, nodeId string, name string) (node *pb.Node, isCreated bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	getResp, err := r.client.GetNode(ctx, &pb.GetNodeRequest{NodeId: nodeId})
	if err != nil && status.Code(err) == codes.NotFound {
		ar, err := r.client.AddNode(ctx, &pb.AddNodeRequest{Node: &pb.Node{NodeId: nodeId, Name: name}, OrgName: orgName})
		if err != nil {
			return nil, false, err
		}
		return ar.GetNode(), true, nil
	} else if err != nil {
		return nil, false, err
	}

	_, err = r.client.UpdateNode(ctx, &pb.UpdateNodeRequest{NodeId: nodeId, Name: name})
	if err != nil {
		return nil, false, err
	}
	getResp.Node.Name = name

	return getResp.Node, false, nil
}

// GetOrg returns list of nodes
func (r *Registry) GetNodes(orgName string) (*pb.GetNodesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
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

func (r *Registry) GetNode(nodeId string) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	res, err := r.client.GetNode(ctx, &pb.GetNodeRequest{
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	res, err := r.client.DeleteNode(ctx, &pb.DeleteNodeRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
