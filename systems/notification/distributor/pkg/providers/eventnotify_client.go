/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventNotifyClientProvider interface {
	GetClient() (pb.EventToNotifyServiceClient, error)
}

type eventNotifyClientProvider struct {
	eventNotifyService pb.EventToNotifyServiceClient
	eventNotifyHost    string
}

func NewEventNotifyClientProvider(eventNotifyHost string) EventNotifyClientProvider {
	return &eventNotifyClientProvider{eventNotifyHost: eventNotifyHost}
}

func (o *eventNotifyClientProvider) GetClient() (pb.EventToNotifyServiceClient, error) {
	if o.eventNotifyService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to event-notify service ", o.eventNotifyHost)

		conn, err := grpc.NewClient(o.eventNotifyHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to event-notify service %s. Error: %v", o.eventNotifyHost, err)

			return nil, fmt.Errorf("failed to connect to remote event-notify service: %w", err)
		}

		o.eventNotifyService = pb.NewEventToNotifyServiceClient(conn)
	}

	return o.eventNotifyService, nil
}
