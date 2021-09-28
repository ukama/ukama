package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"google.golang.org/grpc"
)

type GrpcClientError struct {
	HttpCode int
	Message  string
}

type Registry struct {
	conn    *grpc.ClientConn
	client  pb.RegistryServiceClient
	timeout int
	host    string
}

func NewRegistry(host string, timeout int) *Registry {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
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

func (r *Registry) GetOrg(orgName string) (string, *GrpcClientError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	res, err := r.client.GetOrg(ctx, &pb.GetOrgRequest{Name: orgName})

	return marshallResponse(err, res)
}

// GetOrg returns list of nodes
// org could be empty
func (r *Registry) GetNodes(owner string, orgName string) (string, *GrpcClientError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	res, err := r.client.GetNodes(ctx, &pb.GetNodesRequest{Owner: owner, OrgName: orgName})

	return marshallResponse(err, res)
}
