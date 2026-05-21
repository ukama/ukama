/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package providers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type ControllerClientProvider interface {
	GetClient() (pb.ControllerServiceClient, error)
}

type controllerClientProvider struct {
	controllerService pb.ControllerServiceClient
	controllerHost    string
}

func NewControllerClientProvider(controllerHost string) ControllerClientProvider {
	return &controllerClientProvider{controllerHost: controllerHost}
}

func (o *controllerClientProvider) GetClient() (pb.ControllerServiceClient, error) {
	if o.controllerService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Controller service ", o.controllerHost)

		conn, err := grpc.NewClient(o.controllerHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Controller service %s. Error: %v", o.controllerHost, err)

			return nil, fmt.Errorf("failed to connect to remote controller service: %w", err)
		}

		o.controllerService = pb.NewControllerServiceClient(conn)
	}

	return o.controllerService, nil
}
