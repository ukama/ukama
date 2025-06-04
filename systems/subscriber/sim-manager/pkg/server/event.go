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
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

const (
	handlerTimeoutFactor = 3
)

type SimManagerEventServer struct {
	orgName string
	s       *SimManagerServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimManagerEventServer(orgName string, s *SimManagerServer) *SimManagerEventServer {
	return &SimManagerEventServer{
		orgName: orgName,
		s:       s,
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

		err = handleEventCloudSimManagerSimAllocate(e.RoutingKey, msg, es.s)
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
			err = handleEventCloudProcessorPaymentSuccess(e.RoutingKey, msg, es.s)
			if err != nil {
				return nil, err
			}
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.subscriber.asr_cleanup_completed"):
		msg, err := unmarshalEventSimAsrCleanupCompleted(e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleAsrCleanupCompleted(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.operator.cdr.cdr.create"):
		msg, err := unmarshalOperatorCdrCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudOperatorCdrCreate(e.RoutingKey, msg, es.s)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := unmarshalUkamaAgentCdrCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudUkamaAgentCdrCreate(e.RoutingKey, msg, es.s)
		if err != nil {
			return nil, err
		}
		
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete"):
		msg, err := unmarshalUkamaAgentAsrProfileDelete(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudUkamaAgentAsrProfileDelete(e.RoutingKey, msg, es.s)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

// We auto activate any new allocated sim
func handleEventCloudSimManagerSimAllocate(key string, msg *epb.EventSimAllocation, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	_, err := s.activateSim(ctx, msg.Id)

	return err
}

func handleEventCloudProcessorPaymentSuccess(key string, msg *epb.Payment, s *SimManagerServer) error {
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
func handleEventCloudOperatorCdrCreate(key string, cdr *epb.EventOperatorCdrReport, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	if cdr.Type != ukama.CdrTypeData.String() {
		log.Warnf("Unsupported CDR Type (%s) received for data usage count. Skipping", cdr.Type)
		return nil
	}

	sims, err := s.simRepo.List(cdr.Iccid, "", "", "", ukama.SimTypeOperatorData, ukama.SimStatusUnknown, 0, false, 0, false)
	if err != nil {
		return fmt.Errorf("error while looking up sim for given iccid %q: %w",
			cdr.Iccid, err)
	}

	if len(sims) == 0 {
		return fmt.Errorf("no corresponding sim found for given iccid %q",
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
	}

	route := s.baseRoutingKey.SetAction("usage").SetObject("sim").MustBuild()

	err = s.msgbus.PublishRequest(route, usageMsg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			usageMsg, route, err.Error())
	}

	return err
}
func handleEventCloudUkamaAgentCdrCreate(key string, cdr *epb.CDRReported, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	// FIXED: Changed ukama.SimStatusActive to ukama.SimStatusUnknown to include inactive SIMs
	sims, err := s.simRepo.List("", cdr.Imsi, "", "", ukama.SimTypeUkamaData, ukama.SimStatusUnknown, 0, false, 0, false)
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
	}

	route := s.baseRoutingKey.SetAction("usage").SetObject("sim").MustBuild()

	err = s.msgbus.PublishRequest(route, usageMsg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			usageMsg, route, err.Error())
	}

	return err
}
func handleEventCloudUkamaAgentAsrProfileDelete(key string, asrProfile *epb.Profile, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, asrProfile)

	// FIXED: Changed ukama.SimStatusActive to ukama.SimStatusUnknown to include inactive SIMs
	sims, err := s.simRepo.List(asrProfile.Iccid, "", "", "", ukama.SimTypeUkamaData, ukama.SimStatusUnknown, 0, false, 0, false)
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

	// Rest of the function remains the same...
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
func (es *SimManagerEventServer) handleAsrCleanupCompleted(ctx context.Context, routingKey string, msg *epb.EventSimAsrCleanupCompleted) error {
	log.Infof("Received ASR cleanup completion from ASR service. SubscriberId: %s", msg.SubscriberId)

	if msg.SimResults != nil {
		successCount := 0
		for _, result := range msg.SimResults {
			if result.Success {
				successCount++
			} else {
				log.Warnf("ASR cleanup failed for SIM %s (ICCID: %s)", 
					result.SimId, result.Iccid)
			}
		}
		log.Infof("ASR cleanup summary: %d successful, %d failed", 
			successCount, len(msg.SimResults)-successCount)
	}

	log.Infof("Starting SIM and package cleanup for subscriber: %s", msg.SubscriberId)
	
	_, err := es.s.TerminateSimsForSubscriber(ctx, &pb.TerminateSimsForSubscriberRequest{
		SubscriberId: msg.SubscriberId,
	})
	if err != nil {
		log.Errorf("Failed to terminate SIMs for subscriber %s: %v", msg.SubscriberId, err)

	} else {
		log.Infof("Successfully completed SIM and package cleanup for subscriber: %s", msg.SubscriberId)
	}

	err = es.publishSubscriberDeletionCompleted(msg.SubscriberId)
	if err != nil {
		log.Errorf("Failed to publish subscriber deletion completion: %v", err)
		return err
	}

	log.Infof("Successfully completed full deletion flow for subscriber: %s", msg.SubscriberId)
	return nil
}
func (es *SimManagerEventServer) publishSubscriberDeletionCompleted(subscriberId string) error {
	completionEvent := &epb.EventSubscriberDeletionCompleted{
		SubscriberId: subscriberId,
	}

	route := es.s.baseRoutingKey.
		SetAction("deletion_completed").
		SetObject("sims").
		MustBuild()

	log.Infof("Publishing subscriber deletion completion to Registry at %s: %+v", route, completionEvent)

	err := es.s.PublishEventMessage(route, completionEvent)
	if err != nil {
		log.Errorf("Failed to publish subscriber deletion completion: %v", err)
		return err
	}

	log.Infof("Successfully published deletion completion for subscriber: %s", subscriberId)
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


func unmarshalEventSimAsrCleanupCompleted(msg *anypb.Any) (*epb.EventSimAsrCleanupCompleted, error) {
	p := &epb.EventSimAsrCleanupCompleted{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal EventSimAsrCleanupCompleted message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}