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
	cnuc "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	cpay "github.com/ukama/ukama/systems/common/rest/client/payments"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	subregpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	simpoolpb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

//TODO; Replace all these GetBy with List functions.

const (
	DefaultMinuteDelayForPackageStartDate = 0
	eventPublishErrorMsg                  = "Failed to publish message %+v with key %+v. Errors %v"
	emailDateFormat                       = "January 2, 2006"

	// Payment metadata keys. metadataSimKey ties a payment to a SIM;
	// metadataProvisionedKey marks a payment whose package was already added by
	// the caller (sim allocation), so the payment.success handler must not add it
	// again.
	metadataSimKey         = "sim"
	metadataProvisionedKey = "provisioned"
	metadataProvisionedYes = "true"
)

type SimManagerServer struct {
	simRepo                   sims.SimRepo
	packageRepo               sims.PackageRepo
	agentFactory              adapters.AgentFactory
	packageClient             cdplan.PackageClient
	subscriberRegistryService providers.SubscriberRegistryClientProvider
	simPoolService            providers.SimPoolClientProvider
	tokenCodec                utils.Codec
	msgbus                    mb.MsgBusServiceClient
	baseRoutingKey            msgbus.RoutingKeyBuilder
	orgId                     string
	orgName                   string
	metricsPusher             MetricsPusher
	networkClient             creg.NetworkClient
	nucleusOrgClient          cnuc.OrgClient
	nucleusUserClient         cnuc.UserClient
	paymentClient             cpay.PaymentClient
	pb.UnimplementedSimManagerServiceServer
}

func NewSimManagerServer(
	orgName string, simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	agentFactory adapters.AgentFactory, packageClient cdplan.PackageClient,
	subscriberRegistryService providers.SubscriberRegistryClientProvider,
	simPoolService providers.SimPoolClientProvider,
	tokenCodec utils.Codec,
	msgBus mb.MsgBusServiceClient,
	orgId string,
	pushMetricHost string,
	networkClient creg.NetworkClient,
	nucleusOrgClient cnuc.OrgClient,
	nucleusUserClient cnuc.UserClient,
	paymentClient cpay.PaymentClient,
) *SimManagerServer {
	s := &SimManagerServer{
		orgName:                   orgName,
		simRepo:                   simRepo,
		packageRepo:               packageRepo,
		agentFactory:              agentFactory,
		packageClient:             packageClient,
		subscriberRegistryService: subscriberRegistryService,
		simPoolService:            simPoolService,
		tokenCodec:                tokenCodec,
		msgbus:                    msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).
			SetOrgName(orgName).SetService(pkg.ServiceName),
		orgId:             orgId,
		metricsPusher:     NewMetricsPusher(pushMetricHost),
		networkClient:     networkClient,
		nucleusOrgClient:  nucleusOrgClient,
		nucleusUserClient: nucleusUserClient,
		paymentClient:     paymentClient,
	}

	return s
}

func (s *SimManagerServer) AllocateSim(ctx context.Context, req *pb.AllocateSimRequest) (*pb.SimResponse, error) {
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

	// We validate package duration to make sure it is greater than 0 min and lesser than 1000 years
	if err := validation.ValidatePackageDuration(packageInfo.Duration); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %v",
			fmt.Errorf("invalid package duration: %w", err))
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
		iccid, err := s.tokenCodec.GetIccidFromToken(req.SimToken)
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
		Status:        ukama.SimStatusServiceOff,
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

	// totalData is in bytes and packageInfo.DataVolume depends on pack.DataUnit
	dataUnit := ukama.ParseDataUnitType(packageInfo.DataUnit)
	if dataUnit == ukama.DataUnitTypeUnknown {
		log.Errorf("Invalid data unit type (%s) for data package (%s)", packageInfo.DataUnit, packageInfo.Id)

		return nil, fmt.Errorf("invalid data unit type (%s) for data package (%s)", packageInfo.DataUnit, packageInfo.Id)
	}

	dataUnitInBytes := ukama.ReturnDataUnitsInBytes(dataUnit)
	totalData := packageInfo.DataVolume * dataUnitInBytes

	firstPackage := &sims.Package{
		PackageId:        packageId,
		InitialData:      totalData,
		IsCurrentlyInUse: true,
		DefaultDuration:  packageInfo.Duration,
	}

	err = s.packageRepo.Add(firstPackage, func(pckg *sims.Package, tx *gorm.DB) error {
		firstPackage.Id = uuid.NewV4()
		firstPackage.SimId = sim.Id

		firstPackage.StartDate = time.Now().UTC().Add(time.Minute * DefaultMinuteDelayForPackageStartDate)
		firstPackage.EndDate = validation.CalculateEndDate(firstPackage.StartDate, packageInfo.Duration)

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to add initial package to newly allocated sim. Error %s", err.Error())
	}

	sim.Package = *firstPackage

	if s.paymentClient != nil {
		_, perr := s.paymentClient.Add(cpay.AddPaymentRequest{
			ItemId:        packageId.String(),
			ItemType:      ukama.ItemTypePackage.String(),
			Amount:        fmt.Sprintf("%.2f", packageInfo.Amount),
			Currency:      packageInfo.Currency,
			PayerPhone:    remoteSubResp.Subscriber.PhoneNumber,
			PayerEmail:    remoteSubResp.Subscriber.Email,
			PayerName:     remoteSubResp.Subscriber.Name,
			Description:   packageInfo.Name,
			Country:       packageInfo.Country,
			PaymentMethod: ukama.PaymentMethodCash.String(),
			Metadata: map[string]string{
				metadataSimKey:         sim.Id.String(),
				metadataProvisionedKey: metadataProvisionedYes,
			},
		})
		if perr != nil {
			log.Errorf("failed to record package sale for allocated sim %s: %v",
				sim.Id.String(), perr)
		}
	} else {
		log.Warnf("payment client not configured: sim %s allocated with package %s "+
			"without recording a payment sale", sim.Id.String(), packageId.String())
	}

	err = pushTotalSimsCountMetric(sim.NetworkId.String(), s.simRepo, s.orgId, s.metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim allocation operation: %s", err.Error())
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), s.simRepo, s.orgId, s.metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim allocation operation: %s", err.Error())
	}

	route := s.baseRoutingKey.SetAction("allocate").SetObject("sim").MustBuild()
	evt := &epb.EventSimAllocation{
		Id:                sim.Id.String(),
		SubscriberId:      sim.SubscriberId.String(),
		NetworkId:         sim.NetworkId.String(),
		OrgId:             orgId.String(),
		DataPlanId:        sim.Package.PackageId.String(),
		Iccid:             sim.Iccid,
		Msisdn:            sim.Msisdn,
		Imsi:              sim.Imsi,
		Type:              sim.Type.String(),
		Status:            sim.Status.String(),
		IsPhysical:        sim.IsPhysical,
		TrafficPolicy:     sim.TrafficPolicy,
		SubscriberName:    remoteSubResp.Subscriber.Name,
		SubscriberEmail:   remoteSubResp.Subscriber.Email,
		NetworkName:       netInfo.Name,
		OrgName:           s.orgName,
		OwnerName:         orgOwnerName(s.orgName, s.nucleusOrgClient, s.nucleusUserClient),
		QrCode:            poolSim.QrCode,
		PackageId:         sim.Package.Id.String(),
		PackageName:       packageInfo.Name,
		PackageDataVolume: fmt.Sprintf("%v", packageInfo.DataVolume),
		PackageDataUnit:   packageInfo.DataUnit,
		PackageAmount:     fmt.Sprintf("%v", packageInfo.Amount),
		PackageDuration:   fmt.Sprintf("%v", packageInfo.Duration),
		PackageEndDate:    timestamppb.New(sim.Package.EndDate),
	}

	err = publishEventMessage(route, evt, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evt, route, err)
	}

	log.Infof("Allocating sim to subscriber success: %v", req.GetSubscriberId())

	return &pb.SimResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (s *SimManagerServer) GetSim(ctx context.Context, req *pb.SimRequest) (*pb.SimResponse, error) {
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

	return &pb.SimResponse{Sim: dbSimToPbSim(sim)}, nil
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
			"failure to get agent adapter for sim type: %q", simType)
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

func (s *SimManagerServer) ToggleSimServiceStatus(ctx context.Context, req *pb.ToggleSimServiceStatusRequest) (*pb.ToggleSimServiceStatusResponse, error) {
	log.Infof("Toggling service status for sim: %v", req.GetSimId())

	strStatus := strings.ToLower(req.Status)
	simStatus := ukama.ParseSimStatus(strStatus)

	switch simStatus {
	case ukama.SimStatusServiceOn:
		return s.setSimServiceOnForAgent(ctx, req.SimId)
	case ukama.SimStatusServiceOff:
		return s.setSimServiceOffForAgent(ctx, req.SimId)
	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid status parameter: %s.", strStatus)
	}
}

func (s *SimManagerServer) TerminateSim(ctx context.Context, req *pb.SimRequest) (*pb.TerminateSimResponse, error) {
	//Ths does not terminate the sim, but instead send a termination request to ukama agent
	//which will remove the ASR profile first, then trigger the sim termination request on sim manager.
	log.Infof("Sending terminate sim: %v request to agent", req.GetSimId())

	sim, err := getSim(req.SimId, s.simRepo)
	if err != nil {
		return nil, err
	}

	evtMsg := &epb.EventSimTermination{
		Id:           sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		Type:         sim.Type.String(),
	}

	route := s.baseRoutingKey.SetAction("terminate").SetObject("sim").MustBuild()

	err = publishEventMessage(route, evtMsg, s.msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	log.Infof("Async Sim %s terminate request sent  successfully", req.GetSimId())

	return &pb.TerminateSimResponse{}, nil
}

func (s *SimManagerServer) AddPackageForSim(ctx context.Context, req *pb.AddPackageRequest) (*pb.PackageResponse, error) {
	if err := addPackageForSim(ctx, req.SimId, req.PackageId, req.StartDate, s.simRepo, s.packageRepo, s.packageClient,
		s.orgName, s.orgId, s.metricsPusher, s.nucleusOrgClient, s.nucleusUserClient, s.subscriberRegistryService,
		s.networkClient, s.msgbus, s.baseRoutingKey); err != nil {
		return nil, err
	}

	return &pb.PackageResponse{}, nil
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
		req.FromEndDate, req.ToEndDate, req.IsCurrentlyInUse, req.IsExpired, req.Count, req.Sort)
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

func (s *SimManagerServer) RemovePackageForSim(ctx context.Context, req *pb.PackageRequest) (*pb.PackageResponse, error) {
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
			"simId packageId mismatch: package %s does not belong to the provided sim %s",
			req.GetPackageId(), req.GetSimId())

	}

	if pckg.IsCurrentlyInUse {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot remove currently in-use package (%s) from sim. Set package as not currently in-use first", pckg.Id)
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

	return &pb.PackageResponse{}, nil
}

func (s *SimManagerServer) SetPackageInUseForSim(ctx context.Context, req *pb.PackageRequest) (*pb.PackageResponse, error) {
	if err := setPackageInUseForSim(ctx, req.SimId, req.PackageId, s.simRepo, s.packageRepo, s.agentFactory, s.msgbus, s.baseRoutingKey); err != nil {
		return nil, err
	}

	return &pb.PackageResponse{}, nil
}

func (s *SimManagerServer) UnsetPackageInUseForSim(ctx context.Context, req *pb.PackageRequest) (*pb.PackageResponse, error) {
	if err := unsetPackageInuseForSim(req.SimId, req.PackageId, s.packageRepo); err != nil {
		return nil, err
	}

	return &pb.PackageResponse{}, nil
}

func (s *SimManagerServer) MarkPackageExpiredForSim(ctx context.Context, req *pb.PackageRequest) (*pb.PackageResponse, error) {
	if err := markPackageExpiredForSim(req.SimId, req.PackageId, s.simRepo, s.packageRepo, s.msgbus, s.baseRoutingKey); err != nil {
		return nil, err
	}

	return &pb.PackageResponse{}, nil
}

func (s *SimManagerServer) GenerateSimToken(ctx context.Context, req *pb.SimTokenRequest) (*pb.SimTokenResponse, error) {
	log.Infof("Generating token for sim %v", req.Iccid)

	tok, err := s.tokenCodec.GenerateTokenFromIccid(req.Iccid)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"an unexpected error has occurred while generating token for sim %s. Error: %v",
			req.Iccid, err)
	}

	return &pb.SimTokenResponse{Token: tok}, nil
}

func (s *SimManagerServer) setSimServiceOnForAgent(ctx context.Context, reqSimId string) (*pb.ToggleSimServiceStatusResponse, error) {
	sim, err := getSim(reqSimId, s.simRepo)
	if err != nil {
		return nil, err
	}

	if sim.Status != ukama.SimStatusServiceOff {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim service state: %s is invalid for turning on", sim.Status)
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q while turning service on", sim.Type, reqSimId)
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
		return nil, err
	}

	if err := setSimServiceOn(sim, s.simRepo, s.orgId, s.metricsPusher, s.msgbus, s.baseRoutingKey); err != nil {
		return nil, err
	}

	return &pb.ToggleSimServiceStatusResponse{}, nil
}

func (s *SimManagerServer) setSimServiceOffForAgent(ctx context.Context, reqSimId string) (*pb.ToggleSimServiceStatusResponse, error) {
	sim, err := getSim(reqSimId, s.simRepo)
	if err != nil {
		return nil, err
	}

	if sim.Status != ukama.SimStatusServiceOn {
		return nil, status.Errorf(codes.FailedPrecondition,
			"sim service state: %s is invalid for turning off", sim.Status)
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q while turning service off", sim.Type, reqSimId)
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
		return nil, err
	}

	if err := setSimServiceOff(sim, s.simRepo, s.orgId, s.metricsPusher, s.msgbus, s.baseRoutingKey); err != nil {
		return nil, err
	}

	return &pb.ToggleSimServiceStatusResponse{}, nil
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

func setSimServiceOn(sim *sims.Sim, simRepo sims.SimRepo, orgId string, metricsPusher MetricsPusher,
	msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	log.Infof("Turning service on for sim: %v", sim.Id)

	if sim.Status != ukama.SimStatusServiceOff {
		return status.Errorf(codes.FailedPrecondition,
			"sim service state: %s is invalid for turning on", sim.Status)
	}

	simUpdates := &sims.Sim{
		Id:               sim.Id,
		Status:           ukama.SimStatusServiceOn,
		ActivationsCount: sim.ActivationsCount + 1,
		LastActivatedOn:  time.Now(),
	}

	if sim.FirstActivatedOn.IsZero() {
		simUpdates.FirstActivatedOn = simUpdates.LastActivatedOn
	}

	err := simRepo.Update(simUpdates, nil)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "sim")
	}

	err = pushActiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim service-on operation: %s", err.Error())
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim service-on operation: %s", err.Error())
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

	log.Infof("Sim %s service turned on successfully", sim.Id)

	return nil
}

func setSimServiceOff(sim *sims.Sim, simRepo sims.SimRepo, orgId string, metricsPusher MetricsPusher,
	msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	log.Infof("Turning service off for sim: %v", sim.Id)

	if sim.Status != ukama.SimStatusServiceOn {
		return status.Errorf(codes.FailedPrecondition,
			"sim service state: %s is invalid for turning off", sim.Status)
	}

	simUpdates := &sims.Sim{
		Id:                 sim.Id,
		Status:             ukama.SimStatusServiceOff,
		DeactivationsCount: sim.DeactivationsCount + 1,
	}

	err := simRepo.Update(simUpdates, nil)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "sim")
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim service-off operation: %s", err.Error())
	}

	err = pushActiveSimsCountMetric(sim.NetworkId.String(), simRepo, orgId, metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim service-off operation: %s", err.Error())
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

	log.Infof("Sim %s turned off successfully", sim.Id)

	return nil
}

func addPackageForSim(ctx context.Context, simId, packageId, startDate string, simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	packageClient cdplan.PackageClient, orgName, orgId string, metricsPusher MetricsPusher, nucleusOrgClient cnuc.OrgClient,
	nucleusUserClient cnuc.UserClient, subscriberRegistryService providers.SubscriberRegistryClientProvider, networkClient creg.NetworkClient,
	msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	log.Infof("Adding package %v to sim: %v", packageId, simId)

	formattedStart, err := validation.ValidateDate(startDate)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	sim, err := getSim(simId, simRepo)
	if err != nil {
		return status.Errorf(codes.NotFound,
			"invalid simId while adding package to sim. Error %s", err.Error())
	}

	packageUuid, err := uuid.FromString(packageId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pkgInfo, err := packageClient.Get(packageUuid.String())
	if err != nil {
		return status.Errorf(codes.Internal,
			"failed to retrieve package from data plan system. Error %s", err.Error())
	}

	if !pkgInfo.IsActive {
		return status.Error(codes.FailedPrecondition,
			"cannot set package to sim: data plan package is no more active within its org")
	}

	// We validate package duration to make sure it is greater than 0 min and lesser than 1000 years
	if err := validation.ValidatePackageDuration(pkgInfo.Duration); err != nil {
		return status.Errorf(codes.InvalidArgument, "Error: %v",
			fmt.Errorf("invalid package duration: %w", err))
	}

	pkgInfoSimType := ukama.ParseSimType(pkgInfo.SimType)
	if sim.Type != pkgInfoSimType {
		return status.Errorf(codes.InvalidArgument,
			"invalid sim type: sim (%s) and package (%s)'s sim types mismatch",
			sim.Type, pkgInfoSimType.String())
	}

	// totalData is in bytes and pkgInfo.DataVolume depends on pack.DataUnit
	dataUnit := ukama.ParseDataUnitType(pkgInfo.DataUnit)
	if dataUnit == ukama.DataUnitTypeUnknown {
		log.Errorf("Invalid data unit type (%s) for data package (%s)", pkgInfo.DataUnit, pkgInfo.Id)

		return fmt.Errorf("invalid data unit type (%s) for data package (%s)", pkgInfo.DataUnit, pkgInfo.Id)
	}

	dataUnitInBytes := ukama.ReturnDataUnitsInBytes(dataUnit)
	totalData := pkgInfo.DataVolume * dataUnitInBytes

	pkg := &sims.Package{
		SimId:            sim.Id,
		PackageId:        packageUuid,
		InitialData:      totalData,
		IsCurrentlyInUse: false,
		DefaultDuration:  pkgInfo.Duration,
	}

	packages, err := packageRepo.List(simId, "", "", "", "", "", false, false, 0, true)
	if err != nil {
		log.Errorf("failed to get the sorted list of packages present on sim (%s): %v",
			simId, err)

		return grpc.SqlErrorToGrpc(err, "packages")
	}

	if len(packages) == 0 {
		startDate, err := time.Parse(time.RFC3339, formattedStart)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to parse start date: %v", err)
		}

		pkg.StartDate = startDate
		pkg.EndDate = validation.CalculateEndDate(pkg.StartDate, pkgInfo.Duration)
	} else {
		pkg.StartDate = packages[len(packages)-1].EndDate.Add(time.Minute * DefaultMinuteDelayForPackageStartDate)
		pkg.EndDate = validation.CalculateEndDate(pkg.StartDate, pkgInfo.Duration)
	}

	log.Infof("Package start date: %v, end date: %v", pkg.StartDate, pkg.EndDate)

	err = packageRepo.Add(pkg, func(pckg *sims.Package, tx *gorm.DB) error {
		pckg.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return grpc.SqlErrorToGrpc(err, "package")
	}

	subscriberName, subscriberEmail := subscriberNameAndEmail(ctx, sim.SubscriberId.String(), subscriberRegistryService)

	networkName := ""

	netInfo, err := networkClient.Get(sim.NetworkId.String())
	if err != nil {
		log.Errorf("Failed to get network %s while enriching sim add package event: %v",
			sim.NetworkId.String(), err)
	} else {
		networkName = netInfo.Name
	}

	route := baseRoutingKey.SetAction("addpackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimAddPackage{
		Id:              sim.Id.String(),
		SubscriberId:    sim.SubscriberId.String(),
		Iccid:           sim.Iccid,
		Imsi:            sim.Imsi,
		NetworkId:       sim.NetworkId.String(),
		PackageId:       packageUuid.String(),
		SubscriberName:  subscriberName,
		SubscriberEmail: subscriberEmail,
		NetworkName:     networkName,
		OrgName:         orgName,
		OwnerName:       orgOwnerName(orgName, nucleusOrgClient, nucleusUserClient),
		PackageName:     pkgInfo.Name,
		PackagesCount:   fmt.Sprintf("%v", len(packages)+1),
		PackagesDetails: fmt.Sprintf("$%.2f / %v %s / %d days", pkgInfo.Amount, pkgInfo.DataVolume, pkgInfo.DataUnit, pkgInfo.Duration),
		PackageEndDate:  pkg.EndDate.Format(emailDateFormat),
	}

	err = publishEventMessage(route, evtMsg, msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	return nil
}

func orgOwnerName(orgName string, nucleusOrgClient cnuc.OrgClient, nucleusUserClient cnuc.UserClient) string {
	orgInfo, err := nucleusOrgClient.Get(orgName)
	if err != nil {
		log.Errorf("Failed to get org %s while enriching event: %v", orgName, err)

		return ""
	}

	userInfo, err := nucleusUserClient.GetById(orgInfo.Owner)
	if err != nil {
		log.Errorf("Failed to get org owner %s while enriching event: %v", orgInfo.Owner, err)

		return ""
	}

	return userInfo.Name
}

func subscriberNameAndEmail(ctx context.Context, subscriberId string,
	subscriberRegistryService providers.SubscriberRegistryClientProvider) (string, string) {
	subscriberRegistrySvc, err := subscriberRegistryService.GetClient()
	if err != nil {
		log.Errorf("Failed to get subscriber registry client while enriching event: %v", err)

		return "", ""
	}

	remoteSubResp, err := subscriberRegistrySvc.Get(ctx, &subregpb.GetSubscriberRequest{
		SubscriberId: subscriberId,
	})
	if err != nil {
		log.Errorf("Failed to get subscriber %s while enriching event: %v", subscriberId, err)

		return "", ""
	}

	return remoteSubResp.Subscriber.Name, remoteSubResp.Subscriber.Email
}

func setPackageInUseForSim(ctx context.Context, reqSimId, reqPackageId string, simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	agentFactory adapters.AgentFactory, msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	log.Infof("Setting package %v as in-use for sim: %v", reqPackageId, reqSimId)

	sim, err := getSim(reqSimId, simRepo)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Package.Id != uuid.Nil {
		return status.Errorf(codes.FailedPrecondition,
			"sim currently has package %v as in-use. This package needs to be fully used or marked as expired first",
			sim.Package.Id)
	}

	packageId, err := uuid.FromString(reqPackageId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pkg, err := packageRepo.Get(packageId)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "package")
	}

	if pkg.SimId.String() != reqSimId {
		return status.Errorf(codes.InvalidArgument,
			"simId packageId mismatch: package %s does not belong to the provided sim %s",
			reqPackageId, reqSimId)
	}

	if pkg.IsExpired {
		return status.Errorf(codes.FailedPrecondition,
			"cannot set expired package (%s) as in-use", pkg.Id)
	}

	if pkg.IsCurrentlyInUse {
		return status.Errorf(codes.FailedPrecondition,
			"cannot set already in-use package (%s) as in-use", pkg.Id)
	}

	// We validate package duration to make sure it is greater than 0 min and lesser than 1000 years
	if err := validation.ValidatePackageDuration(pkg.DefaultDuration); err != nil {
		return status.Errorf(codes.InvalidArgument, "Error: %v",
			fmt.Errorf("invalid package duration: %w", err))
	}

	// Update package on Ukama Agent
	simAgent, ok := agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return status.Errorf(codes.InvalidArgument,
			"invalid sim type: %q for sim Id: %q", sim.Type, sim.Id)
	}

	agentRequest := client.AgentRequestData{
		Iccid:        sim.Iccid,
		Imsi:         sim.Imsi,
		NetworkId:    sim.NetworkId.String(),
		PackageId:    pkg.PackageId.String(),
		SimPackageId: pkg.Id.String(),
	}

	log.Infof("Updating package on remote agent for %s sim type with iccid %s",
		sim.Type.String(), sim.Iccid)

	err = simAgent.UpdatePackage(ctx, agentRequest)
	if err != nil {
		log.Errorf("Fail to update package on remote agent for %s sim type with iccid %s. Error: %v",
			sim.Type.String(), sim.Iccid, err)

		return fmt.Errorf("fail to update package on remote agent for %s sim type with iccid %s. Error: %w",
			sim.Type.String(), sim.Iccid, err)
	}

	// Update package on sim manager
	packageToSetInUse := &sims.Package{
		IsCurrentlyInUse: true,
	}

	err = packageRepo.Update([]uuid.UUID{pkg.Id}, packageToSetInUse, func(pckg *sims.Package, tx *gorm.DB) error {
		// update startDate and endDate
		packageToSetInUse.StartDate = time.Now().UTC().
			Add(time.Minute * DefaultMinuteDelayForPackageStartDate)

		packageToSetInUse.EndDate = validation.CalculateEndDate(packageToSetInUse.
			StartDate, pkg.DefaultDuration)

		return nil
	})
	if err != nil {
		log.Errorf("Fail to update package on sim manager for %s sim type with iccid %s, while successfully updated on remote agent. Error: %v",
			sim.Type.String(), sim.Iccid, err)

		return status.Errorf(codes.Internal,
			"failed to set package as in-use. Error %s", err.Error())
	}

	// Publish the event only when both updates are successful
	route := baseRoutingKey.SetAction("activepackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimActivePackage{
		Id:               sim.Id.String(),
		SubscriberId:     sim.SubscriberId.String(),
		Iccid:            sim.Iccid,
		Imsi:             sim.Imsi,
		NetworkId:        sim.NetworkId.String(),
		PackageId:        pkg.Id.String(),
		PlanId:           pkg.PackageId.String(),
		PackageStartDate: timestamppb.New(packageToSetInUse.StartDate),
		PackageEndDate:   timestamppb.New(packageToSetInUse.EndDate),
	}

	err = publishEventMessage(route, evtMsg, msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	return nil
}

func unsetPackageInuseForSim(reqSimId, reqPackageId string, packageRepo sims.PackageRepo) error {
	log.Infof("Unsetting package %v as in-use for sim: %v", reqPackageId, reqSimId)

	packageId, err := uuid.FromString(reqPackageId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pckg, err := packageRepo.Get(packageId)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "package")
	}

	if pckg.SimId.String() != reqSimId {
		return status.Errorf(codes.InvalidArgument,
			"simId packageId mismatch: package %s does not belong to the provided sim %s",
			reqPackageId, reqSimId)
	}

	if !pckg.IsCurrentlyInUse {
		return status.Errorf(codes.FailedPrecondition,
			"cannot set not in-use package (%s) as not in-use", pckg.Id)
	}

	if pckg.IsExpired {
		return status.Errorf(codes.FailedPrecondition,
			"package (%s) has already been marked as expired", pckg.Id)
	}

	inUsePackageToUnset := &sims.Package{
		IsCurrentlyInUse: false,
		// UsedDataAtExpiry; // Get the final volume from ASR,
	}

	err = packageRepo.Update([]uuid.UUID{pckg.Id}, inUsePackageToUnset, func(pckg *sims.Package, tx *gorm.DB) error {
		inUsePackageToUnset.EndDate = time.Now().UTC()

		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal,
			"failed to unset in-use package as . Error %s", err.Error())
	}

	return nil
}

func setNextQueuedPackageInUseIfIdle(ctx context.Context, simId string, simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	agentFactory adapters.AgentFactory, orgId string, metricsPusher MetricsPusher,
	msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	sim, err := getSim(simId, simRepo)
	if err != nil {
		return err
	}

	if sim.Package.Id != uuid.Nil {
		return nil
	}

	packages, err := packageRepo.List(simId, "", "", "", "", "", false, false, 0, true)
	if err != nil {
		return fmt.Errorf("failed to get the list of packages present on sim (%s): %w", simId, err)
	}

	for _, p := range packages {
		if p.IsExpired || p.IsCurrentlyInUse {
			continue
		}

		log.Infof("Sim %s is idle; setting queued package %s as in-use", simId, p.Id.String())

		if err := setPackageInUseForSim(ctx, simId, p.Id.String(), simRepo, packageRepo, agentFactory, msgbus, baseRoutingKey); err != nil {
			return err
		}

		if sim.Status != ukama.SimStatusServiceOff {
			return nil
		}

		log.Infof("Sim %s had been turned off for running out of packages; resuming service now that %s is set in use",
			simId, p.Id.String())

		sim.Package.Id = p.Id

		return setSimServiceOn(sim, simRepo, orgId, metricsPusher, msgbus, baseRoutingKey)
	}

	log.Infof("Sim %s is idle with no queued package to set in use", simId)

	if sim.Status != ukama.SimStatusServiceOn {
		return nil
	}

	log.Infof("Turning off service for sim %s: no queued package to roll over to", simId)

	return setSimServiceOff(sim, simRepo, orgId, metricsPusher, msgbus, baseRoutingKey)
}

func markPackageExpiredForSim(reqSimId, reqPackageId string, simRepo sims.SimRepo, packageRepo sims.PackageRepo,
	msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	log.Infof("Marking package %v as expired for sim: %v", reqPackageId, reqSimId)

	packageId, err := uuid.FromString(reqPackageId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	pckg, err := packageRepo.Get(packageId)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "package")
	}

	if pckg.SimId.String() != reqSimId {
		return status.Errorf(codes.InvalidArgument,
			"simId packageId mismatch: package %s does not belong to the provided sim %s",
			reqPackageId, reqSimId)
	}

	if !pckg.IsCurrentlyInUse {
		log.Warnf("cannot mark not in-use package (%s) as expired. Skipping operation.", pckg.Id)

		return status.Errorf(codes.FailedPrecondition,
			"cannot mark not in-use package (%s) as expired. Skipping operation.", pckg.Id)
	}

	if pckg.IsExpired {
		return status.Errorf(codes.FailedPrecondition,
			"package (%s) has already been marked as expired", pckg.Id)
	}

	sim, err := getSim(reqSimId, simRepo)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "sim")
	}

	packageToExpire := &sims.Package{
		IsCurrentlyInUse: false,
		IsExpired:        true,
	}

	err = packageRepo.Update([]uuid.UUID{pckg.Id}, packageToExpire, func(pckg *sims.Package, tx *gorm.DB) error {
		packageToExpire.EndDate = time.Now().UTC()

		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal,
			"failed to mark package as expired. Error %s", err.Error())
	}

	route := baseRoutingKey.SetAction("expirepackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimPackageExpire{
		Id:              sim.Id.String(),
		StartDate:       pckg.StartDate.String(),
		EndDate:         packageToExpire.EndDate.String(),
		DefaultDuration: pckg.DefaultDuration,
		PackageId:       pckg.Id.String(),
		DataPlanId:      pckg.PackageId.String(),
	}

	err = publishEventMessage(route, evtMsg, msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	return nil
}

func expireAllPackagesForSim(simId string, packageRepo sims.PackageRepo, msgbus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	packages, err := packageRepo.List(simId, "", "", "", "", "", false, false, 0, true)
	if err != nil {
		return fmt.Errorf("failed to get the list of packages present on sim (%s): %w", simId, err)
	}

	var idsToExpire []uuid.UUID
	var activePckg *sims.Package

	for i, pckg := range packages {
		if pckg.IsExpired {
			continue
		}

		idsToExpire = append(idsToExpire, pckg.Id)

		if pckg.IsCurrentlyInUse {
			activePckg = &packages[i]
		}
	}

	if len(idsToExpire) == 0 {
		return nil
	}

	packageToExpire := &sims.Package{
		IsCurrentlyInUse: false,
		IsExpired:        true,
	}

	err = packageRepo.Update(idsToExpire, packageToExpire, func(pckg *sims.Package, tx *gorm.DB) error {
		packageToExpire.EndDate = time.Now().UTC()

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to mark packages as expired on sim (%s): %w", simId, err)
	}

	endDate := packageToExpire.EndDate

	log.Infof("Marked packages %v as expired on sim %s", idsToExpire, simId)

	if activePckg == nil {
		return nil
	}

	route := baseRoutingKey.SetAction("expirepackage").SetObject("sim").MustBuild()
	evtMsg := &epb.EventSimPackageExpire{
		Id:              simId,
		StartDate:       activePckg.StartDate.String(),
		EndDate:         endDate.String(),
		DefaultDuration: activePckg.DefaultDuration,
		PackageId:       activePckg.Id.String(),
		DataPlanId:      activePckg.PackageId.String(),
	}

	err = publishEventMessage(route, evtMsg, msgbus)
	if err != nil {
		log.Errorf(eventPublishErrorMsg, evtMsg, route, err)
	}

	return nil
}

func pushTotalSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId string, metricsPusher MetricsPusher) error {
	log.Infof("Collecting and pushing total sims count metric to push gateway host: %s", metricsPusher.GetPushMetricsHost())

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusUnknown, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting total sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting total sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = metricsPusher.CollectAndPushSimMetrics(pkg.SimMetric, pkg.NumberOfSims,
		float64(len(sims)), map[string]string{"network": networkId, "org": orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing total sims count metric to push gateway: %s", err.Error())

		return fmt.Errorf("error while pushing total sims count metric to push gateway: %w", err)
	}

	return nil
}

func pushActiveSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId string, metricsPusher MetricsPusher) error {
	log.Infof("Collecting and pushing active sims count metric to push gateway host: %s", metricsPusher.GetPushMetricsHost())

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusServiceOn, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting active sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting active sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = metricsPusher.CollectAndPushSimMetrics(pkg.SimMetric, pkg.ActiveCount,
		float64(len(sims)), map[string]string{"network": networkId, "org": orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while active sims count metric to push gateway: %s", err.Error())

		return fmt.Errorf("error while pushing active sims count metric to push gateway: %w", err)
	}

	return nil
}

func pushInactiveSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId string, metricsPusher MetricsPusher) error {
	log.Infof("Collecting and pushing inactive sims count metric to push gateway host: %s", metricsPusher.GetPushMetricsHost())

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusServiceOff, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting inactive sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting inactive sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = metricsPusher.CollectAndPushSimMetrics(pkg.SimMetric, pkg.InactiveCount,
		float64(len(sims)), map[string]string{"network": networkId, "org": orgId}, pkg.SystemName)
	if err != nil {
		log.Errorf("Error while pushing inactive sims count metric to push gateway: %s", err.Error())

		return fmt.Errorf("error while pushing inactive sims count metric to push gateway: %w", err)
	}

	return nil
}

func pushTerminatedSimsCountMetric(networkId string, simRepo sims.SimRepo, orgId string, metricsPusher MetricsPusher) error {
	log.Infof("Collecting and pushing terminated sims count metric to push gateway host: %s", metricsPusher.GetPushMetricsHost())

	sims, err := simRepo.List("", "", "", networkId, ukama.SimTypeUnknown, ukama.SimStatusTerminated, 0, false, 0, false)
	if err != nil {
		log.Errorf("Error while collecting terminated sims count metric for network: %s. Error: %v",
			networkId, err)

		return fmt.Errorf("error while collecting terminated sims count metric for network: %s. Error: %w",
			networkId, grpc.SqlErrorToGrpc(err, "sims"))
	}

	err = metricsPusher.CollectAndPushSimMetrics(pkg.SimMetric, pkg.TerminatedCount,
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
		Id:               pkg.Id.String(),
		PackageId:        pkg.PackageId.String(),
		InitialData:      pkg.InitialData,
		UsedDataAtExpiry: pkg.UsedDataAtExpiry,
		IsCurrentlyInUse: pkg.IsCurrentlyInUse,
		IsExpired:        pkg.IsExpired,
		DefaultDuration:  pkg.DefaultDuration,
		CreatedAt:        pkg.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        pkg.UpdatedAt.Format(time.RFC3339),
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

type MetricsPusher interface {
	CollectAndPushSimMetrics(configMetrics []pmetric.MetricConfig,
		selectedMetric string, Value float64, Labels map[string]string, systemName string) error
	GetPushMetricsHost() string
}

type metricsPusher struct {
	pushMetricHost string
}

func NewMetricsPusher(host string) MetricsPusher {
	return &metricsPusher{
		pushMetricHost: host,
	}
}

func (m *metricsPusher) CollectAndPushSimMetrics(configMetrics []pmetric.MetricConfig,
	selectedMetric string, value float64, labels map[string]string, systemName string) error {
	return pmetric.CollectAndPushSystemMetrics(m.pushMetricHost, configMetrics, selectedMetric,
		value, labels, systemName)

}

func (m *metricsPusher) GetPushMetricsHost() string {
	return m.pushMetricHost
}
