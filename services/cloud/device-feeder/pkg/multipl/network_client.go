package multipl

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type networkClient struct {
	networkClient pb.NetworkServiceClient
	timeoutSecond int
}

type NetworkClient interface {
	GetNodesList(orgName string) (nodes []*pb.Node, err error)
}

func NewNetworkClient(networkHost string, timeoutSecond int) (*networkClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, networkHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Errorf("Could not connect to network: %v", err)
		return nil, err
	}

	return &networkClient{timeoutSecond: timeoutSecond,
		networkClient: pb.NewNetworkServiceClient(conn)}, nil
}

func (r networkClient) GetNodesList(orgName string) (nodes []*pb.Node, err error) {
	log.Info("Getting device list")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()
	resp, err := r.networkClient.GetNodes(ctx, &pb.GetNodesRequest{OrgName: orgName}, grpc_retry.WithMax(3))
	if err != nil {
		log.Errorf("Could not get device list: %v", err)
		return nil, err
	}

	return resp.Nodes, nil
}
