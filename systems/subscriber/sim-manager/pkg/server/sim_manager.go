package server

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
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
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"

	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"

	pkgpb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	simpoolpb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	subregpb "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"
)

const DefaultDaysDelayForPackageStartDate = 1

type SimManagerServer struct {
	simRepo                   sims.SimRepo
	packageRepo               sims.PackageRepo
	agentFactory              adapters.AgentFactory
	packageService            providers.PackageClientProvider
	subscriberRegistryService providers.SubscriberRegistryClientProvider
	simPoolService            providers.SimPoolClientProvider
	key                       string
	msgbus                    mb.MsgBusServiceClient
	baseRoutingKey            msgbus.RoutingKeyBuilder
	pb.UnimplementedSimManagerServiceServer
}

func NewSimManagerServer(simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	agentFactory adapters.AgentFactory, packageService providers.PackageClientProvider,
	subscriberRegistryService providers.SubscriberRegistryClientProvider,
	simPoolService providers.SimPoolClientProvider, key string,
	msgBus mb.MsgBusServiceClient) *SimManagerServer {
	return &SimManagerServer{
		simRepo:                   simRepo,
		packageRepo:               packageRepo,
		agentFactory:              agentFactory,
		packageService:            packageService,
		subscriberRegistryService: subscriberRegistryService,
		simPoolService:            simPoolService,
		key:                       key,
		msgbus:                    msgBus,
		baseRoutingKey:            msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

func (s *SimManagerServer) AllocateSim(ctx context.Context, req *pb.AllocateSimRequest) (*pb.AllocateSimResponse, error) {
	subscriberID, err := uuid.Parse(req.GetSubscriberID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	subRegistrySvc, err := s.subscriberRegistryService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteSubResp, err := subRegistrySvc.Get(ctx,
		&subregpb.GetSubscriberRequest{SubscriberID: subscriberID.String()})
	if err != nil {
		return nil, err
	}

	if remoteSubResp.Subscriber.NetworkID != req.GetNetworkID() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid networkId: subscriber is not registered on the provided network")
	}

	packageID, err := uuid.Parse(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	packageSvc, err := s.packageService.GetClient()
	if err != nil {
		return nil, err
	}

	remotePkgResp, err := packageSvc.Get(ctx,
		&pkgpb.GetPackageRequest{PackageUuid: packageID.String()})
	if err != nil {
		return nil, err
	}

	if remotePkgResp.Package.OrgId != remoteSubResp.Subscriber.OrgID {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to subscriber sim: package does not belong to subscriber registerd org")
	}

	if !remotePkgResp.Package.Active {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to sim: package is no more active within its org")
	}

	strType := strings.ToLower(req.GetSimType())
	simType := sims.ParseType(strType)

	if uint8(simType) != uint8(remotePkgResp.Package.SimType) {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: provided sim type (%s) does not match with package allowed sim type (%s)",
			simType, remotePkgResp.Package.SimType)
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
			&simpoolpb.GetRequest{IsPhysicalSim: false, SimType: simpoolpb.SimType(simType)})

		if err != nil {
			return nil, err
		}

		poolSim = remoteSimPoolResp.Sim
	}

	networkID, err := uuid.Parse(remoteSubResp.Subscriber.NetworkID)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"invalid format of subscriber's network uuid. Error %s", err.Error())
	}

	orgID, err := uuid.Parse(remoteSubResp.Subscriber.OrgID)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"invalid format of subscriber's org uuid. Error %s", err.Error())
	}

	sim := &sims.Sim{
		SubscriberID: subscriberID,
		NetworkID:    networkID,
		OrgID:        orgID,
		Iccid:        poolSim.Iccid,
		Msisdn:       poolSim.Msisdn,
		Type:         simType,
		Status:       sims.SimStatusInactive,
		IsPhysical:   poolSim.IsPhysical,
	}

	err = s.simRepo.Add(sim, func(pckg *sims.Sim, tx *gorm.DB) error {
		txDb := sql.NewDbFromGorm(tx, pkg.IsDebugMode)

		sim.ID = uuid.New()
		startDate := time.Now().AddDate(0, 0, DefaultDaysDelayForPackageStartDate)
		endDate := startDate.Add(time.Duration(remotePkgResp.Package.Duration))

		firstPackage := &sims.Package{
			ID:        uuid.New(),
			SimID:     sim.ID,
			StartDate: startDate,
			EndDate:   endDate,
			PlanID:    packageID,
			IsActive:  false,
		}

		// Adding package to new allocated sim
		err := db.NewPackageRepo(txDb).Add(firstPackage, nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to allocate sim to subscriber. Error %s", err.Error())
	}

	resp := &pb.AllocateSimResponse{Sim: dbSimToPbSim(sim)}

	route := s.baseRoutingKey.SetAction("allocate").SetObject("sim").MustBuild()
	err = s.msgbus.PublishRequest(route, resp.Sim)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return resp, nil
}

func (s *SimManagerServer) GetSim(ctx context.Context, req *pb.GetSimRequest) (*pb.GetSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return &pb.GetSimResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (s *SimManagerServer) GetSimsBySubscriber(ctx context.Context, req *pb.GetSimsBySubscriberRequest) (*pb.GetSimsBySubscriberResponse, error) {
	subID, err := uuid.Parse(req.GetSubscriberID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetBySubscriber(subID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetSimsBySubscriberResponse{
		SubscriberID: req.GetSubscriberID(),
		Sims:         dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) GetSimsByNetwork(ctx context.Context, req *pb.GetSimsByNetworkRequest) (*pb.GetSimsByNetworkResponse, error) {
	netID, err := uuid.Parse(req.GetNetworkID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of network uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetByNetwork(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetSimsByNetworkResponse{
		NetworkID: req.GetNetworkID(),
		Sims:      dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) ToggleSimStatus(ctx context.Context, req *pb.ToggleSimStatusRequest) (*pb.ToggleSimStatusResponse, error) {
	strStatus := strings.ToLower(req.Status)
	simStatus := sims.ParseStatus(strStatus)

	switch simStatus {
	case sims.SimStatusActive:
		return s.activateSim(ctx, req.SimID)
	case sims.SimStatusInactive:
		return s.deactivateSim(ctx, req.SimID)
	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid status parameter: %s.", strStatus)
	}
}

func (s *SimManagerServer) DeleteSim(ctx context.Context, req *pb.DeleteSimRequest) (*pb.DeleteSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
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
			"invalid sim type: %q for sim ID: %q", sim.Type, req.SimID)
	}

	err = simAgent.TerminateSim(ctx, sim.Iccid)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		ID:     sim.ID,
		Status: sims.SimStatusTerminated,
	}

	err = s.simRepo.Update(simUpdates, func(pckg *sims.Sim, tx *gorm.DB) error {
		pckg.TerminatedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
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

	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	packageID, err := uuid.Parse(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	svc, err := s.packageService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteResp, err := svc.Get(ctx, &pkgpb.GetPackageRequest{PackageUuid: packageID.String()})
	if err != nil {
		return nil, err
	}

	if !remoteResp.Package.Active {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to sim: package is no more active within its org")
	}

	if sim.OrgID.String() != remoteResp.Package.OrgId {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid packageID: provided package does not belong to sim org issuer")
	}

	if uint8(sim.Type) != uint8(remoteResp.Package.SimType) {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: sim (%s) and packge (%s) sim types mismatch", sim.Type, remoteResp.Package.SimType.String())
	}

	pkg := &sims.Package{
		SimID:     sim.ID,
		StartDate: startDate,
		EndDate:   startDate.Add(time.Duration(remoteResp.Package.Duration)),
		PlanID:    packageID,
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
		pckg.ID = uuid.New()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	return &pb.AddPackageResponse{}, nil
}

func (s *SimManagerServer) GetPackagesBySim(ctx context.Context, req *pb.GetPackagesBySimRequest) (*pb.GetPackagesBySimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	packages, err := s.packageRepo.GetBySim(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	resp := &pb.GetPackagesBySimResponse{
		SimID:    req.GetSimID(),
		Packages: dbPackagesToPbPackages(packages),
	}

	return resp, nil
}

func (s *SimManagerServer) SetActivePackageForSim(ctx context.Context, req *pb.SetActivePackageRequest) (*pb.SetActivePackageResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
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

	packageID, err := uuid.Parse(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pkg, err := s.packageRepo.Get(packageID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	if pkg.SimID.String() != req.GetSimID() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid simID: packageID does not belong to the provided simID: %s", req.GetSimID())
	}

	if pkg.IsExpired() {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set expired package as active: package end date is %s", pkg.EndDate)
	}

	if pkg.IsActive {
		return &pb.SetActivePackageResponse{}, nil
	}

	newPackageToActivate := &sims.Package{
		ID:       pkg.ID,
		IsActive: true,
	}

	err = s.packageRepo.Update(newPackageToActivate, func(pckg *sims.Package, tx *gorm.DB) error {
		// if there is already an active package
		if sim.Package.ID != uuid.Nil {
			// get it
			currentActivePackage := &sims.Package{
				ID: sim.Package.ID,
			}

			// then deactivate it
			result := tx.Model(currentActivePackage).Update("is_active", false)
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}

			if result.Error != nil {
				return result.Error
			}
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
	packageID, err := uuid.Parse(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pckg, err := s.packageRepo.Get(packageID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	if pckg.SimID.String() != req.GetSimID() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid simID: packageID does not belong to the provided simID: %s", req.GetSimID())
	}

	err = s.packageRepo.Delete(packageID, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	return &pb.RemovePackageResponse{}, nil
}

func (s *SimManagerServer) activateSim(ctx context.Context, reqSimID string) (*pb.ToggleSimStatusResponse, error) {
	simID, err := uuid.Parse(reqSimID)
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
			"invalid sim type: %q for sim ID: %q", sim.Type, reqSimID)
	}

	err = simAgent.ActivateSim(ctx, sim.Iccid)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		ID:               sim.ID,
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

	return &pb.ToggleSimStatusResponse{}, nil
}

func (s *SimManagerServer) deactivateSim(ctx context.Context, reqSimID string) (*pb.ToggleSimStatusResponse, error) {
	simID, err := uuid.Parse(reqSimID)
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
			"invalid sim type: %q for sim ID: %q", sim.Type, reqSimID)
	}

	err = simAgent.DeactivateSim(ctx, sim.Iccid)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		ID:                 sim.ID,
		Status:             sims.SimStatusInactive,
		DeactivationsCount: sim.DeactivationsCount + 1}

	err = s.simRepo.Update(simUpdates, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return &pb.ToggleSimStatusResponse{}, nil
}

func dbSimToPbSim(sim *sims.Sim) *pb.Sim {
	res := &pb.Sim{
		Id:                 sim.ID.String(),
		SubscriberID:       sim.SubscriberID.String(),
		NetworkID:          sim.NetworkID.String(),
		OrgID:              sim.OrgID.String(),
		Iccid:              sim.Iccid,
		Msisdn:             sim.Msisdn,
		Imsi:               sim.Imsi,
		Type:               sim.Type.String(),
		Status:             sim.Status.String(),
		IsPhysical:         sim.IsPhysical,
		ActivationsCount:   sim.ActivationsCount,
		DeactivationsCount: sim.DeactivationsCount,
	}

	if sim.Package.ID != uuid.Nil {
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
		Id: pkg.ID.String(),
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
