package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	db2 "github.com/ukama/ukama/systems/registry/network/pkg/db"
	"google.golang.org/protobuf/types/known/timestamppb"

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
		Network: dbNtwkToPbNtwk(nt),
		Org:     req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	nt, err := n.netRepo.Get(req.OrgName, req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org/network")
	}

	return &pb.GetResponse{
		Network: dbNtwkToPbNtwk(nt),
		Org:     req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) GetByOrg(ctx context.Context, req *pb.GetByOrgRequest) (*pb.GetByOrgResponse, error) {
	ntwks, err := n.netRepo.GetByOrg(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org/network")
	}

	resp := &pb.GetByOrgResponse{
		Org:      req.OrgName,
		Networks: dbNtwksToPbNtwks(ntwks),
	}

	return resp, nil
}

func (n *NetworkServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	logrus.Infof("Deleting network %s", req.Name)

	err := n.netRepo.Delete(req.OrgName, req.Name)
	if err != nil {
		logrus.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	// publish event before returning resp

	return &pb.DeleteResponse{}, nil
}

func dbNtwkToPbNtwk(ntwk *db.Network) *pb.Network {
	return &pb.Network{
		Id:   uint64(ntwk.ID),
		Name: ntwk.Name,
		// Org:           ntwk.Org.Name,
		IsDeactivated: ntwk.Deactivated,
		CreatedAt:     timestamppb.New(ntwk.CreatedAt),
	}
}

func dbNtwksToPbNtwks(ntwks []db.Network) []*pb.Network {
	res := []*pb.Network{}

	for _, n := range ntwks {
		res = append(res, dbNtwkToPbNtwk(&n))
	}

	return res
}
