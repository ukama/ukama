/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SimManagerClientProvider creates a local client to interact with
// a remote instance of  Org service.
type SimManagerClientProvider interface {
	GetSimManagerService() (pb.SimManagerServiceClient, error)
}

type simManagerClientProvider struct {
	simManagerService pb.SimManagerServiceClient
	simManagerHost    string
}

func NewSimManagerClientProvider(simManagerHost string) SimManagerClientProvider {
	return &simManagerClientProvider{simManagerHost: simManagerHost}
}

func (u *simManagerClientProvider) GetSimManagerService() (pb.SimManagerServiceClient, error) {
	if u.simManagerService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to SimManager service ", u.simManagerHost)

		conn, err := grpc.NewClient(u.simManagerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to SimManager service %s. Error: %v", u.simManagerHost, err)

			return nil, fmt.Errorf("failed to connect to remote SimManager service: %w", err)
		}

		u.simManagerService = pb.NewSimManagerServiceClient(conn)
	}

	return u.simManagerService, nil
}
