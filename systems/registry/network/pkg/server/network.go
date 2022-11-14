package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	db2 "github.com/ukama/ukama/systems/registry/network/pkg/db"

	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	"github.com/sirupsen/logrus"
)

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	orgRepo db2.OrgRepo
	netRepo db2.NetRepo
}

func NewNetworkServer(orgRepo db2.OrgRepo, netRepo db2.NetRepo) *NetworkServer {
	return &NetworkServer{
		orgRepo: orgRepo,
		netRepo: netRepo,
	}
}

func (n *NetworkServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	org, err := n.orgRepo.MakeUserOrgExist(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	nt, err := n.netRepo.Add(org.ID, req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.AddResponse{
		Network: &pb.Network{
			Name: nt.Name,
		},
		Org: req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	nt, err := n.netRepo.Get(req.OrgName, req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org/network")
	}

	return &pb.GetResponse{
		Network: &pb.Network{
			Name: nt.Name,
		},
		Org: nt.Org.Name,
	}, nil
}

func (n *NetworkServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	logrus.Infof("Deleting network %s", req.Name)

	err := n.netRepo.Delete(req.OrgName, req.Name)
	if err != nil {
		logrus.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	resp := &pb.DeleteResponse{}

	// publish event before returning resp

	return resp, nil
}
