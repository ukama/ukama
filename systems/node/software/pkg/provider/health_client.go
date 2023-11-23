/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	"context"
	"time"

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

func (u *healthClientProvider) GetClient() (pb.HealhtServiceClient, error) {
	if u.healthService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to Health service ", u.healthHost)

		conn, err := grpc.DialContext(ctx, u.healthHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Health service %s. Error: %v", u.healthHost, err)
		}

		u.healthService = pb.NewHealhtServiceClient(conn)
	}

	return u.healthService, nil
}
