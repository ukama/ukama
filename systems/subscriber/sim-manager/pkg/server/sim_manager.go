package server

import (
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pmetric "github.com/ukama/ukama/systems/common/metrics"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"

	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"

	subregpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	simpoolpb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

const DefaultDaysDelayForPackageStartDate = 1

type SimManagerServer struct {
	simRepo                   sims.SimRepo
	packageRepo               sims.PackageRepo
	agentFactory              adapters.AgentFactory
	packageClient             providers.PackageClient
	subscriberRegistryService providers.SubscriberRegistryClientProvider
	simPoolService            providers.SimPoolClientProvider
	key                       string
	msgbus                    mb.MsgBusServiceClient
	baseRoutingKey            msgbus.RoutingKeyBuilder
	pb.UnimplementedSimManagerServiceServer
	org                string
	pushMetricHost     string
	notificationClient providers.NotificationClient
	networkClient      providers.NetworkInfoClient
}

func NewSimManagerServer(simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	agentFactory adapters.AgentFactory, packageClient providers.PackageClient,
	subscriberRegistryService providers.SubscriberRegistryClientProvider,
	simPoolService providers.SimPoolClientProvider, key string,
	msgBus mb.MsgBusServiceClient,
	org string,
	pushMetricHost string,
	notificationClient providers.NotificationClient,
	networkClient providers.NetworkInfoClient,
) *SimManagerServer {
	return &SimManagerServer{
		simRepo:                   simRepo,
		packageRepo:               packageRepo,
		agentFactory:              agentFactory,
		packageClient:             packageClient,
		subscriberRegistryService: subscriberRegistryService,
		simPoolService:            simPoolService,
		key:                       key,
		msgbus:                    msgBus,
		baseRoutingKey:            msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		org:                       org,
		pushMetricHost:            pushMetricHost,
		notificationClient:        notificationClient,
		networkClient:             networkClient,
	}
}

func (s *SimManagerServer) AllocateSim(ctx context.Context, req *pb.AllocateSimRequest) (*pb.AllocateSimResponse, error) {

	subscriberID, err := uuid.FromString(req.GetSubscriberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	subRegistrySvc, err := s.subscriberRegistryService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteSubResp, err := subRegistrySvc.Get(ctx,
		&subregpb.GetSubscriberRequest{SubscriberId: subscriberID.String()})
	if err != nil {
		return nil, err
	}

	if remoteSubResp.Subscriber.NetworkId != req.GetNetworkId() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid networkId: subscriber is not registered on the provided network")
	}

	packageID, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	packageInfo, err := s.packageClient.GetPackageInfo(packageID.String())
	// think about how to handle different types of rest errors
	if err != nil {
		return nil, err
	}

	if packageInfo.OrgId != remoteSubResp.Subscriber.OrgId {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to subscriber sim: package does not belong to subscriber registerd org")
	}

	if !packageInfo.IsActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to sim: package is no more active within its org")
	}

	strType := strings.ToLower(req.GetSimType())
	simType := sims.ParseType(strType)
	pkgInfoSimType := sims.ParseType(packageInfo.SimType)

	if simType != pkgInfoSimType {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: provided sim type (%s) does not match with package allowed sim type (%s)",
			simType, pkgInfoSimType)
	}

	poolSim := new(simpoolpb.Sim)

	simPoolSvc, err := s.simPoolService.GetClient()
	if err != nil {
		return nil, err
	}

	if req.SimToken != "" {
		iccid, err := utils.GetIccidFromToken(req.SimToken, s.key)
		if err != nil {
			return nil, status.Errorf(codes.Internal,
				"an unknown error occured while getting iccid from sim token. Error %s", err.Error())
		}

		remoteSimPoolResp, err := simPoolSvc.GetByIccid(ctx, &simpoolpb.GetByIccidRequest{Iccid: iccid})
		if err != nil {
			return nil, err
		}

		poolSim = remoteSimPoolResp.Sim

	} else {
		remoteSimPoolResp, err := simPoolSvc.Get(ctx,
			&simpoolpb.GetRequest{IsPhysicalSim: false, SimType: simType.String()})

		if err != nil {
			return nil, err
		}

		poolSim = remoteSimPoolResp.Sim
	}

	networkID, err := uuid.FromString(remoteSubResp.Subscriber.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"invalid format of subscriber's network uuid. Error %s", err.Error())
	}

	orgID, err := uuid.FromString(remoteSubResp.Subscriber.OrgId)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"invalid format of subscriber's org uuid. Error %s", err.Error())
	}

	sim := &sims.Sim{
		SubscriberId: subscriberID,
		NetworkId:    networkID,
		OrgId:        orgID,
		Iccid:        poolSim.Iccid,
		Msisdn:       poolSim.Msisdn,
		Type:         simType,
		Status:       sims.SimStatusInactive,
		IsPhysical:   poolSim.IsPhysical,
	}

	err = s.simRepo.Add(sim, func(pckg *sims.Sim, tx *gorm.DB) error {
		sim.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to allocate sim to subscriber. Error %s", err.Error())
	}

	firstPackage := &sims.Package{
		PackageId: packageID,
		IsActive:  false,
	}

	err = s.packageRepo.Add(firstPackage, func(pckg *sims.Package, tx *gorm.DB) error {
		firstPackage.Id = uuid.NewV4()
		firstPackage.SimId = sim.Id

		firstPackage.StartDate = time.Now().AddDate(0, 0, DefaultDaysDelayForPackageStartDate)
		firstPackage.EndDate = firstPackage.StartDate.Add(time.Duration(packageInfo.Duration))

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to add initial package to newlly allocated sim. Error %s", err.Error())
	}

	resp := &pb.AllocateSimResponse{Sim: dbSimToPbSim(sim)}

	route := s.baseRoutingKey.SetAction("allocate").SetObject("sim").MustBuild()

	err = s.msgbus.PublishRequest(route, resp.Sim)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	netInfo, err := s.networkClient.GetNetworkInfo(remoteSubResp.Subscriber.NetworkId, remoteSubResp.Subscriber.OrgId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network not found for that org %s", err.Error())
	}

	emailBody, err := utils.GenerateEmailBody(netInfo.Name, remoteSubResp.Subscriber.FirstName, poolSim.QrCode)
	if err != nil {
		return nil, err
	}
	if poolSim.QrCode != "" {
		err = s.notificationClient.SendEmail(providers.SendEmailReq{
			To:      []string{remoteSubResp.Subscriber.Email},
			Subject: "[Ukama] " + netInfo.Name + " invited you to use their network",
			Body:    emailBody,
			Values:  map[string]string{"SubscriberID": remoteSubResp.Subscriber.SubscriberId},
		})
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "Unable to send email %s", err.Error())
		}

	}

	simsCount, _, _, _, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("failed to get Sims counts: %s", err.Error())
	}

	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.NumberOfSubscribers, float64(simsCount), map[string]string{"network": req.NetworkId, "org": s.org}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}
	return resp, nil
}

func (s *SimManagerServer) GetSim(ctx context.Context, req *pb.GetSimRequest) (*pb.GetSimResponse, error) {
	simID, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, req.SimId)
	}

	_, err = simAgent.GetSim(ctx, sim.Iccid)
	if err != nil {
		return nil, err
	}

	return &pb.GetSimResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (s *SimManagerServer) ListSims(ctx context.Context, req *pb.ListSimsRequest) (*pb.ListSimsResponse, error) {
	return nil, nil
}

func (s *SimManagerServer) GetSimsBySubscriber(ctx context.Context, req *pb.GetSimsBySubscriberRequest) (*pb.GetSimsBySubscriberResponse, error) {
	subID, err := uuid.FromString(req.GetSubscriberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetBySubscriber(subID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetSimsBySubscriberResponse{
		SubscriberId: req.GetSubscriberId(),
		Sims:         dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) GetSimsByNetwork(ctx context.Context, req *pb.GetSimsByNetworkRequest) (*pb.GetSimsByNetworkResponse, error) {
	netID, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of network uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetByNetwork(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetSimsByNetworkResponse{
		NetworkId: req.GetNetworkId(),
		Sims:      dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) ToggleSimStatus(ctx context.Context, req *pb.ToggleSimStatusRequest) (*pb.ToggleSimStatusResponse, error) {
	strStatus := strings.ToLower(req.Status)
	simStatus := sims.ParseStatus(strStatus)

	switch simStatus {
	case sims.SimStatusActive:
		return s.activateSim(ctx, req.SimId)
	case sims.SimStatusInactive:
		return s.deactivateSim(ctx, req.SimId)
	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid status parameter: %s.", strStatus)
	}

}

func (s *SimManagerServer) DeleteSim(ctx context.Context, req *pb.DeleteSimRequest) (*pb.DeleteSimResponse, error) {
	simID, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != sims.SimStatusInactive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim state: %s is invalid for deletion", sim.Status)
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, req.SimId)
	}

	err = simAgent.TerminateSim(ctx, sim.Iccid)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		Id:     sim.Id,
		Status: sims.SimStatusTerminated,
	}

	err = s.simRepo.Update(simUpdates, func(pckg *sims.Sim, tx *gorm.DB) error {
		pckg.TerminatedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	route := s.baseRoutingKey.SetAction("delete").SetObject("sim").MustBuild()
	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	_, _, _, terminatedCount, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("Failed to get terminated sim counts: %s", err.Error())
	}

	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.TerminatedCount, float64(terminatedCount), map[string]string{"org": s.org}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing terminateSimCount metric to pushgateway %s", err.Error())
	}

	return &pb.DeleteSimResponse{}, nil
}

func (s *SimManagerServer) AddPackageForSim(ctx context.Context, req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	if err := req.GetStartDate().CheckValid(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid time format for package start_date. Error %s", err.Error())
	}

	startDate := req.GetStartDate().AsTime()

	if startDate.Before(time.Now()) {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package start date on the past: package start date is %s", startDate)
	}

	simID, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	packageID, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pkgInfo, err := s.packageClient.GetPackageInfo(packageID.String())
	if err != nil {
		return nil, err
	}

	if !pkgInfo.IsActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to sim: package is no more active within its org")
	}

	if sim.OrgId.String() != pkgInfo.OrgId {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid packageID: provided package does not belong to sim org issuer")
	}

	pkgInfoSimType := sims.ParseType(pkgInfo.SimType)

	if sim.Type != pkgInfoSimType {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: sim (%s) and packge (%s) sim types mismatch", sim.Type, pkgInfoSimType.String())
	}

	pkg := &sims.Package{
		SimId:     sim.Id,
		StartDate: startDate,
		EndDate:   startDate.Add(time.Duration(pkgInfo.Duration)),
		PackageId: packageID,
		IsActive:  false,
	}

	overlappingPackages, err := s.packageRepo.GetOverlap(pkg)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	if len(overlappingPackages) > 0 {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to sim: package validity period overlaps with %d or more other packaes set for this sim",
			len(overlappingPackages))
	}

	err = s.packageRepo.Add(pkg, func(pckg *sims.Package, tx *gorm.DB) error {
		pckg.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	route := s.baseRoutingKey.SetAction("addpackage").SetObject("sim").MustBuild()
	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.AddPackageResponse{}, nil
}

func (s *SimManagerServer) GetPackagesBySim(ctx context.Context, req *pb.GetPackagesBySimRequest) (*pb.GetPackagesBySimResponse, error) {
	simID, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	packages, err := s.packageRepo.GetBySim(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	resp := &pb.GetPackagesBySimResponse{
		SimId:    req.GetSimId(),
		Packages: dbPackagesToPbPackages(packages),
	}

	return resp, nil
}

func (s *SimManagerServer) SetActivePackageForSim(ctx context.Context, req *pb.SetActivePackageRequest) (*pb.SetActivePackageResponse, error) {
	simID, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != sims.SimStatusActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set active package on non active sim: sim's status is is %s", sim.Status)
	}

	packageID, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pkg, err := s.packageRepo.Get(packageID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	if pkg.SimId.String() != req.GetSimId() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid simID: packageID does not belong to the provided simID: %s", req.GetSimId())
	}

	if pkg.IsExpired() {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set expired package as active: package end date is %s", pkg.EndDate)
	}

	if pkg.IsActive {
		return &pb.SetActivePackageResponse{}, nil
	}

	newPackageToActivate := &sims.Package{
		Id:       pkg.Id,
		IsActive: true,
	}

	err = s.packageRepo.Update(newPackageToActivate, func(pckg *sims.Package, tx *gorm.DB) error {
		// if there is already an active package
		if sim.Package.Id != uuid.Nil {
			// get it
			currentActivePackage := &sims.Package{
				Id: sim.Package.Id,
			}

			// then deactivate it
			result := tx.Model(currentActivePackage).Update("active", false)
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}

			if result.Error != nil {
				return result.Error
			}
		}
		route := s.baseRoutingKey.SetAction("activepackage").SetObject("sim").MustBuild()
		err = s.msgbus.PublishRequest(route, req)
		if err != nil {
			logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to set package as active. Error %s", err.Error())
	}

	return &pb.SetActivePackageResponse{}, nil
}

func (s *SimManagerServer) RemovePackageForSim(ctx context.Context, req *pb.RemovePackageRequest) (*pb.RemovePackageResponse, error) {
	packageID, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pckg, err := s.packageRepo.Get(packageID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	if pckg.SimId.String() != req.GetSimId() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid simID: packageID does not belong to the provided simID: %s", req.GetSimId())
	}

	err = s.packageRepo.Delete(packageID, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}
	route := s.baseRoutingKey.SetAction("removepackage").SetObject("sim").MustBuild()
	err = s.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.RemovePackageResponse{}, nil
}

func (s *SimManagerServer) activateSim(ctx context.Context, reqSimID string) (*pb.ToggleSimStatusResponse, error) {

	simID, err := uuid.FromString(reqSimID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != sims.SimStatusInactive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim state: %s is invalid for activation", sim.Status)
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, reqSimID)
	}

	err = simAgent.ActivateSim(ctx, sim.Iccid)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		Id:               sim.Id,
		Status:           sims.SimStatusActive,
		ActivationsCount: sim.ActivationsCount + 1,
		LastActivatedOn:  time.Now(),
	}

	if sim.FirstActivatedOn.IsZero() {
		simUpdates.FirstActivatedOn = simUpdates.LastActivatedOn
	}

	err = s.simRepo.Update(simUpdates, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}
	msg := &pb.ToggleSimStatusRequest{
		SimId: reqSimID,
	}
	route := s.baseRoutingKey.SetAction("activate").SetObject("sim").MustBuild()
	err = s.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
	}

	_, activeCount, _, _, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("Failed to get activated Sims counts: %s", err.Error())
	}
	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.ActiveCount, float64(activeCount), map[string]string{"org": s.org}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing activateCount metric to pushgateway %s", err.Error())
	}

	return &pb.ToggleSimStatusResponse{}, nil
}

func (s *SimManagerServer) deactivateSim(ctx context.Context, reqSimID string) (*pb.ToggleSimStatusResponse, error) {

	simID, err := uuid.FromString(reqSimID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != sims.SimStatusActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim state: %s is invalid for deactivation", sim.Status)
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, reqSimID)
	}

	err = simAgent.DeactivateSim(ctx, sim.Iccid)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		Id:                 sim.Id,
		Status:             sims.SimStatusInactive,
		DeactivationsCount: sim.DeactivationsCount + 1}

	err = s.simRepo.Update(simUpdates, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}
	msg := &pb.ToggleSimStatusRequest{
		SimId: reqSimID,
	}
	route := s.baseRoutingKey.SetAction("deactivate").SetObject("sim").MustBuild()
	err = s.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
	}
	_, _, inactiveCount, _, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("failed to get inactive Sim counts: %s", err.Error())
	}
	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.InactiveCount, float64(inactiveCount), map[string]string{"org": s.org}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while push inactive metrics to pushgateway: %s", err.Error())
	}
	return &pb.ToggleSimStatusResponse{}, nil
}

func dbSimToPbSim(sim *sims.Sim) *pb.Sim {
	res := &pb.Sim{
		Id:                 sim.Id.String(),
		SubscriberId:       sim.SubscriberId.String(),
		NetworkId:          sim.NetworkId.String(),
		OrgId:              sim.OrgId.String(),
		Iccid:              sim.Iccid,
		Msisdn:             sim.Msisdn,
		Imsi:               sim.Imsi,
		Type:               sim.Type.String(),
		Status:             sim.Status.String(),
		IsPhysical:         sim.IsPhysical,
		ActivationsCount:   sim.ActivationsCount,
		DeactivationsCount: sim.DeactivationsCount,
	}

	if sim.Package.Id != uuid.Nil {
		res.Package = dbPackageToPbPackage(&sim.Package)
	}

	if !sim.FirstActivatedOn.IsZero() {
		res.FirstActivatedOn = timestamppb.New(sim.FirstActivatedOn)
	}

	if !sim.LastActivatedOn.IsZero() {
		res.LastActivatedOn = timestamppb.New(sim.LastActivatedOn)
	}

	if sim.AllocatedAt != 0 {
		res.AllocatedAt = timestamppb.New(sim.LastActivatedOn)
	}

	return res
}

func dbSimsToPbSims(sims []sims.Sim) []*pb.Sim {
	res := []*pb.Sim{}

	for _, s := range sims {
		res = append(res, dbSimToPbSim(&s))
	}

	return res
}

func dbPackageToPbPackage(pkg *sims.Package) *pb.Package {
	res := &pb.Package{
		Id: pkg.Id.String(),
	}

	if !pkg.EndDate.IsZero() {
		res.EndDate = timestamppb.New(pkg.EndDate)
	}

	if !pkg.StartDate.IsZero() {
		res.StartDate = timestamppb.New(pkg.StartDate)
	}

	return res
}

func dbPackagesToPbPackages(packages []sims.Package) []*pb.Package {
	res := []*pb.Package{}

	for _, s := range packages {
		res = append(res, dbPackageToPbPackage(&s))
	}

	return res
}

func isEsim(isEsim bool) (qrCode string, e error) {

	if isEsim {
		qrCode = "esim"

		return qrCode, nil
	} else {
		qrCode = "sim"

		return "", nil
	}

}
