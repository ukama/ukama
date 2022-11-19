package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"github.com/ukama/ukama/systems/registry/network/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"

	"github.com/sirupsen/logrus"
)

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	netRepo    db.NetRepo
	orgRepo    db.OrgRepo
	siteRepo   db.SiteRepo
	orgService providers.OrgClientProvider
}

func NewNetworkServer(netRepo db.NetRepo, orgRepo db.OrgRepo, siteRepo db.SiteRepo,
	orgService providers.OrgClientProvider) *NetworkServer {
	return &NetworkServer{
		netRepo:    netRepo,
		orgRepo:    orgRepo,
		siteRepo:   siteRepo,
		orgService: orgService,
	}
}

func (n *NetworkServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Get the Org locally
	orgName := req.GetOrgName()
	networkName := req.GetName()
	org, err := n.orgRepo.GetByName(orgName)
	if err != nil {
		logrus.Infof("lookup for org %s remotely", orgName)

		svc, err := n.orgService.GetClient()
		if err != nil {
			return nil, err
		}

		remoteOrg, err := svc.GetByName(ctx, &orgpb.GetByNameRequest{Name: orgName})
		if err != nil {
			return nil, err
		}

		// What should we do if the remote org exists but is deactivated.
		if remoteOrg.Org.IsDeactivated {
			return nil, status.Errorf(codes.FailedPrecondition,
				"org is deactivated: cannot add network to it")
		}

		logrus.Infof("Adding remove org %s to local org repo", orgName)
		org = &db.Org{Name: remoteOrg.Org.Name,
			Deactivated: remoteOrg.Org.IsDeactivated}

		err = n.orgRepo.Add(org)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "org")
		}
	}

	logrus.Infof("Adding network %s", networkName)
	nt, err := n.netRepo.Add(org.ID, networkName)
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

func (n *NetworkServer) GetSite(ctx context.Context, req *pb.GetSiteRequest) (*pb.GetSiteResponse, error) {
	site, err := n.siteRepo.Get(uint(req.SiteID))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	return &pb.GetSiteResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (n *NetworkServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	ntwk, err := n.netRepo.Get(uint(req.GetNetID()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	sites, err := n.siteRepo.GetByNetwork(ntwk.ID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	resp := &pb.GetByNetworkResponse{
		NetworkID: uint64(ntwk.ID),
		Sites:     dbSitesToPbSites(sites),
	}

	return resp, nil
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

func dbSitesToPbSites(sites []db.Site) []*pb.Site {
	res := []*pb.Site{}

	for _, s := range sites {
		res = append(res, dbSiteToPbSite(&s))
	}

	return res
}
