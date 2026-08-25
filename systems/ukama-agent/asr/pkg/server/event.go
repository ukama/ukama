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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/utils"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

type AsrEventServer struct {
	asrRepo    db.AsrRecordRepo
	gutiRepo   db.GutiRepo
	s          *AsrRecordServer
	network    registry.NetworkClient
	factory    factory.SimFactoryClient
	msgbus     mb.MsgBusServiceClient
	pc         pm.Controller
	allowedToS int64
	orgName    string
	epb.UnimplementedEventNotificationServiceServer
}

func NewAsrEventServer(asrRepo db.AsrRecordRepo, s *AsrRecordServer, gutiRepo db.GutiRepo, factory factory.SimFactoryClient,
	network registry.NetworkClient, pc pm.Controller, msgBus mb.MsgBusServiceClient, aToS int64, org string) *AsrEventServer {
	return &AsrEventServer{
		asrRepo:    asrRepo,
		gutiRepo:   gutiRepo,
		factory:    factory,
		network:    network,
		msgbus:     msgBus,
		pc:         pc,
		allowedToS: aToS,
		orgName:    org,
		s:          s,
	}
}

func (ae *AsrEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(ae.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := cpb.UnmarshalProtoEvent[epb.CDRReported](e.Msg)
		if err != nil {
			log.Errorf("Error while unmarshaling CDRReported event proto: %v", err)

			return nil, fmt.Errorf("error while unmarshaling CDRReported event proto: %w", err)
		}

		err = ae.handleEventCDRCreate(e.RoutingKey, msg)
		if err != nil {
			log.Errorf("Error while handling CDR create Event: %v", err)

			return nil, fmt.Errorf("error while handling CDR create Event: %w", err)
		}

	case msgbus.PrepareRoute(ae.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimAllocation](e.Msg)
		if err != nil {
			log.Errorf("Error while unmarshaling EventSimAllocation proto: %v", err)

			return nil, fmt.Errorf("error while unmarshaling EventSimAllocation proto: %w", err)
		}

		err = ae.handleSimManagerSimAllocateEvent(e.RoutingKey, msg)
		if err != nil {
			log.Errorf("Error while handling sim manage SimAllocate Event: %v", err)

			return nil, fmt.Errorf("error while handling sim manage SimAllocate Event: %w", err)
		}

	case msgbus.PrepareRoute(ae.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.terminate"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimTermination](e.Msg)
		if err != nil {
			log.Errorf("Error while unmarshaling EventSimTermination proto: %v", err)

			return nil, fmt.Errorf("error while unmarshaling EventSimTermination proto: %w", err)
		}

		err = ae.handleSimManagerSimTerminateEvent(e.RoutingKey, msg)
		if err != nil {
			log.Errorf("Error while handling sim manage SimTerminate Event: %v", err)

			return nil, fmt.Errorf("error while handling sim manage SimTerminate Event: %w", err)
		}

	case msgbus.PrepareRoute(ae.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.serviceon"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimServiceOn](e.Msg)
		if err != nil {
			log.Errorf("Error while unmarshaling EventSimServiceOn proto: %v", err)

			return nil, fmt.Errorf("error while unmarshaling EventSimServiceOn proto: %w", err)
		}

		err = ae.handleSimManagerSimServiceOnEvent(e.RoutingKey, msg)
		if err != nil {
			log.Errorf("Error while handling sim manager SimServiceOn Event: %v", err)

			return nil, fmt.Errorf("error while handling sim manager SimServiceOn Event: %w", err)
		}

	case msgbus.PrepareRoute(ae.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.serviceoff"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimServiceOff](e.Msg)
		if err != nil {
			log.Errorf("Error while unmarshaling EventSimServiceOff proto: %v", err)

			return nil, fmt.Errorf("error while unmarshaling EventSimServiceOff proto: %w", err)
		}

		err = ae.handleSimManagerSimServiceOffEvent(e.RoutingKey, msg)
		if err != nil {
			log.Errorf("Error while handling sim manager SimServiceOff Event: %v", err)

			return nil, fmt.Errorf("error while handling sim manager SimServiceOff Event: %w", err)
		}
	default:
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (ae *AsrEventServer) handleEventCDRCreate(key string, cdr *epb.CDRReported) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	err := ae.s.UpdateAndSyncAsrProfileFromCdr(cdr.GetImsi())
	if err != nil {
		log.Errorf("Failed to update the active subscriber %s. Error: %v", cdr.Imsi, err)

		return fmt.Errorf("eailed to update the active subscriber %s. Error: %w", cdr.Imsi, err)
	}

	return nil
}

func (ae *AsrEventServer) handleSimManagerSimAllocateEvent(key string, sim *epb.EventSimAllocation) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	if sim.Type != ukama.SimTypeUkamaData.String() {
		log.Infof("Sim type %s is not supported by ukama agent. Skipping...", sim.Type)

		return nil
	}

	_, err := createProfile(sim.Iccid, sim.Imsi, sim.PackageId, sim.DataPlanId, sim.NetworkId, pm.PolicyInput{
		TotalData: sim.PackageTotalData,
		Dlbr:      sim.PackageDlbr,
		Ulbr:      sim.PackageUlbr,
		StartTime: uint64(sim.PackageStartDate.AsTime().Unix()),
		EndTime:   uint64(sim.PackageEndDate.AsTime().Unix()),
	}, ae.network, ae.factory, ae.asrRepo, ae.pc, ae.allowedToS)
	if err != nil {
		log.Errorf("Failed to create subscriber profile for sim %s. Error: %v", sim.Imsi, err)

		// Publish ASR activation failure for rollback on sim manager and sim pool if necessary
		e := &epb.AsrProfileDeleted{
			Subscriber: &epb.Subscriber{
				Iccid: sim.Iccid,
			},
		}

		if err := utils.PublishEvent(ae.orgName, msgbus.ACTION_CRUD_DELETE, activeSubscriberEventObject, e,
			ae.msgbus); err != nil {
			log.Warnf("Failed to publish subcriber delete as rollback event for sim %s allocation failure. Error: %v", sim.Imsi, err)
		}

		return fmt.Errorf("failed to create profile for sim %s. Error: %v", sim.Imsi, err)
	}

	return nil
}

func (ae *AsrEventServer) handleSimManagerSimServiceOnEvent(key string, sim *epb.EventSimServiceOn) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	if err := updateServiceStatus(sim.Iccid, true, ae.asrRepo, ae.pc); err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			log.Infof("No ASR profile for sim %s; skipping service status update", sim.Iccid)

			return nil
		}

		log.Errorf("Failed to turn on ASR service status for sim %s. Error: %v", sim.Iccid, err)

		return fmt.Errorf("failed to turn on ASR service status for sim %s. Error: %w", sim.Iccid, err)
	}

	return nil
}

func (ae *AsrEventServer) handleSimManagerSimServiceOffEvent(key string, sim *epb.EventSimServiceOff) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	if err := updateServiceStatus(sim.Iccid, false, ae.asrRepo, ae.pc); err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			log.Infof("No ASR profile for sim %s; skipping service status update", sim.Iccid)

			return nil
		}

		log.Errorf("Failed to turn off ASR service status for sim %s. Error: %v", sim.Iccid, err)

		return fmt.Errorf("failed to turn off ASR service status for sim %s. Error: %w", sim.Iccid, err)
	}

	return nil
}

func (ae *AsrEventServer) handleSimManagerSimTerminateEvent(key string, sim *epb.EventSimTermination) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	if sim.Type != ukama.SimTypeUkamaData.String() {
		log.Infof("Sim type %s is not supported by ukama agent. Skipping...", sim.Type)

		return nil
	}

	_, err := deleteProfile(sim.Iccid, ae.asrRepo, ae.pc)
	if err != nil {
		log.Errorf("Failed to remove subscriber profile for sim %s. Error: %v", sim.Imsi, err)

		return fmt.Errorf("failed to remove subscriber profile for sim %s. Error: %v", sim.Imsi, err)
	}

	return nil
}
