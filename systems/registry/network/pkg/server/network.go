package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	pmetric "github.com/ukama/ukama/systems/common/metrics"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/network/pkg"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"github.com/ukama/ukama/systems/registry/network/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"

	"github.com/sirupsen/logrus"
)

const uuidParsingError = "Error parsing UUID"

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	netRepo    db.NetRepo
	orgRepo    db.OrgRepo
	siteRepo   db.SiteRepo
	orgService providers.OrgClientProvider
	org string
	pushGatewayHost string
}

func NewNetworkServer(netRepo db.NetRepo, orgRepo db.OrgRepo, siteRepo db.SiteRepo,
	orgService providers.OrgClientProvider,org string,pushGatewayHost string) *NetworkServer {
	return &NetworkServer{
		netRepo:    netRepo,
		orgRepo:    orgRepo,
		siteRepo:   siteRepo,
		orgService: orgService,
		org : org,
		pushGatewayHost: pushGatewayHost,
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
networkId:= uuid.NewV4()
	network := &db.Network{
		ID:    networkId,
		Name:  networkName,
		OrgID: org.ID,
	}

	logrus.Infof("Adding network %s", networkName)
	err = n.netRepo.Add(network)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}
	networkCount, err := n.netRepo.GetNetworkCount()
	if err != nil {
		logrus.Errorf("failed to get network counts: %s", err.Error())
	}
	

	err = pmetric.CollectAndPushSimMetrics(n.pushGatewayHost, pkg.NetworkMetric, pkg.NumberOfNetwork, float64(networkCount), map[string]string{"network": networkId.String(), "org": n.org},pkg.SystemName)
	if err != nil {
		logrus.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}

	return &pb.AddResponse{
		Network: dbNtwkToPbNtwk(network),
		Org:     req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	netID, err := uuid.FromString(req.NetworkID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	nt, err := n.netRepo.Get(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.GetResponse{
		Network: dbNtwkToPbNtwk(nt),
	}, nil
}

func (n *NetworkServer) GetByName(ctx context.Context, req *pb.GetByNameRequest) (*pb.GetByNameResponse, error) {
	nt, err := n.netRepo.GetByName(req.GetOrgName(), req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "mapping org/network")
	}

	return &pb.GetByNameResponse{
		Network: dbNtwkToPbNtwk(nt),
		Org:     req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) GetByOrg(ctx context.Context, req *pb.GetByOrgRequest) (*pb.GetByOrgResponse, error) {
	orgID, err := uuid.FromString(req.OrgID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	ntwks, err := n.netRepo.GetByOrg(orgID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "networks")
	}

	resp := &pb.GetByOrgResponse{
		OrgID:    req.OrgID,
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
	networkCount, err := n.netRepo.GetNetworkCount()
	if err != nil {
		logrus.Errorf("failed to get network counts: %s", err.Error())
	}
	err = pmetric.CollectAndPushSimMetrics(n.pushGatewayHost, pkg.NetworkMetric, pkg.NumberOfNetwork, float64(networkCount), map[string]string{"network": "", "org": n.org},pkg.SystemName)
	if err != nil {
		logrus.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}

	return &pb.DeleteResponse{}, nil
}

func (n *NetworkServer) AddSite(ctx context.Context, req *pb.AddSiteRequest) (*pb.AddSiteResponse, error) {
	netID, err := uuid.FromString(req.NetworkID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	// We need to improve ukama/common/sql for more sql errors like foreign keys violations
	// which will allow us to skip these extra calls to DBs
	ntwk, err := n.netRepo.Get(netID)
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
	siteID, err := uuid.FromString(req.SiteID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	site, err := n.siteRepo.Get(siteID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	return &pb.GetSiteResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (n *NetworkServer) GetSiteByName(ctx context.Context, req *pb.GetSiteByNameRequest) (*pb.GetSiteResponse, error) {
	netID, err := uuid.FromString(req.NetworkID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	ntwk, err := n.netRepo.Get(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	site, err := n.siteRepo.GetByName(ntwk.ID, req.SiteName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	return &pb.GetSiteResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (n *NetworkServer) GetSitesByNetwork(ctx context.Context, req *pb.GetSitesByNetworkRequest) (*pb.GetSitesByNetworkResponse, error) {
	netID, err := uuid.FromString(req.NetworkID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	ntwk, err := n.netRepo.Get(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	sites, err := n.siteRepo.GetByNetwork(ntwk.ID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	resp := &pb.GetSitesByNetworkResponse{
		NetworkID: ntwk.ID.String(),
		Sites:     dbSitesToPbSites(sites),
	}

	return resp, nil
}

func dbNtwkToPbNtwk(ntwk *db.Network) *pb.Network {
	return &pb.Network{
		Id:            ntwk.ID.String(),
		Name:          ntwk.Name,
		OrgID:         ntwk.OrgID.String(),
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
		Id:            site.ID.String(),
		Name:          site.Name,
		NetworkID:     site.NetworkID.String(),
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
