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

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	cnotif "github.com/ukama/ukama/systems/common/rest/client/notification"
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
	mailerClient              cnotif.MailerClient
	subscriberRegistryService providers.SubscriberRegistryClientProvider
	msgbus                    mb.MsgBusServiceClient
	baseRoutingKey            msgbus.RoutingKeyBuilder
	orgId                     string
	orgName                   string
	metricsPusher             MetricsPusher
	s                         *SimManagerServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimManagerEventServer(orgName, orgId string, simRepo sims.SimRepo, packageRepo sims.PackageRepo, agentFactory adapters.AgentFactory,
	packageClient cdplan.PackageClient, subscriberRegistryService providers.SubscriberRegistryClientProvider,
	networkClient creg.NetworkClient, mailerClient cnotif.MailerClient, nucleusOrgClient cnuc.OrgClient,
	nucleusUserClient cnuc.UserClient, msgBus mb.MsgBusServiceClient, pushMetricHost string, s *SimManagerServer) *SimManagerEventServer {
	return &SimManagerEventServer{
		simRepo:                   simRepo,
		packageRepo:               packageRepo,
		agentFactory:              agentFactory,
		packageClient:             packageClient,
		networkClient:             networkClient,
		nucleusOrgClient:          nucleusOrgClient,
		nucleusUserClient:         nucleusUserClient,
		mailerClient:              mailerClient,
		subscriberRegistryService: subscriberRegistryService,
		msgbus:                    msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).
			SetOrgName(orgName).SetService(pkg.ServiceName),
		orgName:       orgName,
		orgId:         orgId,
		metricsPusher: NewMetricsPusher(pushMetricHost),
		s:             s,
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
		msg, err := cpb.UnmarshalProtoEvent[epb.Profile](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleUkamaAgentAsrProfileCreateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete"):
		msg, err := cpb.UnmarshalProtoEvent[epb.Profile](e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleUkamaAgentAsrProfileDeleteEvent(e.RoutingKey, msg)
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

	simId, ok := metadata["sim"]
	if !ok {
		return fmt.Errorf("missing sim metadata for successful package payment: %s", msg.ItemId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	startDate := time.Now().UTC().Format(time.RFC3339)

	log.Infof("Adding package %s to sim %s", msg.ItemId, simId)

	return addPackageForSim(ctx, simId, msg.ItemId, startDate, es.simRepo, es.packageRepo, es.packageClient,
		es.orgName, es.orgId, es.metricsPusher, es.nucleusOrgClient, es.nucleusUserClient,
		es.subscriberRegistryService, es.networkClient, es.mailerClient, es.msgbus, es.baseRoutingKey)
}

func (es *SimManagerEventServer) handleOperatorCdrCreateEvent(key string, cdr *epb.EventOperatorCdrReport) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	if cdr.Type != ukama.CdrTypeData.String() {
		log.Warnf("Unsupported CDR Type (%s) received for data usage count. Skipping", cdr.Type)

		return nil
	}

	operatorSims, err := es.simRepo.List(cdr.Iccid, "", "", "", ukama.SimTypeOperatorData, ukama.SimStatusActive, 0, false, 0, false)
	if err != nil {
		return fmt.Errorf("error while looking up sim for given iccid %q: %w",
			cdr.Iccid, err)
	}

	if len(operatorSims) == 0 {
		return fmt.Errorf("no corresponding active sim found for given iccid %q",
			cdr.Iccid)
	}

	if len(operatorSims) > 1 {
		return fmt.Errorf("inconsistent state: multiple sim found for given iccid %q",
			cdr.Iccid)
	}

	sim := operatorSims[0]

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

	ukamaSims, err := es.simRepo.List("", cdr.Imsi, "", "", ukama.SimTypeUkamaData, ukama.SimStatusActive, 0, false, 0, false)
	if err != nil {
		return fmt.Errorf("error while looking up sim for given imsi %q: %w",
			cdr.Imsi, err)
	}

	if len(ukamaSims) == 0 {
		return fmt.Errorf("no corresponding sim found for given imsi %q",
			cdr.Imsi)
	}

	if len(ukamaSims) > 1 {
		return fmt.Errorf("inconsistent state: multiple sim found for given imsi %q",
			cdr.Imsi)
	}

	sim := ukamaSims[0]

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

// We activate any new allocated sim as long as ARS registration was successful
func (es *SimManagerEventServer) handleUkamaAgentAsrProfileCreateEvent(key string, asrProfile *epb.Profile) error {
	log.Infof("Keys %s and Proto is: %+v", key, asrProfile)

	sim, err := es.getSimFromIccid(asrProfile.Iccid)
	if err != nil {
		log.Errorf("Error while looking up sim from iccid %s. Error: %v",
			asrProfile.Iccid, err)

		return fmt.Errorf("error while looking up sim from iccid %s. Error: %w",
			asrProfile.Iccid, err)
	}

	if sim.Type != ukama.SimTypeUkamaData {
		log.Errorf("Invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())

		return fmt.Errorf("invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	return activateSim(ctx, sim.Id.String(), es.simRepo, es.agentFactory, es.orgId, es.metricsPusher, es.msgbus, es.baseRoutingKey)
}

func (es *SimManagerEventServer) handleUkamaAgentAsrProfileDeleteEvent(key string, asrProfile *epb.Profile) error {
	log.Infof("Keys %s and Proto is: %+v", key, asrProfile)

	sim, err := es.getSimFromIccid(asrProfile.Iccid)
	if err != nil {
		log.Errorf("Error while looking up sim from iccid %s. Error: %v",
			asrProfile.Iccid, err)

		return fmt.Errorf("error while looking up sim from iccid %s. Error: %w",
			asrProfile.Iccid, err)
	}

	if sim.Type != ukama.SimTypeUkamaData {
		log.Errorf("Invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())

		return fmt.Errorf("invalid sim type: sim must be of type %s, not %s",
			ukama.SimTypeUkamaData.String(), sim.Type.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	log.Infof("terminating package %s on sim %s", asrProfile.SimPackage, sim.Id.String())

	err = terminatePackageForSim(ctx, sim.Id.String(), asrProfile.SimPackage, es.simRepo,
		es.packageRepo, es.msgbus, es.baseRoutingKey)
	if err != nil {
		log.Errorf("Failed to terminate active package %s on sim %s. Error: %v",
			asrProfile.SimPackage, sim.Id.String(), err)

		return fmt.Errorf("failed to terminate active package %s on sim %s. Error: %w",
			asrProfile.SimPackage, sim.Id.String(), err)
	}

	// Get next package to activate if any
	packages, err := es.packageRepo.List(sim.Id.String(), "", "", "", "", "", false, false, 0, true)
	if err != nil {
		log.Errorf("failed to get the sorted list of packages present on sim (%s): %v",
			sim.Id.String(), err)

		return fmt.Errorf("failed to get the sorted list of packages present on sim (%s): %w",
			sim.Id.String(), err)
	}

	if len(packages) > 1 {
		var p sims.Package

		var i int
		for i, p = range packages {
			if p.Id.String() == asrProfile.SimPackage {
				break
			}
		}

		if i <= len(packages)-2 {
			nextPackage := packages[i+1]

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
			defer cancel()

			log.Infof("activating package %s on sim %s", nextPackage.Id.String(), sim.Id.String())

			err = setActivePackageForSim(ctx, sim.Id.String(), nextPackage.Id.String(), es.simRepo, es.packageRepo,
				es.agentFactory, es.msgbus, es.baseRoutingKey)
			if err != nil {
				log.Errorf("Failed to activate next package %s for sim %s. Error: %v",
					nextPackage.Id.String(), sim.Id.String(), err)

				return fmt.Errorf("failed to activate next package %s for sim %s. Error: %w",
					nextPackage.Id.String(), sim.Id.String(), err)
			}
		}
	}

	return nil
}

func (es *SimManagerEventServer) getSimFromIccid(iccid string) (*sims.Sim, error) {
	ukamaSims, err := es.simRepo.List(iccid, "", "", "", ukama.SimTypeUnknown, ukama.SimStatusUnknown, 0, false, 0, false)
	if err != nil {
		log.Errorf("Sim list error for given iccid %q: %v",
			iccid, err)

		return nil, fmt.Errorf("sim list error for given iccid %q: %w",
			iccid, err)
	}

	if len(ukamaSims) == 0 {
		log.Errorf("No corresponding sim found for given iccid %q",
			iccid)

		return nil, fmt.Errorf("no corresponding sim found for given iccid %q",
			iccid)
	}

	if len(ukamaSims) > 1 {
		log.Errorf("Inconsistent state: multiple sims found for given iccid %q",
			iccid)

		return nil, fmt.Errorf("inconsistent state: multiple sims found for given iccid %q",
			iccid)
	}

	sim := ukamaSims[0]

	return &sim, nil
}
