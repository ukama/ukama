/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/emailTemplate"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"

	log "github.com/sirupsen/logrus"
	pmetric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	cnotif "github.com/ukama/ukama/systems/common/rest/client/notification"
	cnuc "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	subregpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	simpoolpb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

//TODO; Replace all these GetBy with List functions.

const DefaultMinuteDelayForPackageStartDate = 1

type SimManagerServer struct {
	simRepo                   sims.SimRepo
	packageRepo               sims.PackageRepo
	agentFactory              adapters.AgentFactory
	packageClient             cdplan.PackageClient
	subscriberRegistryService providers.SubscriberRegistryClientProvider
	simPoolService            providers.SimPoolClientProvider
	key                       string
	msgbus                    mb.MsgBusServiceClient
	baseRoutingKey            msgbus.RoutingKeyBuilder
	orgId                     string
	orgName                   string
	pushMetricHost            string
	mailerClient              cnotif.MailerClient
	networkClient             creg.NetworkClient
	nucleusOrgClient          cnuc.OrgClient
	nucleusUserClient         cnuc.UserClient
	pb.UnimplementedSimManagerServiceServer
}

func NewSimManagerServer(
	orgName string, simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	agentFactory adapters.AgentFactory, packageClient cdplan.PackageClient,
	subscriberRegistryService providers.SubscriberRegistryClientProvider,
	simPoolService providers.SimPoolClientProvider, key string,
	msgBus mb.MsgBusServiceClient,
	orgId string,
	pushMetricHost string,
	mailerClient cnotif.MailerClient,
	networkClient creg.NetworkClient,
	nucleusOrgClient cnuc.OrgClient,
	nucleusUserClient cnuc.UserClient,

) *SimManagerServer {
	return &SimManagerServer{
		orgName:                   orgName,
		simRepo:                   simRepo,
		packageRepo:               packageRepo,
		agentFactory:              agentFactory,
		packageClient:             packageClient,
		subscriberRegistryService: subscriberRegistryService,
		simPoolService:            simPoolService,
		key:                       key,
		msgbus:                    msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).
			SetOrgName(orgName).SetService(pkg.ServiceName),
		orgId:             orgId,
		pushMetricHost:    pushMetricHost,
		mailerClient:      mailerClient,
		networkClient:     networkClient,
		nucleusOrgClient:  nucleusOrgClient,
		nucleusUserClient: nucleusUserClient,
	}
}

func (s *SimManagerServer) AllocateSim(ctx context.Context, req *pb.AllocateSimRequest) (*pb.AllocateSimResponse, error) {
	log.Infof("Allocating sim to subscriber: %v", req.GetSubscriberId())

	subscriberId, err := uuid.FromString(req.GetSubscriberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	subRegistrySvc, err := s.subscriberRegistryService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteSubResp, err := subRegistrySvc.Get(ctx,
		&subregpb.GetSubscriberRequest{SubscriberId: subscriberId.String()})
	if err != nil {
		return nil, err
	}

	if remoteSubResp.Subscriber.NetworkId != req.GetNetworkId() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid networkId: subscriber is not registered on the provided network")
	}

	packageId, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	packageInfo, err := s.packageClient.Get(packageId.String())
	// think about how to handle different types of rest errors
	if err != nil {
		return nil, err
	}

	if !packageInfo.IsActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to sim: package is no more active within its org")
	}

	strType := strings.ToLower(req.GetSimType())
	simType := ukama.ParseSimType(strType)
	pkgInfoSimType := ukama.ParseSimType(packageInfo.SimType)

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

		remoteSimPoolResp, err := simPoolSvc.GetByIccid(ctx,
			&simpoolpb.GetByIccidRequest{Iccid: iccid})
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

	networkId, err := uuid.FromString(remoteSubResp.Subscriber.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"invalid format of subscriber's network uuid. Error %s", err.Error())
	}

	netInfo, err := s.networkClient.Get(remoteSubResp.Subscriber.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network not found for that org %s", err.Error())
	}

	orgId, err := uuid.FromString(s.orgId)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"invalid format of subscriber's org uuid. Error %s", err.Error())
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(simType)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim with lCCID: %q", simType, poolSim.Iccid)
	}

	_, err = simAgent.BindSim(ctx, poolSim.Iccid)
	if err != nil {
		return nil, err
	}

	var trafficPolicy uint32

	// zero value traffic policy means pick the traffic policy of the upper layer
	if req.TrafficPolicy != 0 {
		trafficPolicy = req.TrafficPolicy
	} else if packageInfo.TrafficPolicy != 0 {
		trafficPolicy = packageInfo.TrafficPolicy
	} else {
		trafficPolicy = netInfo.TrafficPolicy
	}

	sim := &sims.Sim{
		SubscriberId:  subscriberId,
		NetworkId:     networkId,
		Iccid:         poolSim.Iccid,
		Msisdn:        poolSim.Msisdn,
		Type:          simType,
		Status:        ukama.SimStatusInactive,
		IsPhysical:    poolSim.IsPhysical,
		TrafficPolicy: trafficPolicy,
		SyncStatus:    ukama.StatusTypePending,
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
		PackageId: packageId,
		IsActive:  false,
	}

	err = s.packageRepo.Add(firstPackage, func(pckg *sims.Package, tx *gorm.DB) error {
		firstPackage.Id = uuid.NewV4()
		firstPackage.SimId = sim.Id

		firstPackage.StartDate = time.Now().Add(time.Minute * DefaultMinuteDelayForPackageStartDate)
		firstPackage.EndDate = firstPackage.StartDate.Add(time.Hour * 24 * time.Duration(packageInfo.Duration))

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to add initial package to newlly allocated sim. Error %s", err.Error())
	}

	sim.Package = *firstPackage
	resp := &pb.AllocateSimResponse{Sim: dbSimToPbSim(sim)}

	route := s.baseRoutingKey.SetAction("allocate").SetObject("sim").MustBuild()

	evt := &epb.EventSimAllocation{
		Id:            sim.Id.String(),
		SubscriberId:  sim.SubscriberId.String(),
		NetworkId:     sim.NetworkId.String(),
		OrgId:         orgId.String(),
		DataPlanId:    sim.Package.PackageId.String(),
		Iccid:         sim.Iccid,
		Msisdn:        sim.Msisdn,
		Imsi:          sim.Imsi,
		Type:          sim.Type.String(),
		Status:        sim.Status.String(),
		IsPhysical:    sim.IsPhysical,
		PackageId:     sim.Package.Id.String(),
		TrafficPolicy: sim.TrafficPolicy,
	}
	orgInfos, err := s.nucleusOrgClient.Get(s.orgName)
	if err != nil {
		return nil, err
	}

	userInfos, err := s.nucleusUserClient.GetById(orgInfos.Owner)
	if err != nil {
		return nil, err
	}

	if poolSim.QrCode != "" && !poolSim.IsPhysical {
		err = s.mailerClient.SendEmail(cnotif.SendEmailReq{
			To:           []string{remoteSubResp.Subscriber.Email},
			TemplateName: emailTemplate.EmailTemplateSimAllocation,
			Values: map[string]interface{}{
				emailTemplate.EmailKeySubscriber: remoteSubResp.Subscriber.Name,
				emailTemplate.EmailKeyNetwork:    netInfo.Name,
				emailTemplate.EmailKeyName:       userInfos.Name,
				emailTemplate.EmailKeyQRCode:     poolSim.QrCode,
				emailTemplate.EmailKeyVolume:     fmt.Sprintf("%v", packageInfo.DataVolume),
				emailTemplate.EmailKeyUnit:       packageInfo.DataUnit,
				emailTemplate.EmailKeyOrg:        s.orgName,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	_ = s.PublishEventMessage(route, evt)

	simsCount, _, _, _, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("failed to get Sims counts: %s", err.Error())
	}

	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.NumberOfSubscribers,
		float64(simsCount), map[string]string{"network": req.NetworkId, "org": s.orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}
	return resp, nil
}

func (s *SimManagerServer) GetSim(ctx context.Context, req *pb.GetSimRequest) (*pb.GetSimResponse, error) {
	log.Infof("Getting sim: %v", req.GetSimId())

	sim, err := s.getSim(req.SimId)
	if err != nil {
		return nil, err
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

func (s *SimManagerServer) GetUsages(ctx context.Context, req *pb.UsageRequest) (*pb.UsageResponse, error) {
	log.Infof("Getting Usages matching: %v", req)

	if req.Type == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid value. Cdr type cannot be empty while getting usages")
	}

	var simType ukama.SimType
	var simIccid string

	if req.SimId != "" {
		sim, err := s.getSim(req.SimId)
		if err != nil {
			return nil, err
		}

		simType = sim.Type
		simIccid = sim.Iccid
	} else {
		simType = ukama.ParseSimType(req.SimType)
		if simType == ukama.SimTypeUnknown {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid value for sim type: %s", req.SimType)
		}
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(simType)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"failure to get agentDeactivateSim for sim type: %q", simType)
	}

	u, c, err := simAgent.GetUsages(ctx, simIccid, req.Type, req.From, req.To, req.Region)
	if err != nil {
		return nil, err
	}

	usage, ok := u.(map[string]any)
	if !ok {
		return nil, status.Errorf(codes.Internal,
			"an unexpected error has occured while unpacking usage response. Type is not map[string]any")
	}

	cost, ok := c.(map[string]any)
	if !ok {
		return nil, status.Errorf(codes.Internal,
			"an unexpected error has occured while unpacking cost response. Type is not map[string]any")
	}

	usageProtoMsg, err := structpb.NewStruct(usage)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to marshall usages map response to proto message. Error %s", err)
	}

	costProtoMsg, err := structpb.NewStruct(cost)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to marshall cost map response to proto message. Error %s", err)
	}

	return &pb.UsageResponse{
		Usage: usageProtoMsg,
		Cost:  costProtoMsg,
	}, nil
}

func (s *SimManagerServer) ListSims(ctx context.Context, req *pb.ListSimsRequest) (*pb.ListSimsResponse, error) {
	return nil, nil
}

func (s *SimManagerServer) GetSimsBySubscriber(ctx context.Context, req *pb.GetSimsBySubscriberRequest) (*pb.GetSimsBySubscriberResponse, error) {
	log.Infof("Getting sims for subscriber: %v", req.SubscriberId)

	subId, err := uuid.FromString(req.GetSubscriberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetBySubscriber(subId)
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
	log.Infof("Getting sims for network: %v", req.NetworkId)

	netId, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of network uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetByNetwork(netId)
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
	log.Infof("Toggling status for sim: %v", req.GetSimId())

	strStatus := strings.ToLower(req.Status)
	simStatus := ukama.ParseSimStatus(strStatus)

	switch simStatus {
	case ukama.SimStatusActive:
		return s.activateSim(ctx, req.SimId)
	case ukama.SimStatusInactive:
		return s.deactivateSim(ctx, req.SimId)
	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid status parameter: %s.", strStatus)
	}
}

func (s *SimManagerServer) DeleteSim(ctx context.Context, req *pb.DeleteSimRequest) (*pb.DeleteSimResponse, error) {
	log.Infof("Deleting sim: %v", req.GetSimId())

	sim, err := s.getSim(req.SimId)
	if err != nil {
		return nil, err
	}

	if sim.Status != ukama.SimStatusInactive {
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
		Status: ukama.SimStatusTerminated,
	}

	err = s.simRepo.Update(simUpdates, func(pckg *sims.Sim, tx *gorm.DB) error {
		pckg.TerminatedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	evtMsg := &epb.EventSimTermination{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
	}
	route := s.baseRoutingKey.SetAction("terminate").SetObject("sim").MustBuild()
	_ = s.PublishEventMessage(route, evtMsg)

	_, _, _, terminatedCount, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("Failed to get terminated sim counts: %s", err.Error())
	}

	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.TerminatedCount,
		float64(terminatedCount), map[string]string{"org": s.orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing terminateSimCount metric to pushgateway %s", err.Error())
	}

	return &pb.DeleteSimResponse{}, nil
}

func (s *SimManagerServer) AddPackageForSim(ctx context.Context, req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	log.Infof("Adding package %v to sim: %v", req.GetPackageId(), req.GetSimId())

	formattedStart, err := validation.ValidateDate(req.GetStartDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if err := validation.IsFutureDate(formattedStart); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())

	}

	startDate, err := time.Parse(time.RFC3339, formattedStart)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse start date: %v", err)
	}

	sim, err := s.getSim(req.SimId)
	if err != nil {
		return nil, err
	}

	packageId, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pkgInfo, err := s.packageClient.Get(packageId.String())
	if err != nil {
		return nil, err
	}

	if !pkgInfo.IsActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set package to sim: package is no more active within its org")
	}

	pkgInfoSimType := ukama.ParseSimType(pkgInfo.SimType)

	if sim.Type != pkgInfoSimType {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: sim (%s) and packge (%s) sim types mismatch",
			sim.Type, pkgInfoSimType.String())
	}

	pkg := &sims.Package{
		SimId:     sim.Id,
		StartDate: startDate,
		EndDate:   startDate.Add(time.Hour * 24 * time.Duration(pkgInfo.Duration)),
		PackageId: packageId,
		IsActive:  false,
	}

	log.Infof("Package start date: %v, end date: %v", pkg.StartDate, pkg.EndDate)

	// overlappingPackages, err := s.packageRepo.GetOverlap(pkg)
	// if err != nil {
	// 	return nil, grpc.SqlErrorToGrpc(err, "packages")
	// }

	// if len(overlappingPackages) > 0 {
	// 	return nil, status.Errorf(codes.FailedPrecondition,
	// 		"cannot set package to sim: package validity period overlaps with %d or more other packaes set for this sim",
	// 		len(overlappingPackages))
	// }

	err = s.packageRepo.Add(pkg, func(pckg *sims.Package, tx *gorm.DB) error {
		pckg.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	evtMsg := &epb.EventSimAddPackage{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    packageId.String(),
	}
	route := s.baseRoutingKey.SetAction("addpackage").SetObject("sim").MustBuild()
	_ = s.PublishEventMessage(route, evtMsg)

	return &pb.AddPackageResponse{}, nil
}

func (s *SimManagerServer) GetPackagesBySim(ctx context.Context, req *pb.GetPackagesBySimRequest) (*pb.GetPackagesBySimResponse, error) {
	log.Infof("Getting packages for sim: %v", req.GetSimId())

	simId, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	packages, err := s.packageRepo.GetBySim(simId)
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
	log.Infof("Setting package %v as active for sim: %v", req.GetPackageId(), req.GetSimId())

	sim, err := s.getSim(req.SimId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != ukama.SimStatusActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set active package on non active sim: sim's status is is %s", sim.Status)
	}

	packageId, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pkg, err := s.packageRepo.Get(packageId)
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

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to set package as active. Error %s", err.Error())
	}

	evtMsg := &epb.EventSimActivation{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    pkg.Id.String(),
	}
	route := s.baseRoutingKey.SetAction("activepackage").SetObject("sim").MustBuild()
	_ = s.PublishEventMessage(route, evtMsg)

	/* Update package for opertaor */
	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, sim.Id)
	}

	opReq := adapters.ReqData{
		Iccid:     sim.Iccid,
		Imsi:      sim.Imsi,
		NetworkId: sim.NetworkId.String(),
		PackageId: sim.Package.Id.String(),
		SimId:     sim.Id.String(),
	}

	err = simAgent.UpdatePackage(ctx, opReq)
	if err != nil {
		return nil, err
	}

	return &pb.SetActivePackageResponse{}, nil
}

func (s *SimManagerServer) RemovePackageForSim(ctx context.Context, req *pb.RemovePackageRequest) (*pb.RemovePackageResponse, error) {
	log.Infof("Removing package %v for sim: %v", req.GetPackageId(), req.GetSimId())

	packageId, err := uuid.FromString(req.GetPackageId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pckg, err := s.packageRepo.Get(packageId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	if pckg.SimId.String() != req.GetSimId() {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid simID: packageID does not belong to the provided simID: %s", req.GetSimId())
	}

	sim, err := s.getSim(req.SimId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	err = s.packageRepo.Delete(packageId, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	evtMsg := &epb.EventSimRemovePackage{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    packageId.String(),
	}
	route := s.baseRoutingKey.SetAction("removepackage").SetObject("sim").MustBuild()
	_ = s.PublishEventMessage(route, evtMsg)

	return &pb.RemovePackageResponse{}, nil
}

func (s *SimManagerServer) activateSim(ctx context.Context, reqSimId string) (*pb.ToggleSimStatusResponse, error) {
	log.Infof("Activating sim: %v", reqSimId)

	sim, err := s.getSim(reqSimId)
	if err != nil {
		return nil, err
	}

	if sim.Status != ukama.SimStatusInactive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim state: %s is invalid for activation", sim.Status)
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, reqSimId)
	}

	req := adapters.ReqData{
		Iccid:     sim.Iccid,
		Imsi:      sim.Imsi,
		NetworkId: sim.NetworkId.String(),
		PackageId: sim.Package.Id.String(),
		SimId:     sim.Id.String(),
	}
	err = simAgent.ActivateSim(ctx, req)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		Id:               sim.Id,
		Status:           ukama.SimStatusActive,
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

	evtMsg := &epb.EventSimActivation{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.Id.String(),
	}
	route := s.baseRoutingKey.SetAction("activate").SetObject("sim").MustBuild()
	_ = s.PublishEventMessage(route, evtMsg)

	_, activeCount, _, _, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("Failed to get activated Sims counts: %s", err.Error())
	}

	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.ActiveCount,
		float64(activeCount), map[string]string{"org": s.orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing activateCount metric to pushgateway %s", err.Error())
	}

	return &pb.ToggleSimStatusResponse{}, nil
}

func (s *SimManagerServer) deactivateSim(ctx context.Context, reqSimId string) (*pb.ToggleSimStatusResponse, error) {
	log.Infof("Deactivating sim: %v", reqSimId)

	sim, err := s.getSim(reqSimId)
	if err != nil {
		return nil, err
	}

	if sim.Status != ukama.SimStatusActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim state: %s is invalid for deactivation", sim.Status)
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, reqSimId)
	}

	req := adapters.ReqData{
		Iccid:     sim.Iccid,
		Imsi:      sim.Imsi,
		NetworkId: sim.NetworkId.String(),
		SimId:     sim.Id.String(),
	}
	err = simAgent.DeactivateSim(ctx, req)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		Id:                 sim.Id,
		Status:             ukama.SimStatusInactive,
		DeactivationsCount: sim.DeactivationsCount + 1}

	err = s.simRepo.Update(simUpdates, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	evtMsg := &epb.EventSimDeactivation{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.Id.String(),
	}
	route := s.baseRoutingKey.SetAction("deactivate").SetObject("sim").MustBuild()
	_ = s.PublishEventMessage(route, evtMsg)

	_, _, inactiveCount, _, err := s.simRepo.GetSimMetrics()
	if err != nil {
		log.Errorf("failed to get inactive Sim counts: %s", err.Error())
	}

	err = pmetric.CollectAndPushSimMetrics(s.pushMetricHost, pkg.SimMetric, pkg.InactiveCount,
		float64(inactiveCount), map[string]string{"org": s.orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while push inactive metrics to pushgateway: %s", err.Error())
	}

	return &pb.ToggleSimStatusResponse{}, nil
}

func (s *SimManagerServer) getSim(simId string) (*sims.Sim, error) {
	parsedSimId, err := uuid.FromString(simId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(parsedSimId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return sim, nil
}

func (s *SimManagerServer) PublishEventMessage(route string, msg protoreflect.ProtoMessage) error {

	err := s.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
	}
	return err

}

func dbSimToPbSim(sim *sims.Sim) *pb.Sim {
	res := &pb.Sim{
		Id:                 sim.Id.String(),
		SubscriberId:       sim.SubscriberId.String(),
		NetworkId:          sim.NetworkId.String(),
		Iccid:              sim.Iccid,
		Msisdn:             sim.Msisdn,
		Imsi:               sim.Imsi,
		Type:               sim.Type.String(),
		Status:             sim.Status.String(),
		IsPhysical:         sim.IsPhysical,
		TrafficPolicy:      sim.TrafficPolicy,
		ActivationsCount:   sim.ActivationsCount,
		DeactivationsCount: sim.DeactivationsCount,
		SyncStatus:         sim.SyncStatus.String(),
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
		Id:        pkg.Id.String(),
		PackageId: pkg.PackageId.String(),
		IsActive:  pkg.IsActive,
		CreatedAt: pkg.CreatedAt.Format(time.RFC3339),
		UpdatedAt: pkg.UpdatedAt.Format(time.RFC3339),
	}

	if !pkg.EndDate.IsZero() {
		res.EndDate = pkg.EndDate.Format(time.RFC3339)
	}

	if !pkg.StartDate.IsZero() {
		res.StartDate = pkg.StartDate.Format(time.RFC3339)
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
