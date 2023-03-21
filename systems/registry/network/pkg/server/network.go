package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/network/pkg"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"github.com/ukama/ukama/systems/registry/network/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"

	"github.com/sirupsen/logrus"
)

const uuidParsingError = "Error parsing UUID"

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	netRepo        db.NetRepo
	orgRepo        db.OrgRepo
	siteRepo       db.SiteRepo
	orgService     providers.OrgClientProvider
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewNetworkServer(netRepo db.NetRepo, orgRepo db.OrgRepo, siteRepo db.SiteRepo,
	orgService providers.OrgClientProvider, msgBus mb.MsgBusServiceClient) *NetworkServer {
	return &NetworkServer{
		netRepo:        netRepo,
		orgRepo:        orgRepo,
		siteRepo:       siteRepo,
		orgService:     orgService,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
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

		// What should we do if the remote org exists but is deactivated?
		// For now we simply abort.
		if remoteOrg.Org.IsDeactivated {
			return nil, status.Errorf(codes.FailedPrecondition,
				"org is deactivated: cannot add network to it")
		}

		remoteOrgID, err := uuid.FromString(remoteOrg.Org.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid remote org id: %v", err)
		}

		logrus.Infof("Adding remove org %s to local org repo", orgName)
		org = &db.Org{
			Id:          remoteOrgID,
			Name:        remoteOrg.Org.Name,
			Deactivated: remoteOrg.Org.IsDeactivated}

		err = n.orgRepo.Add(org)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "org")
		}
	}

	network := &db.Network{
		Name:  networkName,
		OrgId: org.Id,
	}

	logrus.Infof("Adding network %s", networkName)
	err = n.netRepo.Add(network, func(*db.Network, *gorm.DB) error {
		network.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	route := n.baseRoutingKey.SetAction("add").SetObject("network").MustBuild()

	err = n.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.AddResponse{
		Network: dbNtwkToPbNtwk(network),
		Org:     req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	netID, err := uuid.FromString(req.NetworkId)
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
	orgID, err := uuid.FromString(req.OrgId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	ntwks, err := n.netRepo.GetByOrg(orgID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "networks")
	}

	resp := &pb.GetByOrgResponse{
		OrgId:    req.OrgId,
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
	route := n.baseRoutingKey.SetAction("delete").SetObject("network").MustBuild()
	err = n.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.DeleteResponse{}, nil
}

func (n *NetworkServer) AddSite(ctx context.Context, req *pb.AddSiteRequest) (*pb.AddSiteResponse, error) {
	netID, err := uuid.FromString(req.NetworkId)
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
		NetworkId: ntwk.Id,
		Name:      req.SiteName,
	}

	err = n.siteRepo.Add(site, func(*db.Site, *gorm.DB) error {
		site.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	route := n.baseRoutingKey.SetAction("add").SetObject("site").MustBuild()

	err = n.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.AddSiteResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (n *NetworkServer) GetSite(ctx context.Context, req *pb.GetSiteRequest) (*pb.GetSiteResponse, error) {
	siteID, err := uuid.FromString(req.SiteId)
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
	netID, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	ntwk, err := n.netRepo.Get(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	site, err := n.siteRepo.GetByName(ntwk.Id, req.SiteName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "site")
	}

	return &pb.GetSiteResponse{
		Site: dbSiteToPbSite(site)}, nil
}

func (n *NetworkServer) GetSitesByNetwork(ctx context.Context, req *pb.GetSitesByNetworkRequest) (*pb.GetSitesByNetworkResponse, error) {
	netID, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	ntwk, err := n.netRepo.Get(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	sites, err := n.siteRepo.GetByNetwork(ntwk.Id)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sites")
	}

	resp := &pb.GetSitesByNetworkResponse{
		NetworkId: ntwk.Id.String(),
		Sites:     dbSitesToPbSites(sites),
	}

	return resp, nil
}

func dbNtwkToPbNtwk(ntwk *db.Network) *pb.Network {
	return &pb.Network{
		Id:            ntwk.Id.String(),
		Name:          ntwk.Name,
		OrgId:         ntwk.OrgId.String(),
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
		Id:            site.Id.String(),
		Name:          site.Name,
		NetworkId:     site.NetworkId.String(),
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
