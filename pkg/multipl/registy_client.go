package multipl

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"google.golang.org/grpc"
	"time"
)

type registryClient struct {
	registryClient pb.RegistryServiceClient
	timeoutSecond  int
}

type RegistryClient interface {
	GetNodesList(orgName string) (nodes []*pb.Node, err error)
}

func NewRegistryClient(registryHost string, timeoutSecond int) (*registryClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, registryHost, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Errorf("Could not connect to registry: %v", err)
		return nil, err
	}

	return &registryClient{timeoutSecond: timeoutSecond,
		registryClient: pb.NewRegistryServiceClient(conn)}, nil
}

func (r registryClient) GetNodesList(orgName string) (nodes []*pb.Node, err error) {
	log.Info("Getting device list")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()
	resp, err := r.registryClient.GetNodes(ctx, &pb.GetNodesRequest{OrgName: orgName}, grpc_retry.WithMax(3))
	if err != nil {
		log.Errorf("Could not get device list: %v", err)
		return nil, err
	}

	if len(resp.Orgs) == 0 {
		log.Errorf("unexpected format of response")
		return nil, err
	}

	return resp.Orgs[0].Nodes, nil
}
