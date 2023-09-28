package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/types"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/network/pkg"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"github.com/ukama/ukama/systems/registry/network/pkg/providers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	metric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

const uuidParsingError = "Error parsing UUID"

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	orgName        string
	netRepo        db.NetRepo
	orgRepo        db.OrgRepo
	siteRepo       db.SiteRepo
	orgService     providers.OrgClientProvider
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
}

func NewNetworkServer(orgName string, netRepo db.NetRepo, orgRepo db.OrgRepo, siteRepo db.SiteRepo,
	orgService providers.OrgClientProvider, msgBus mb.MsgBusServiceClient, pushGateway string) *NetworkServer {
	return &NetworkServer{
		orgName:        orgName,
		netRepo:        netRepo,
		orgRepo:        orgRepo,
		siteRepo:       siteRepo,
		orgService:     orgService,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
	}
}

func (n *NetworkServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Get the Org locally
	orgName := req.GetOrgName()
	networkName := req.GetName()

	org, err := n.orgRepo.GetByName(orgName)
	if err != nil {
		log.Infof("lookup for org %s remotely", orgName)

		// svc, err := n.orgService.GetClient()
		// if err != nil {
		// 	return nil, err
		// }

		remoteOrg, err := n.orgService.GetByName(orgName)
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

		log.Infof("Adding remove org %s to local org repo", orgName)
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
		Name:             networkName,
		OrgId:            org.Id,
		AllowedCountries: req.AllowedCountries,
		AllowedNetworks:  req.AllowedNetworks,
		Budget:           req.Budget,
		Overdraft:        req.Overdraft,
		TrafficPolicy:    req.TrafficPolicy,
		PaymentLinks:     req.PaymentLinks,
		SyncStatus:       types.SyncStatusPending,
	}

	log.Infof("Adding network %s", networkName)
	err = n.netRepo.Add(network, func(*db.Network, *gorm.DB) error {
		network.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	route := n.baseRoutingKey.SetAction("add").SetObject("network").MustBuild()

	evt := &epb.NetworkCreatedEvent{
		Id:               network.Id.String(),
		Name:             network.Name,
		OrgId:            network.OrgId.String(),
		AllowedCountries: network.AllowedCountries,
		AllowedNetworks:  network.AllowedNetworks,
		Budget:           network.Budget,
		Overdraft:        network.Overdraft,
		TrafficPolicy:    network.TrafficPolicy,
		PaymentLinks:     network.PaymentLinks,
		IsDeactivated:    network.Deactivated,
	}

	err = n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			evt, route, err.Error())
	}

	n.pushNetworkCount(org.Id)

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
	log.Infof("Deleting network %s", req.Name)

	org, err := n.orgRepo.GetByName(req.OrgName)
	if err != nil {
		log.Errorf("Failed to find org %s. Errors %s", req.OrgName, err.Error())
		return nil, err
	}

	err = n.netRepo.Delete(req.OrgName, req.Name)
	if err != nil {
		log.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	// publish event before returning resp
	route := n.baseRoutingKey.SetAction("delete").SetObject("network").MustBuild()
	err = n.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	n.pushNetworkCount(org.Id)

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
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	n.pushSiteCount(ntwk.OrgId, ntwk.Id)

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
		Id:               ntwk.Id.String(),
		Name:             ntwk.Name,
		OrgId:            ntwk.OrgId.String(),
		AllowedCountries: ntwk.AllowedCountries,
		AllowedNetworks:  ntwk.AllowedNetworks,
		Budget:           ntwk.Budget,
		Overdraft:        ntwk.Overdraft,
		TrafficPolicy:    ntwk.TrafficPolicy,
		PaymentLinks:     ntwk.PaymentLinks,
		IsDeactivated:    ntwk.Deactivated,
		SyncStatus:       ntwk.SyncStatus.String(),
		CreatedAt:        timestamppb.New(ntwk.CreatedAt),
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

func (n *NetworkServer) pushNetworkCount(orgId uuid.UUID) {
	networkCount, err := n.netRepo.GetNetworkCount(orgId)
	if err != nil {
		log.Errorf("failed to get network counts: %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NetworkMetric, pkg.NumberOfNetworks, float64(networkCount), map[string]string{"org": orgId.String()}, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing network count metric to pushgateway %s", err.Error())
	}
}

func (n *NetworkServer) pushSiteCount(orgId uuid.UUID, netId uuid.UUID) {
	siteCount, err := n.siteRepo.GetSiteCount(netId)
	if err != nil {
		log.Errorf("failed to get site count: %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NetworkMetric, pkg.NumberOfSites, float64(siteCount), map[string]string{"org": orgId.String(), "network": netId.String()}, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing network count metric to pushgateway %s", err.Error())
	}
}

func (n *NetworkServer) PushMetrics() error {

	// Push Network count metric per org to pushgateway
	orgs, err := n.orgRepo.GetAll()
	if err != nil {
		log.Errorf("Failed to get all networks. Error %s", err.Error())
		return err
	}

	for _, org := range orgs {
		n.pushNetworkCount(org.Id)
	}

	// Push Site count metric per network to pushgateway
	networks, err := n.netRepo.GetAll()
	if err != nil {
		log.Errorf("Failed to get all networks. Error %s", err.Error())
		return err
	}

	for _, network := range networks {
		n.pushSiteCount(network.OrgId, network.Id)
	}

	return nil
}
