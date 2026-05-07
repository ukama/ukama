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
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type HealthClientProvider interface {
	GetClient() (pb.HealhtServiceClient, error)
}

type healthClientProvider struct {
	healthService pb.HealhtServiceClient
	healthHost    string
}

func NewHealthClientProvider(healthHost string) HealthClientProvider {
	return &healthClientProvider{healthHost: healthHost}
}

func (o *healthClientProvider) GetClient() (pb.HealhtServiceClient, error) {
	if o.healthService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Health service ", o.healthHost)

		conn, err := grpc.NewClient(o.healthHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Health service %s. Error: %v", o.healthHost, err)

			return nil, fmt.Errorf("failed to connect to remote health service: %w", err)
		}

		o.healthService = pb.NewHealhtServiceClient(conn)
	}

	return o.healthService, nil
}
