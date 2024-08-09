package pkg

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"

	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"google.golang.org/grpc"
)

type NodeIpResolver interface {
	Resolve(nodeId ukama.NodeID) (string, error)
}

type nodeIpResolver struct {
	netClient     pb.NnsClient
	timeoutSecond int
}

func NewNodeIpResolver(netHost string, timeoutSecond int) (*nodeIpResolver, error) {
	conn, err := grpc.NewClient(netHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Errorf("Could not connect to network service: %v", err)
		return nil, err
	}

	return &nodeIpResolver{timeoutSecond: timeoutSecond, netClient: pb.NewNnsClient(conn)}, nil
}

func (r *nodeIpResolver) Resolve(nodeId ukama.NodeID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeoutSecond)*time.Second)
	defer cancel()
	res, err := r.netClient.Get(ctx, &pb.GetNodeIPRequest{NodeId: nodeId.String()})
	if err != nil {
		return "", err
	}
	return res.Ip, nil
}
