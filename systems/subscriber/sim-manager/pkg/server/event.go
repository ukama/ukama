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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

const (
	handlerTimeoutFactor = 3
)

type SimManagerEventServer struct {
	simRepo        sims.SimRepo
	agentFactory   adapters.AgentFactory
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	orgId          string
	orgName        string
	pushMetricHost string
	s              *SimManagerServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimManagerEventServer(orgName, orgId string, simRepo sims.SimRepo, agentFactory adapters.AgentFactory,
	msgBus mb.MsgBusServiceClient, pushMetricHost string, s *SimManagerServer) *SimManagerEventServer {
	return &SimManagerEventServer{
		simRepo:      simRepo,
		agentFactory: agentFactory,
		msgbus:       msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).
			SetOrgName(orgName).SetService(pkg.ServiceName),
		orgName:        orgName,
		orgId:          orgId,
		pushMetricHost: pushMetricHost,
		s:              s,
	}
}

func (es *SimManagerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := epb.UnmarshalEventSimAllocation(e.Msg, "EventSimAllocate")
		if err != nil {
			return nil, err
		}

		err = es.handleSimManagerSimAllocateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.payments.processor.payment.success"):
		msg, err := unmarshalProcessorPaymentSuccess(e.Msg)
		if err != nil {
			return nil, err
		}

		paymentStatus := ukama.ParseStatusType(msg.Status)
		itemType := ukama.ParseItemType(msg.ItemType)

		if paymentStatus == ukama.StatusTypeCompleted && itemType == ukama.ItemTypePackage {
			err = handleProcessorPaymentSuccessEvent(e.RoutingKey, msg, es.s)
			if err != nil {
				return nil, err
			}
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.operator.cdr.cdr.create"):
		msg, err := unmarshalOperatorCdrCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleOperatorCdrCreateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := unmarshalUkamaAgentCdrCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleUkamaAgentCdrCreateEvent(e.RoutingKey, msg, es.s)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete"):
		msg, err := unmarshalUkamaAgentAsrProfileDelete(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleUkamaAgentAsrProfileDeleteEvent(e.RoutingKey, msg, es.s)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

// We auto activate any new allocated sim
func (es *SimManagerEventServer) handleSimManagerSimAllocateEvent(key string, msg *epb.EventSimAllocation) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	return activateSim(ctx, msg.Id, es.simRepo, es.agentFactory, es.orgId, es.pushMetricHost, es.msgbus, es.baseRoutingKey)
}

func handleProcessorPaymentSuccessEvent(key string, msg *epb.Payment, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

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

	addReq := &pb.AddPackageRequest{
		SimId:     simId,
		PackageId: msg.ItemId,
		StartDate: time.Now().UTC().Format(time.RFC3339),
	}

	log.Infof("Adding package %s to sim %s", addReq.PackageId, addReq.SimId)

	_, err = s.AddPackageForSim(ctx, addReq)

	return err
}

func (es *SimManagerEventServer) handleOperatorCdrCreateEvent(key string, cdr *epb.EventOperatorCdrReport) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	if cdr.Type != ukama.CdrTypeData.String() {
		log.Warnf("Unsupported CDR Type (%s) received for data usage count. Skipping", cdr.Type)

		return nil
	}

	sims, err := es.simRepo.List(cdr.Iccid, "", "", "", ukama.SimTypeOperatorData, ukama.SimStatusActive, 0, false, 0, false)
	if err != nil {
		return fmt.Errorf("error while looking up sim for given iccid %q: %w",
			cdr.Iccid, err)
	}

	if len(sims) == 0 {
		return fmt.Errorf("no corresponding active sim found for given iccid %q",
			cdr.Iccid)
	}

	if len(sims) > 1 {
		return fmt.Errorf("inconsistent state: multiple sim found for given iccid %q",
			cdr.Iccid)
	}

	sim := sims[0]

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

	return err
}

func handleUkamaAgentCdrCreateEvent(key string, cdr *epb.CDRReported, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	sims, err := s.simRepo.List("", cdr.Imsi, "", "", ukama.SimTypeUkamaData, ukama.SimStatusActive, 0, false, 0, false)
	if err != nil {
		return fmt.Errorf("error while looking up sim for given imsi %q: %w",
			cdr.Imsi, err)
	}

	if len(sims) == 0 {
		return fmt.Errorf("no corresponding sim found for given imsi %q",
			cdr.Imsi)
	}

	if len(sims) > 1 {
		return fmt.Errorf("inconsistent state: multiple sim found for given imsi %q",
			cdr.Imsi)
	}

	sim := sims[0]

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

	route := s.baseRoutingKey.SetAction("usage").SetObject("sim").MustBuild()

	err = s.msgbus.PublishRequest(route, usageMsg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			usageMsg, route, err.Error())
	}

	return err
}

func handleUkamaAgentAsrProfileDeleteEvent(key string, asrProfile *epb.Profile, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, asrProfile)

	sims, err := s.simRepo.List(asrProfile.Iccid, "", "", "", ukama.SimTypeUkamaData, ukama.SimStatusActive, 0, false, 0, false)
	if err != nil {
		return fmt.Errorf("error while looking up sim for given iccid %q: %w",
			asrProfile.Iccid, err)
	}

	if len(sims) == 0 {
		return fmt.Errorf("no corresponding sim found for given iccid %q",
			asrProfile.Iccid)
	}

	if len(sims) > 1 {
		return fmt.Errorf("inconsistent state: multiple sim found for given iccid %q",
			asrProfile.Iccid)
	}

	sim := sims[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	termReq := &pb.TerminatePackageRequest{
		SimId:     sim.Id.String(),
		PackageId: asrProfile.SimPackage,
	}

	log.Infof("terminating package %s on sim %s", termReq.PackageId, termReq.SimId)

	_, err = s.TerminatePackageForSim(ctx, termReq)
	if err != nil {
		return fmt.Errorf("failed to terminate active package %s on sim %s. Error: %w",
			termReq.PackageId, termReq.SimId, err)
	}

	// Get next package to activate if any
	packages, err := s.packageRepo.List(termReq.SimId, "", "", "", "", "", false, false, 0, true)
	if err != nil {
		log.Errorf("failed to get the sorted list of packages present on sim (%s): %v",
			termReq.SimId, err)

		return fmt.Errorf("failed to get the sorted list of packages present on sim (%s): %w",
			termReq.SimId, err)
	}

	if len(packages) > 1 {
		var p db.Package

		var i int
		for i, p = range packages {
			if p.Id.String() == termReq.PackageId {
				break
			}
		}

		if i <= len(packages)-2 {
			nextPackage := packages[i+1]

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
			defer cancel()

			activeReq := &pb.SetActivePackageRequest{
				SimId:     sim.Id.String(),
				PackageId: nextPackage.Id.String(),
			}

			log.Infof("activating package %s on sim %s", activeReq.PackageId, activeReq.SimId)

			_, err = s.SetActivePackageForSim(ctx, activeReq)
			if err != nil {
				return fmt.Errorf("failed to activate next package %s for sim %s. Error: %w",
					activeReq.PackageId, activeReq.SimId, err)
			}
		}
	}

	return nil
}

func unmarshalProcessorPaymentSuccess(msg *anypb.Any) (*epb.Payment, error) {
	p := &epb.Payment{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal payment message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalOperatorCdrCreate(msg *anypb.Any) (*epb.EventOperatorCdrReport, error) {
	p := &epb.EventOperatorCdrReport{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal EventOperatorCdrReport message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalUkamaAgentCdrCreate(msg *anypb.Any) (*epb.CDRReported, error) {
	p := &epb.CDRReported{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal UkamaAgent CDRReported message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalUkamaAgentAsrProfileDelete(msg *anypb.Any) (*epb.Profile, error) {
	p := &epb.Profile{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal UkamaAgent ASR profile message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}
