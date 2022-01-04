package pkg

import (
	"context"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc"
	"time"
)

type NodeIpResolver interface {
	Resolve(nodeId ukama.NodeID) (string, error)
}

type deviceIpResolver struct {
	netClient     pb.NetworkingClient
	timeoutSecond int
}

func NewDeviceIpResolver(registryHost string, timeoutSecond int) (*deviceIpResolver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, registryHost, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		logrus.Errorf("Could not connect to registry: %v", err)
		return nil, err
	}

	return &deviceIpResolver{timeoutSecond: timeoutSecond, netClient: pb.NewNetworkingClient(conn)}, nil
}

func (r *deviceIpResolver) Resolve(nodeId ukama.NodeID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()
	res, err := r.netClient.ResolveNodeIp(ctx, &pb.ResolveNodeIpRequest{NodeId: nodeId.String()})
	if err != nil {
		return "", err
	}
	return res.Ip, nil
}
