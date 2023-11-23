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
	pb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type EventServer struct {
	pb.UnimplementedEventNotificationServiceServer
}

func NewEventServer() *EventServer {
	return &EventServer{}
}

func (e *EventServer) EventNotification(ctx context.Context, evt *pb.Event) (*pb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", evt.RoutingKey, evt.Msg)
	return &pb.EventResponse{}, nil
}
