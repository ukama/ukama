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
	"encoding/json"
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	cnuc "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

const (
	handlerTimeoutFactor = 3
)

type SimManagerEventServer struct {
	simRepo                   sims.SimRepo
	packageRepo               sims.PackageRepo
	agentFactory              adapters.AgentFactory
	packageClient             cdplan.PackageClient
	networkClient             creg.NetworkClient
	nucleusOrgClient          cnuc.OrgClient
	nucleusUserClient         cnuc.UserClient
	subscriberRegistryService providers.SubscriberRegistryClientProvider
	msgbus                    mb.MsgBusServiceClient
	baseRoutingKey            msgbus.RoutingKeyBuilder
	orgId                     string
	orgName                   string
	metricsPusher             MetricsPusher
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimManagerEventServer(orgName, orgId string, simRepo sims.SimRepo, packageRepo sims.PackageRepo, agentFactory adapters.AgentFactory,
	packageClient cdplan.PackageClient, subscriberRegistryService providers.SubscriberRegistryClientProvider,
	networkClient creg.NetworkClient, nucleusOrgClient cnuc.OrgClient,
	nucleusUserClient cnuc.UserClient, msgBus mb.MsgBusServiceClient, pushMetricHost string) *SimManagerEventServer {
	return &SimManagerEventServer{
		simRepo:                   simRepo,
		packageRepo:               packageRepo,
		agentFactory:              agentFactory,
		packageClient:             packageClient,
		networkClient:             networkClient,
		nucleusOrgClient:          nucleusOrgClient,
		nucleusUserClient:         nucleusUserClient,
		subscriberRegistryService: subscriberRegistryService,
		msgbus:                    msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).
			SetOrgName(orgName).SetService(pkg.ServiceName),
		orgName:       orgName,
		orgId:         orgId,
		metricsPusher: NewMetricsPusher(pushMetricHost),
	}
}

func (es *SimManagerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.payments.processor.payment.success"):
		msg, err := cpb.UnmarshalProtoEvent[epb.Payment](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleProcessorPaymentSuccessEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.operator.cdr.cdr.create"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventOperatorCdrReport](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleOperatorCdrCreateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := cpb.UnmarshalProtoEvent[epb.CDRReported](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleUkamaAgentCdrCreateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.create"):
		msg, err := cpb.UnmarshalProtoEvent[epb.AsrActivated](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleUkamaAgentAsrProfileCreateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete"):
		msg, err := cpb.UnmarshalProtoEvent[epb.AsrInactivated](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleUkamaAgentAsrProfileDeleteEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.policy.violation"):
		msg, err := cpb.UnmarshalProtoEvent[epb.PolicyViolation](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleUkamaAgentAsrPolicyViolationEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (es *SimManagerEventServer) handleProcessorPaymentSuccessEvent(key string, msg *epb.Payment) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	paymentStatus := ukama.ParseStatusType(msg.Status)
	itemType := ukama.ParseItemType(msg.ItemType)

	if paymentStatus != ukama.StatusTypeCompleted || itemType != ukama.ItemTypePackage {
		return fmt.Errorf("payment of %s with status %s is not valid for event handling",
			paymentStatus, itemType)
	}

	metadata := map[string]string{}

	err := json.Unmarshal(msg.Metadata, &metadata)
	if err != nil {
		return fmt.Errorf("failed to Unmarshal payment metadata as map[string]string: %w", err)
	}

	simId, ok := metadata[metadataSimKey]
	if !ok {
		return fmt.Errorf("missing sim metadata for successful package payment: %s", msg.ItemId)
	}

	// Payments recorded by sim allocation already provisioned the package
	// (added inline in AllocateSim), so skip adding it again here. The payment
	// is still a valid sale record; only the package-add is suppressed.
	if metadata[metadataProvisionedKey] == metadataProvisionedYes {
		log.Infof("payment for sim %s is already provisioned; skipping package add", simId)

		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	startDate := time.Now().UTC().Format(time.RFC3339)

	log.Infof("Adding package %s to sim %s", msg.ItemId, simId)

	return addPackageForSim(ctx, simId, msg.ItemId, startDate, es.simRepo, es.packageRepo, es.packageClient,
		es.orgName, es.orgId, es.metricsPusher, es.nucleusOrgClient, es.nucleusUserClient,
		es.subscriberRegistryService, es.networkClient, es.msgbus, es.baseRoutingKey)
}

func (es *SimManagerEventServer) handleOperatorCdrCreateEvent(key string, cdr *epb.EventOperatorCdrReport) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	if cdr.Type != ukama.CdrTypeData.String() {
		log.Errorf("Invalid cdr type: cdr must be of type %s, not %s",
			ukama.CdrTypeData.String(), cdr.Type)

		return fmt.Errorf("invalid cdr type: cdr must be of type %s, not %s",
			ukama.CdrTypeData.String(), cdr.Type)
	}

	sim, err := es.getSimFromIccidOrImsi(cdr.Iccid, "")
	if err != nil {
		log.Errorf("Error while looking up sim for operator agent CDR create event. Error: %v",
			err)

		return fmt.Errorf("error while looking up sim for operator agent CDR create event. Error: %w",
			err)
	}

	if sim.Type != ukama.SimTypeOperatorData {
		log.Errorf("Invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeOperatorData.String(), sim.Type.String())

		return fmt.Errorf("invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeOperatorData.String(), sim.Type.String())
	}

	usageMsg := &epb.EventSimUsage{
		SimId:        sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		NetworkId:    sim.NetworkId.String(),
		Type:         cdr.Type,
		BytesUsed:    cdr.Duration,
		StartTime:    cdr.ConnectTime,
		EndTime:      cdr.CloseTime,
		Id:           cdr.Id,
		// OrgId:        s.OrgId.String(),
		// SessionId: msg.InventoryId,
	}

	route := es.baseRoutingKey.SetAction("usage").SetObject("sim").MustBuild()

	err = es.msgbus.PublishRequest(route, usageMsg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			usageMsg, route, err.Error())
	}

	return nil
}

func (es *SimManagerEventServer) handleUkamaAgentCdrCreateEvent(key string, cdr *epb.CDRReported) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	sim, err := es.getSimFromIccidOrImsi("", cdr.Imsi)
	if err != nil {
		log.Errorf("Error while looking up sim for ukama agent CDR create event. Error: %v",
			err)

		return fmt.Errorf("error while looking up sim for ukama agent CDR create event. Error: %w",
			err)
	}

	if sim.Type != ukama.SimTypeUkamaData {
		log.Errorf("Invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())

		return fmt.Errorf("invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())
	}

	usageMsg := &epb.EventSimUsage{
		SimId:        sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		NetworkId:    sim.NetworkId.String(),
		Type:         ukama.CdrTypeData.String(),
		BytesUsed:    cdr.TotalBytes,
		StartTime:    cdr.StartTime,
		EndTime:      cdr.EndTime,
		// Id:           cdr.Id,
		// OrgId:        s.OrgId.String(),
		// SessionId:    cdr.Session,
	}

	route := es.baseRoutingKey.SetAction("usage").SetObject("sim").MustBuild()

	err = es.msgbus.PublishRequest(route, usageMsg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			usageMsg, route, err.Error())
	}

	return nil
}

func (es *SimManagerEventServer) handleUkamaAgentAsrProfileCreateEvent(key string, asrProfile *epb.AsrActivated) error {
	log.Infof("Keys %s and Proto is: %+v", key, asrProfile)

	sim, err := es.getSimFromIccidOrImsi(asrProfile.Subscriber.Iccid, "")
	if err != nil {
		log.Errorf("Error while looking up sim for ukama agent ASR create event. Error: %v",
			err)

		return fmt.Errorf("error while looking up sim for ukama agent ASR create event. Error: %w",
			err)
	}

	if sim.Type != ukama.SimTypeUkamaData {
		log.Errorf("Invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())

		return fmt.Errorf("invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())
	}

	return setSimServiceOn(sim, es.simRepo, es.orgId, es.metricsPusher, es.msgbus, es.baseRoutingKey)
}

func (es *SimManagerEventServer) handleUkamaAgentAsrProfileDeleteEvent(key string, asrProfile *epb.AsrInactivated) error {
	log.Infof("Keys %s and Proto is: %+v", key, asrProfile)

	sim, err := es.getSimFromIccidOrImsi(asrProfile.Subscriber.Iccid, "")
	if err != nil {
		log.Errorf("Error while looking up sim for ukama agent ASR delete event. Error: %v",
			err)

		return fmt.Errorf("error while looking up sim for ukama agent ASR delte event. Error: %w",
			err)
	}

	if sim.Type != ukama.SimTypeUkamaData {
		log.Errorf("Invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())

		return fmt.Errorf("invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())
	}

	log.Infof("Terminating all packages on sim %s", sim.Id.String())

	err = expireAllPackagesForSim(sim.Id.String(), es.packageRepo, es.msgbus, es.baseRoutingKey)
	if err != nil {
		log.Errorf("Failed to expire packages on sim %s. Error: %v", sim.Id.String(), err)

		return fmt.Errorf("failed to expire packages on sim %s. Error: %w", sim.Id.String(), err)
	}

	// Terminating sim

	simUpdates := &sims.Sim{
		Id:     sim.Id,
		Status: ukama.SimStatusTerminated,
	}

	err = es.simRepo.Update(simUpdates, func(sim *sims.Sim, tx *gorm.DB) error {
		sim.TerminatedAt = time.Now().UTC()

		return nil
	})
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "sim")
	}

	err = pushTerminatedSimsCountMetric(sim.NetworkId.String(), es.simRepo, es.orgId, es.metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim terminate operation: %s", err.Error())
	}

	err = pushInactiveSimsCountMetric(sim.NetworkId.String(), es.simRepo, es.orgId, es.metricsPusher)
	if err != nil {
		log.Errorf("Error while pushing metrics on sim terminate operation: %s", err.Error())
	}

	return nil
}

func (es *SimManagerEventServer) handleUkamaAgentAsrPolicyViolationEvent(key string, violation *epb.PolicyViolation) error {
	log.Infof("Keys %s and Proto is: %+v", key, violation)

	sim, err := es.getSimFromIccidOrImsi(violation.Profile.Iccid, "")
	if err != nil {
		log.Errorf("Error while looking up sim for ukama agent policy violation event. Error: %v",
			err)

		return fmt.Errorf("error while looking up sim for ukama agent policy violation event. Error: %w",
			err)
	}

	if sim.Type != ukama.SimTypeUkamaData {
		log.Errorf("Invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())

		return fmt.Errorf("invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())
	}

	packageId, err := uuid.FromString(violation.Profile.SimPackage)
	if err != nil {
		log.Errorf("Invalid sim package uuid %q on policy violation for sim %s. Error: %v",
			violation.Profile.SimPackage, sim.Id.String(), err)

		return fmt.Errorf("invalid sim package uuid %q on policy violation for sim %s. Error: %w",
			violation.Profile.SimPackage, sim.Id.String(), err)
	}

	err = es.packageRepo.Update([]uuid.UUID{packageId},
		&sims.Package{UsedDataAtExpiry: violation.Profile.TotalDataBytes}, nil)
	if err != nil {
		log.Errorf("Failed to record usage for package %s on sim %s. Error: %v",
			packageId, sim.Id.String(), err)

		return fmt.Errorf("failed to record usage for package %s on sim %s. Error: %w",
			packageId, sim.Id.String(), err)
	}

	reason := ukama.ParsePolicyViolationReason(violation.Reason)
	if reason != ukama.PolicyViolationReasonPackageExpired {
		log.Infof("Policy violation (%s) recorded for sim %s. No further action for this reason yet.",
			reason, sim.Id.String())

		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	log.Infof("Expiring package %s on sim %s", packageId, sim.Id.String())

	err = markPackageExpiredForSim(sim.Id.String(), packageId.String(), es.simRepo,
		es.packageRepo, es.msgbus, es.baseRoutingKey)
	if err != nil {
		log.Errorf("Failed to mark package %s as expired on sim %s. Error: %v",
			packageId, sim.Id.String(), err)

		return fmt.Errorf("failed to mark package %s as expired on sim %s. Error: %w",
			packageId, sim.Id.String(), err)
	}

	// Get next package to set in-use, if any
	packages, err := es.packageRepo.List(sim.Id.String(), "", "", "", "", "", false, false, 0, true)
	if err != nil {
		log.Errorf("failed to get the sorted list of packages present on sim (%s): %v",
			sim.Id.String(), err)

		return fmt.Errorf("failed to get the sorted list of packages present on sim (%s): %w",
			sim.Id.String(), err)
	}

	var i int
	var p sims.Package

	for i, p = range packages {
		if p.Id == packageId {
			break
		}
	}

	if i > len(packages)-2 {
		log.Warnf("Sim %s has no queued package to roll over to after package %s expired.",
			sim.Id.String(), packageId.String())

		return nil
	}

	nextPackage := packages[i+1]

	log.Infof("Setting next package %s as in-use on sim %s", nextPackage.Id.String(), sim.Id.String())

	err = setPackageInUseForSim(ctx, sim.Id.String(), nextPackage.Id.String(), es.simRepo, es.packageRepo,
		es.agentFactory, es.msgbus, es.baseRoutingKey)
	if err != nil {
		log.Errorf("Failed to set next package %s as in-use for sim %s. Error: %v",
			nextPackage.Id.String(), sim.Id.String(), err)

		return fmt.Errorf("failed to set next package %s as in-use for sim %s. Error: %w",
			nextPackage.Id.String(), sim.Id.String(), err)
	}

	return nil
}

func (es *SimManagerEventServer) getSimFromIccidOrImsi(iccid, imsi string) (*sims.Sim, error) {
	ukamaSims, err := es.simRepo.List(iccid, imsi, "", "", ukama.SimTypeUnknown, ukama.SimStatusUnknown, 0, false, 0, false)
	if err != nil {
		log.Errorf("Sim list error for given (iccid: %q, imisi: %q). Error:: %v",
			iccid, imsi, err)

		return nil, fmt.Errorf("sim list error for given (iccid: %q, imisi: %q). Error:: %w",
			iccid, imsi, err)
	}

	if len(ukamaSims) == 0 {
		log.Errorf("No corresponding sim found for given (iccid %q, imsi: %q)",
			iccid, imsi)

		return nil, fmt.Errorf("no corresponding sim found for given (iccid: %q, imsi: %q)",
			iccid, imsi)
	}

	if len(ukamaSims) > 1 {
		log.Errorf("Inconsistent state: multiple sims found for given (iccid: %q, imsi: %q)",
			iccid, imsi)

		return nil, fmt.Errorf("inconsistent state: multiple sims found for given (iccid %q, imsi: %q)",
			iccid, imsi)
	}

	sim := ukamaSims[0]

	return &sim, nil
}
