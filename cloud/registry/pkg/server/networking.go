package server

import (
	"context"
	"fmt"
	"github.com/ukama/ukamaX/cloud/registry/internal/db"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/common/ukama"
	"net"
)

type NetworkingServer struct {
	pb.UnimplementedNetworkingServer
	netRepo db.NetRepo
}

func NewNetworkingServer(netRepo db.NetRepo) *NetworkingServer {
	return &NetworkingServer{netRepo: netRepo}
}

func (n *NetworkingServer) ResolveNodeIp(c context.Context, req *pb.ResolveNodeIpRequest) (*pb.ResolveNodeIpResponse, error) {
	nd, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, err
	}

	ip, err := n.netRepo.GetIP(nd)
	if err != nil {
		return nil, err
	}

	return &pb.ResolveNodeIpResponse{
		Ip: ip.String(),
	}, nil
}

func (n *NetworkingServer) SetNodeIp(c context.Context, req *pb.SetNodeIpRequest) (*pb.SetNodeIpResponse, error) {
	nd, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(req.Ip)
	if ip == nil {
		return nil, fmt.Errorf("not valid ip")
	}

	err = n.netRepo.SetIp(nd, &ip)
	if err != nil {
		return nil, err
	}
	return &pb.SetNodeIpResponse{}, nil
}
