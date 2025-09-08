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
	"github.com/ukama/ukama/systems/common/rest/client"
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

const (
	DefaultMinuteDelayForPackageStartDate = 1
	eventPublishErrorMsg                  = "Failed to publish message %+v with key %+v. Errors %v"
)

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
		return nil, status.Error(codes.InvalidArgument,
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
		return nil, status.Error(codes.FailedPrecondition,
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
				"an unknown error occurred while getting iccid from sim token. Error %s", err.Error())
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
		PackageId:       packageId,
		IsActive:        true,
		DefaultDuration: packageInfo.Duration,
	}

	err = s.packageRepo.Add(firstPackage, func(pckg *sims.Package, tx *gorm.DB) error {
		firstPackage.Id = uuid.NewV4()
		firstPackage.SimId = sim.Id

		firstPackage.StartDate = time.Now().UTC().Add(time.Minute * DefaultMinuteDelayForPackageStartDate)
		firstPackage.EndDate = firstPackage.StartDate.Add(time.Hour * 24 * time.Duration(packageInfo.Duration))

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to add initial package to newly allocated sim. Error %s", err.Error())
	}

	sim.Package = *firstPackage

	simAgent, ok := s.agentFactory.GetAgentAdapter(simType)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim with ICCID: %q", simType, poolSim.Iccid)
	}

	agentRequest := client.AgentRequestData{
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.PackageId.String(),
		SimPackageId: sim.Package.Id.String(),
	}

	log.Infof("Activating sim on remote agent with request: %v", agentRequest)
	_, err = simAgent.BindSim(ctx, agentRequest)
	if err != nil {
		// TODO: think of rolling back the DB transaction on sim manager
		// if agent operation fails.

		return nil, status.Errorf(codes.Internal,
			"error while activating sim type %s on remote agent with request: %v", simType, agentRequest)
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
				emailTemplate.EmailKeyEndDate:    sim.Package.EndDate.Format("January 2, 2006"),
				emailTemplate.EmailKeyPackage:    packageInfo.Name,
				emailTemplate.EmailKeyDuration:   fmt.Sprintf("%v", packageInfo.Duration),
				emailTemplate.EmailKeyAmount:     fmt.Sprintf("%v", packageInfo.Amount),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	err = pushTotalSimsCountMetric(sim.NetworkId.String(), s.simRepo, s.orgId, s.pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim allocation operation: %s", err.Error())
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), s.simRepo, s.orgId, s.pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim allocation operation: %s", err.Error())
	}

	route := s.baseRoutingKey.SetAction("allocate").SetObject("sim").MustBuild()
	evt := &epb.EventSimAllocation{
		Id:             sim.Id.String(),
		SubscriberId:   sim.SubscriberId.String(),
		NetworkId:      sim.NetworkId.String(),
		OrgId:          orgId.String(),
		DataPlanId:     sim.Package.PackageId.String(),
		Iccid:          sim.Iccid,
		Msisdn:         sim.Msisdn,
		Imsi:           sim.Imsi,
		Type:           sim.Type.String(),
		Status:         sim.Status.String(),
		IsPhysical:     sim.IsPhysical,
		PackageId:      sim.Package.Id.String(),
		TrafficPolicy:  sim.TrafficPolicy,
		PackageEndDate: timestamppb.New(sim.Package.EndDate),
	}

	err = publishEventMessage(route, evt, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evt, route, err)
	}

	log.Infof("Allocating sim to subscriber success: %v", req.GetSubscriberId())

	return &pb.AllocateSimResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (s *SimManagerServer) GetSim(ctx context.Context, req *pb.GetSimRequest) (*pb.GetSimResponse, error) {
	log.Infof("Getting sim: %v", req.GetSimId())

	sim, err := getSim(req.SimId, s.simRepo)
	if err != nil {
		return nil, err
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, req.SimId)
	}

	log.Infof("Getting sim %s active record info from remote agent, if any...", sim.Iccid)
	_, err = simAgent.GetSim(ctx, sim.Iccid)
	if err != nil {
		log.Warnf("Failed to get active record info for sim %s. Error: %v", sim.Iccid, err)
		log.Warnf("Please make sure sim %s is properly configured and allocated", sim.Iccid)
	}

	return &pb.GetSimResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (s *SimManagerServer) GetUsages(ctx context.Context, req *pb.UsageRequest) (*pb.UsageResponse, error) {
	log.Infof("Getting Usages matching: %v", req)

	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument,
			"invalid value. Cdr type cannot be empty while getting usages")
	}

	var simType ukama.SimType
	var simIccid string

	if req.SimId != "" {
		sim, err := getSim(req.SimId, s.simRepo)
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
		return nil, status.Error(codes.Internal,
			"an unexpected error has occurred while unpacking usage response. Type is not map[string]any")
	}

	cost, ok := c.(map[string]any)
	if !ok {
		return nil, status.Error(codes.Internal,
			"an unexpected error has occurred while unpacking cost response. Type is not map[string]any")
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
	log.Infof("Getting sims matching: %v", req)

	if req.SubscriberId != "" {
		subscriberId, err := uuid.FromString(req.GetSubscriberId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for subscriber uuid: %s. Error %v", req.SubscriberId, err)
		}

		req.SubscriberId = subscriberId.String()
	}

	if req.NetworkId != "" {
		networkId, err := uuid.FromString(req.GetNetworkId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for network uuid: %s. Error %v", req.NetworkId, err)
		}

		req.NetworkId = networkId.String()
	}

	simType := ukama.SimTypeUnknown
	if req.SimType != "" {
		simType = ukama.ParseSimType(req.SimType)
		if simType == ukama.SimTypeUnknown {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid value for sim type: %s", req.SimType)
		}
	}

	simStatus := ukama.SimStatusUnknown
	if req.SimStatus != "" {
		simStatus = ukama.ParseSimStatus(req.SimStatus)
		if simStatus == ukama.SimStatusUnknown {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid value for sim status: %s", req.SimStatus)
		}
	}

	sims, err := s.simRepo.List(req.Iccid, req.Imsi, req.SubscriberId, req.NetworkId,
		simType, simStatus, req.TrafficPolicy, req.IsPhysical, req.Count, req.Sort)
	if err != nil {
		log.Errorf("Error while getting list of sims matching the given filters: %v",
			err)

		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	return &pb.ListSimsResponse{Sims: dbSimsToPbSims(sims)}, nil
}

// Deprecated: Use pkg.server.ListSims with subscriberId as filtering param instead.
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

// Deprecated: Use pkg.server.ListSims with networkId as filtering param instead.
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

func (s *SimManagerServer) TerminateSim(ctx context.Context, req *pb.TerminateSimRequest) (*pb.TerminateSimResponse, error) {
	log.Infof("Terminating sim: %v", req.GetSimId())

	sim, err := getSim(req.SimId, s.simRepo)
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
		pckg.TerminatedAt = time.Now().UTC()

		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	err = pushTerminatedSimsCountMetric(sim.NetworkId.String(), s.simRepo, s.orgId, s.pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim terminate operation: %s", err.Error())
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), s.simRepo, s.orgId, s.pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim terminate operation: %s", err.Error())
	}

	evtMsg := &epb.EventSimTermination{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
	}

	route := s.baseRoutingKey.SetAction("terminate").SetObject("sim").MustBuild()

	err = publishEventMessage(route, evtMsg, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	log.Infof("Sim %s terminated successfully", req.GetSimId())

	return &pb.TerminateSimResponse{}, nil
}

func (s *SimManagerServer) AddPackageForSim(ctx context.Context, req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	log.Infof("Adding package %v to sim: %v", req.GetPackageId(), req.GetSimId())

	formattedStart, err := validation.ValidateDate(req.GetStartDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sim, err := getSim(req.SimId, s.simRepo)
	if err != nil {
		return nil, status.Errorf(codes.NotFound,
			"invalid simId while adding package to sim. Error %s", err.Error())
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
		return nil, status.Error(codes.FailedPrecondition,
			"cannot set package to sim: data plan package is no more active within its org")
	}

	pkgInfoSimType := ukama.ParseSimType(pkgInfo.SimType)

	if sim.Type != pkgInfoSimType {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: sim (%s) and package (%s)'s sim types mismatch",
			sim.Type, pkgInfoSimType.String())
	}

	pkg := &sims.Package{
		SimId:           sim.Id,
		PackageId:       packageId,
		IsActive:        false,
		DefaultDuration: pkgInfo.Duration,
	}

	packages, err := s.packageRepo.List(req.SimId, "", "", "", "", "", false, false, 0, true)
	if err != nil {
		log.Errorf("failed to get the sorted list of packages present on sim (%s): %v",
			req.SimId, err)

		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	if len(packages) == 0 {
		startDate, err := time.Parse(time.RFC3339, formattedStart)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to parse start date: %v", err)
		}

		pkg.StartDate = startDate
		pkg.EndDate = pkg.StartDate.Add(time.Hour * 24 * time.Duration(pkgInfo.Duration))
		pkg.IsActive = true
	} else {
		pkg.StartDate = packages[len(packages)-1].EndDate.Add(time.Minute * DefaultMinuteDelayForPackageStartDate)
		pkg.EndDate = pkg.StartDate.Add(time.Hour * 24 * time.Duration(pkgInfo.Duration))
	}

	log.Infof("Package start date: %v, end date: %v", pkg.StartDate, pkg.EndDate)

	err = s.packageRepo.Add(pkg, func(pckg *sims.Package, tx *gorm.DB) error {
		pckg.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	route := s.baseRoutingKey.SetAction("addpackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimAddPackage{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    packageId.String(),
	}

	err = publishEventMessage(route, evtMsg, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	orgInfos, err := s.nucleusOrgClient.Get(s.orgName)
	if err != nil {
		return nil, err
	}

	userInfos, err := s.nucleusUserClient.GetById(orgInfos.Owner)
	if err != nil {
		return nil, err
	}

	subscriberRegistrySvc, err := s.subscriberRegistryService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteSubResp, err := subscriberRegistrySvc.Get(ctx, &subregpb.GetSubscriberRequest{
		SubscriberId: sim.SubscriberId.String(),
	})
	if err != nil {
		return nil, err
	}

	netInfo, err := s.networkClient.Get(sim.NetworkId.String())
	if err != nil {
		return nil, err
	}

	err = s.mailerClient.SendEmail(cnotif.SendEmailReq{
		To:           []string{remoteSubResp.Subscriber.Email},
		TemplateName: emailTemplate.EmailTemplatePackageAddition,
		Values: map[string]interface{}{
			emailTemplate.EmailKeySubscriber:      remoteSubResp.Subscriber.Name,
			emailTemplate.EmailKeyNetwork:         netInfo.Name,
			emailTemplate.EmailKeyName:            userInfos.Name,
			emailTemplate.EmailKeyOrg:             s.orgName,
			emailTemplate.EmailKeyPackagesCount:   fmt.Sprintf("%v", len(packages)+1),
			emailTemplate.EmailKeyPackagesDetails: fmt.Sprintf("$%.2f / %v %s / %d days", pkgInfo.Amount, pkgInfo.DataVolume, pkgInfo.DataUnit, pkgInfo.Duration),
			emailTemplate.EmailKeyExpiration:      pkg.EndDate.Format("January 2, 2006"),
			emailTemplate.EmailKeyPackage:         pkgInfo.Name,
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddPackageResponse{}, nil
}

func (s *SimManagerServer) ListPackagesForSim(ctx context.Context, req *pb.ListPackagesForSimRequest) (*pb.ListPackagesForSimResponse, error) {
	log.Infof("Getting packages  matching: %v", req)

	simId, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for sim uuid: %s. Error %v", req.SimId, err)
	}

	req.SimId = simId.String()

	if req.DataPlanId != "" {
		dataPlanId, err := uuid.FromString(req.GetDataPlanId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for data plan uuid: %s. Error %v", req.DataPlanId, err)
		}

		req.DataPlanId = dataPlanId.String()
	}

	if req.FromStartDate != "" {
		fromStartDate, err := validation.ValidateDate(req.GetFromStartDate())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		req.FromStartDate = fromStartDate
	}

	if req.ToStartDate != "" {
		toStartDate, err := validation.ValidateDate(req.GetToStartDate())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		req.ToStartDate = toStartDate
	}

	if req.FromEndDate != "" {
		fromEndDate, err := validation.ValidateDate(req.GetFromEndDate())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		req.FromEndDate = fromEndDate
	}

	if req.ToEndDate != "" {
		toEndDate, err := validation.ValidateDate(req.GetToEndDate())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		req.ToEndDate = toEndDate
	}

	packages, err := s.packageRepo.List(req.SimId, req.DataPlanId, req.FromStartDate, req.ToStartDate,
		req.FromEndDate, req.ToEndDate, req.IsActive, req.AsExpired, req.Count, req.Sort)
	if err != nil {
		log.Errorf("Error while getting list of packages present on sim (%s) matching the given filters: %v",
			req.SimId, err)

		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	return &pb.ListPackagesForSimResponse{Packages: dbPackagesToPbPackages(packages)}, nil
}

// Deprecated: Use pkg.server.ListPackagesForSim with simId as filtering param instead.
func (s *SimManagerServer) GetPackagesForSim(ctx context.Context, req *pb.GetPackagesForSimRequest) (*pb.GetPackagesForSimResponse, error) {
	log.Infof("Getting packages for sim: %v", req.GetSimId())

	simId, err := uuid.FromString(req.GetSimId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	packages, err := s.packageRepo.GetBySim(simId)
	if err != nil {
		log.Errorf("Failed to get the list of packages present on sim (%s): %v",
			req.SimId, err)

		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	resp := &pb.GetPackagesForSimResponse{
		SimId:    req.GetSimId(),
		Packages: dbPackagesToPbPackages(packages),
	}

	return resp, nil
}

func (s *SimManagerServer) SetActivePackageForSim(ctx context.Context, req *pb.SetActivePackageRequest) (*pb.SetActivePackageResponse, error) {
	log.Infof("Setting package %v as active for sim: %v", req.GetPackageId(), req.GetSimId())

	sim, err := getSim(req.SimId, s.simRepo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != ukama.SimStatusActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set active package on non active sim: sim's status is %s", sim.Status)
	}

	if sim.Package.Id != uuid.Nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim currently has package %v as active. This package needs to expire first",
			sim.Package.Id)
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
			"simID packageID mismatch: package %s does not belong to the provided sim %s",
			req.GetPackageId(), req.GetSimId())
	}

	if pkg.AsExpired {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set expired package (%s) as active", pkg.Id)
	}

	if pkg.IsActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot set already active package (%s) as active", pkg.Id)
	}

	// Update package on sim manager
	newPackageToActivate := &sims.Package{
		Id:       pkg.Id,
		IsActive: true,
	}

	err = s.packageRepo.Update(newPackageToActivate, func(pckg *sims.Package, tx *gorm.DB) error {
		// update startDate and endDate
		newPackageToActivate.StartDate = time.Now().UTC().
			Add(time.Minute * DefaultMinuteDelayForPackageStartDate)

		newPackageToActivate.EndDate = newPackageToActivate.StartDate.
			Add(time.Hour * 24 * time.Duration(pkg.DefaultDuration))

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to set package as active. Error %s", err.Error())
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, sim.Id)
	}

	agentRequest := client.AgentRequestData{
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.PackageId.String(),
		SimPackageId: sim.Package.Id.String(),
	}

	log.Infof("Updating package on remote agent for %s sim type with iccid %s",
		sim.Type.String(), sim.Iccid)

	err = simAgent.ActivateSim(ctx, agentRequest)
	if err != nil {
		log.Infof("Fail to update package on remote agent for %s sim type with iccid %s. Error: %v",
			sim.Type.String(), sim.Iccid, err)

		// TODO: think of rolling back the package update DB transaction on sim manager
		// if agent package update fails.

		return nil, fmt.Errorf("fail to update package on remote agent for %s sim type with iccid %s. Error: %w",
			sim.Type.String(), sim.Iccid, err)
	}

	// Publish the event only when both updates are successful
	route := s.baseRoutingKey.SetAction("activepackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimActivePackage{
		Id:               sim.Id.String(),
		SubscriberId:     sim.SubscriberId.String(),
		Iccid:            sim.Iccid,
		Imsi:             sim.Imsi,
		NetworkId:        sim.NetworkId.String(),
		PackageId:        pkg.Id.String(),
		PlanId:           pkg.PackageId.String(),
		PackageStartDate: timestamppb.New(newPackageToActivate.StartDate),
		PackageEndDate:   timestamppb.New(newPackageToActivate.EndDate),
	}

	err = publishEventMessage(route, evtMsg, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	return &pb.SetActivePackageResponse{}, nil
}

func (s *SimManagerServer) TerminatePackageForSim(ctx context.Context, req *pb.TerminatePackageRequest) (*pb.TerminatePackageResponse, error) {
	log.Infof("Terminating package %v for sim: %v", req.GetPackageId(), req.GetSimId())

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
			"simID packageID mismatch: package %s does not belong to the provided sim %s",
			req.GetPackageId(), req.GetSimId())
	}

	if !pckg.IsActive {
		log.Warnf("cannot terminate inactive package (%s). Skipping operation.", pckg.Id)

		return &pb.TerminatePackageResponse{}, nil
	}

	if pckg.AsExpired {
		return nil, status.Errorf(codes.FailedPrecondition,
			"package (%s) has already been marked as terminated", pckg.Id)
	}

	sim, err := getSim(req.SimId, s.simRepo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != ukama.SimStatusActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot terminate active package on non active sim: sim's status is is %s", sim.Status)
	}

	packageToTerminate := &sims.Package{
		Id:        pckg.Id,
		IsActive:  false,
		AsExpired: true,
	}

	err = s.packageRepo.Update(packageToTerminate, func(pckg *sims.Package, tx *gorm.DB) error {
		// update endDate
		packageToTerminate.EndDate = time.Now().UTC()

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to terminate package. Error %s", err.Error())
	}

	route := s.baseRoutingKey.SetAction("expirepackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimPackageExpire{
		Id:              sim.Id.String(),
		StartDate:       pckg.StartDate.String(),
		EndDate:         packageToTerminate.EndDate.String(),
		DefaultDuration: pckg.DefaultDuration,
		PackageId:       pckg.Id.String(),
		DataPlanId:      pckg.PackageId.String(),
	}

	err = publishEventMessage(route, evtMsg, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	return &pb.TerminatePackageResponse{}, nil
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
			"simID packageID mismatch: package %s does not belong to the provided sim %s",
			req.GetPackageId(), req.GetSimId())

	}

	if pckg.IsActive {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot remove active package (%s) from sim. Set package as not active first", pckg.Id)
	}

	sim, err := getSim(req.SimId, s.simRepo)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	err = s.packageRepo.Delete(packageId, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	route := s.baseRoutingKey.SetAction("removepackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimRemovePackage{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    packageId.String(),
	}

	err = publishEventMessage(route, evtMsg, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	return &pb.RemovePackageResponse{}, nil
}

func (s *SimManagerServer) activateSim(ctx context.Context, reqSimId string) (*pb.ToggleSimStatusResponse, error) {
	if err := activateSim(ctx, reqSimId, s.simRepo, s.agentFactory, s.orgId, s.pushMetricHost, s.msgbus, s.baseRoutingKey); err != nil {
		return nil, err
	}

	return &pb.ToggleSimStatusResponse{}, nil
}

func (s *SimManagerServer) deactivateSim(ctx context.Context, reqSimId string) (*pb.ToggleSimStatusResponse, error) {
	if err := deactivateSim(ctx, reqSimId, s.simRepo, s.agentFactory, s.orgId, s.pushMetricHost, s.msgbus, s.baseRoutingKey); err != nil {
		return nil, err
	}

	return &pb.ToggleSimStatusResponse{}, nil
}

func getSim(simId string, simRepo sims.SimRepo) (*sims.Sim, error) {
	parsedSimId, err := uuid.FromString(simId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := simRepo.Get(parsedSimId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return sim, nil
}

func activateSim(ctx context.Context, reqSimId string, simRepo sims.SimRepo, agentFactory adapters.AgentFactory,
	orgId, pushMetricHost string, msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	log.Infof("Activating sim: %v", reqSimId)

	sim, err := getSim(reqSimId, simRepo)
	if err != nil {
		return err
	}

	if sim.Status != ukama.SimStatusInactive {
		return status.Errorf(codes.FailedPrecondition,
			"sim state: %s is invalid for activation", sim.Status)
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

	err = simRepo.Update(simUpdates, nil)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "sim")
	}

	simAgent, ok := agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, reqSimId)
	}

	agentRequest := client.AgentRequestData{
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.PackageId.String(),
		SimPackageId: sim.Package.Id.String(),
	}

	err = simAgent.ActivateSim(ctx, agentRequest)
	if err != nil {
		// TODO: think of rolling back the DB transaction on sim manager
		// if agent operation fails.

		return err
	}

	err = pushActiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim activation operation: %s", err.Error())
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim activation operation: %s", err.Error())
	}

	route := baseRoutingKey.SetAction("activate").SetObject("sim").MustBuild()

	evtMsg := &epb.EventSimActivation{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.Id.String(),
	}

	err = publishEventMessage(route, evtMsg, msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	log.Infof("Sim %s activated successfully", reqSimId)

	return nil
}

func deactivateSim(ctx context.Context, reqSimId string, simRepo sims.SimRepo, agentFactory adapters.AgentFactory,
	orgId, pushMetricHost string, msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	log.Infof("Deactivating sim: %v", reqSimId)

	sim, err := getSim(reqSimId, simRepo)
	if err != nil {
		return err
	}

	if sim.Status != ukama.SimStatusActive {
		return status.Errorf(codes.FailedPrecondition,
			"sim state: %s is invalid for deactivation", sim.Status)
	}

	simAgent, ok := agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, reqSimId)
	}

	simUpdates := &sims.Sim{
		Id:                 sim.Id,
		Status:             ukama.SimStatusInactive,
		DeactivationsCount: sim.DeactivationsCount + 1}

	err = simRepo.Update(simUpdates, nil)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "sim")
	}

	agentRequest := client.AgentRequestData{
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.PackageId.String(),
		SimPackageId: sim.Package.Id.String(),
	}
	err = simAgent.DeactivateSim(ctx, agentRequest)
	if err != nil {
		// TODO: think of rolling back the DB transaction on sim manager
		// if agent operation fails.

		return err
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim deactivation operation: %s", err.Error())
	}

	err = pushActiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, pushMetricHost)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim deactivation operation: %s", err.Error())
	}

	route := baseRoutingKey.SetAction("deactivate").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimDeactivation{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    sim.Package.Id.String(),
	}

	err = publishEventMessage(route, evtMsg, msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	log.Infof("Sim %s deactivated successfully", reqSimId)

	return nil
}

func pushTotalSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId, pushMetricHost string) error {
	log.Infof("Collecting and pushing total sims count metric to push gateway host: %s", pushMetricHost)

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusUnknown, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting total sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting total sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = pmetric.CollectAndPushSimMetrics(pushMetricHost, pkg.SimMetric, pkg.NumberOfSims,
		float64(len(sims)), map[string]string{"network": networkId, "org": orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing total sims count metric to push gateway: %s", err.Error())

		return fmt.Errorf("error while pushing total sims count metric to push gateway: %w", err)
	}

	return nil
}

func pushActiveSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId, pushMetricHost string) error {
	log.Infof("Collecting and pushing active sims count metric to push gateway host: %s", pushMetricHost)

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusActive, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting active sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting active sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = pmetric.CollectAndPushSimMetrics(pushMetricHost, pkg.SimMetric, pkg.ActiveCount,
		float64(len(sims)), map[string]string{"network": networkId, "org": orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while active sims count metric to push gateway: %s", err.Error())

		return fmt.Errorf("error while pushing active sims count metric to push gateway: %w", err)
	}

	return nil
}

func pushInactiveSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId, pushMetricHost string) error {
	log.Infof("Collecting and pushing inactive sims count metric to push gateway host: %s", pushMetricHost)

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusInactive, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting inactive sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting inactive sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = pmetric.CollectAndPushSimMetrics(pushMetricHost, pkg.SimMetric, pkg.InactiveCount,
		float64(len(sims)), map[string]string{"network": networkId, "org": orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing inactive sims count metric to push gateway: %s", err.Error())

		return fmt.Errorf("error while pushing inactive sims count metric to push gateway: %w", err)
	}

	return nil
}

func pushTerminatedSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId, pushMetricHost string) error {
	log.Infof("Collecting and pushing terminated sims count metric to push gateway host: %s", pushMetricHost)

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusTerminated, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting terminated sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting terminated sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = pmetric.CollectAndPushSimMetrics(pushMetricHost, pkg.SimMetric, pkg.TerminatedCount,
		float64(len(sims)), map[string]string{"network": networkId, "org": orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing terminated sims count metric to push gateway: %s", err.Error())

		return fmt.Errorf("error while pushing terminated sims count metric to push gateway: %w", err)
	}

	return nil
}

func publishEventMessage(route string, msg protoreflect.ProtoMessage, msgbus mb.MsgBusServiceClient) error {
	err := msgbus.PublishRequest(route, msg)
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

	//TODO: remove usage of timestamp and update this.
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
		Id:              pkg.Id.String(),
		PackageId:       pkg.PackageId.String(),
		IsActive:        pkg.IsActive,
		DefaultDuration: pkg.DefaultDuration,
		AsExpired:       pkg.AsExpired,
		CreatedAt:       pkg.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       pkg.UpdatedAt.Format(time.RFC3339),
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
