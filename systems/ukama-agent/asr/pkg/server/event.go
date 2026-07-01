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
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cpb "github.com/ukama/ukama/systems/common/pb/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

const (
	handlerTimeoutFactor = 3
)

type AsrEventServer struct {
	asrRepo        db.AsrRecordRepo
	gutiRepo       db.GutiRepo
	s              *AsrRecordServer
	network        registry.NetworkClient
	factory        factory.SimFactoryClient
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pc             pm.Controller
	allowedToS     int64
	orgName        string
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

func (as *AsrEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(as.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := cpb.UnmarshalProtoEvent[epb.CDRReported](e.Msg)
		if err != nil {
			log.Errorf("Error while unmarshaling CDRReported event proto: %v", err)

			return nil, fmt.Errorf("error while unmarshaling CDRReported event proto: %w", err)
		}

		err = as.handleEventCDRCreate(e.RoutingKey, msg)
		if err != nil {
			log.Errorf("Error while handling CDR create Event: %v", err)

			return nil, fmt.Errorf("error while handling CDR create Event: %w", err)
		}

	case msgbus.PrepareRoute(as.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := cpb.UnmarshalProtoEvent[epb.EventSimAllocation](e.Msg)
		if err != nil {
			log.Errorf("Error while unmarshaling EventSimAllocation proto: %v", err)

			return nil, fmt.Errorf("error while unmarshaling EventSimAllocation proto: %w", err)
		}

		err = as.handleSimManagerSimAllocateEvent(e.RoutingKey, msg)
		if err != nil {
			log.Errorf("Error while handling sim manage SimAllocate Event: %v", err)

			return nil, fmt.Errorf("error while handling sim manage SimAllocate Event: %w", err)
		}
	default:
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (as *AsrEventServer) handleEventCDRCreate(key string, cdr *epb.CDRReported) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	err := as.s.UpdateAndSyncAsrProfile(cdr.GetImsi())
	if err != nil {
		log.Errorf("Failed to update the active subscriber %s. Error: %v", cdr.Imsi, err)

		return fmt.Errorf("eailed to update the active subscriber %s. Error: %w", cdr.Imsi, err)
	}

	return nil
}

func (as *AsrEventServer) handleSimManagerSimAllocateEvent(key string, sim *epb.EventSimAllocation) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	if sim.Type != ukama.SimTypeUkamaData.String() {
		log.Infof("Sim type %s is not supported by ukama agent. Skipping...", sim.Type)

		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	_, err := activate(ctx, sim.Iccid, sim.Imsi, sim.PackageId, sim.DataPlanId, sim.NetworkId,
		as.network, as.factory, as.asrRepo, as.pc, as.allowedToS, as.msgbus, as.baseRoutingKey)
	if err != nil {
		log.Errorf("Failed to activate sim %s. Error: %v", sim.Imsi, err)

		//TODO: publish activation failure for rollback on sim manager and sim pool if necessary

		return fmt.Errorf("failed to activate sim %s. Error: %w", sim.Imsi, err)
	}

	return nil
}
