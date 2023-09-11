package multipl

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type registryClient struct {
	registryClient nodepb.NodeServiceClient
	timeoutSecond  int
}

type RegistryClient interface {
	GetNodesList() (nodes []*nodepb.Node, err error)
}

func NewRegistryClient(registryHost string, timeoutSecond int) (*registryClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, registryHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Errorf("Could not connect to registry: %v", err)
		return nil, err
	}

	return &registryClient{timeoutSecond: timeoutSecond,
		registryClient: nodepb.NewNodeServiceClient(conn)}, nil
}

func (r registryClient) GetNodesList(orgName string) (nodes []*nodepb.Node, err error) {
	log.Info("Getting device list")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()
	resp, err := r.registryClient.GetNodes(ctx, &nodepb.GetNodesRequest{Free: true}, grpc_retry.WithMax(3))
	if err != nil {
		log.Errorf("Could not get device list: %v", err)
		return nil, err
	}

	return resp.Node, nil
}