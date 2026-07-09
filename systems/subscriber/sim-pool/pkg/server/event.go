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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
)

type SimPoolEventServer struct {
	orgName     string
	simPoolRepo db.SimRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimPoolEventServer(orgName string, simPoolRepo db.SimRepo) *SimPoolEventServer {
	return &SimPoolEventServer{
		orgName:     orgName,
		simPoolRepo: simPoolRepo,
	}
}
func (l *SimPoolEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := epb.UnmarshalEventSimAllocation(e.Msg, "EventSimAllocate")
		if err != nil {
			return nil, err
		}
		err = handleEventCloudSimManagerSimAllocate(e.RoutingKey, msg, l)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("handler not registered for %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func handleEventCloudSimManagerSimAllocate(key string, msg *epb.EventSimAllocation, l *SimPoolEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	err := l.simPoolRepo.UpdateStatus(msg.Iccid, true, false)
	if err != nil {
		return err
	}
	return err
}
