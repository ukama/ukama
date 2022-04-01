package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"

	"github.com/ukama/ukamaX/common/rest"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
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

func NewRegistryFromClient(registryClient pb.RegistryServiceClient) *Registry {
	return &Registry{
		host:    "localhost",
		timeout: 1,
		conn:    nil,
		client:  registryClient,
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

func (r *Registry) Add(orgName string, nodeId string) (*pb.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	resp, err := r.client.AddNode(ctx, &pb.AddNodeRequest{Node: &pb.Node{NodeId: nodeId}, OrgName: orgName})
	if err != nil {
		return nil, err
	}
	return resp.Node, nil
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
