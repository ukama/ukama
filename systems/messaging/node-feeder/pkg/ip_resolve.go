package pkg

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/cloud/net/pb/gen"
	"github.com/ukama/ukama/systems/common/ukama"
	"google.golang.org/grpc"
)

type NodeIpResolver interface {
	Resolve(nodeId ukama.NodeID) (string, error)
}

type deviceIpResolver struct {
	netClient     pb.NnsClient
	timeoutSecond int
}

func NewDeviceIpResolver(netHost string, timeoutSecond int) (*deviceIpResolver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, netHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		logrus.Errorf("Could not connect to network service: %v", err)
		return nil, err
	}

	return &deviceIpResolver{timeoutSecond: timeoutSecond, netClient: pb.NewNnsClient(conn)}, nil
}

func (r *deviceIpResolver) Resolve(nodeId ukama.NodeID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()
	res, err := r.netClient.Get(ctx, &pb.GetRequest{NodeId: nodeId.String()})
	if err != nil {
		return "", err
	}
	return res.Ip, nil
}
