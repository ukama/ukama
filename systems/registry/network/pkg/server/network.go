package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	"github.com/sirupsen/logrus"
)

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	orgRepo  db.OrgRepo
	netRepo  db.NetRepo
	siteRepo db.SiteRepo
}

func NewNetworkServer(orgRepo db.OrgRepo, netRepo db.NetRepo, siteRepo db.SiteRepo) *NetworkServer {
	return &NetworkServer{
		orgRepo:  orgRepo,
		netRepo:  netRepo,
		siteRepo: siteRepo,
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
	nt, err := n.netRepo.GetByName(req.OrgName, req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "mapping org/network")
	}

	return &pb.GetResponse{
		Network: dbNtwkToPbNtwk(nt),
		Org:     req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) GetByOrg(ctx context.Context, req *pb.GetByOrgRequest) (*pb.GetByOrgResponse, error) {
	org, err := n.orgRepo.GetByName(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	ntwks, err := n.netRepo.GetAllByOrgId(org.ID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "networks")
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

func (n *NetworkServer) AddSite(ctx context.Context, req *pb.AddSiteRequest) (*pb.AddSiteResponse, error) {
	// We need to improve ukama/common/sql for more sql errors like foreign keys violations
	// which will allow us to skip these extra calls to DBs
	ntwk, err := n.netRepo.Get(uint(req.NetID))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	site := &db.Site{
		NetworkID: ntwk.ID,
		Name:      req.SiteName,
	}

	err = n.siteRepo.Add(site)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	return &pb.AddSiteResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func dbNtwkToPbNtwk(ntwk *db.Network) *pb.Network {
	return &pb.Network{
		Id:            uint64(ntwk.ID),
		Name:          ntwk.Name,
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

func dbSiteToPbSite(site *db.Site) *pb.Site {
	return &pb.Site{
		Id:            uint64(site.ID),
		Name:          site.Name,
		NetworkID:     uint64(site.NetworkID),
		IsDeactivated: site.Deactivated,
		CreatedAt:     timestamppb.New(site.CreatedAt),
	}
}
