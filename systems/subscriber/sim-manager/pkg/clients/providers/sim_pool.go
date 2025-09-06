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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

// SimPoolClientProvider creates a local client to interact with
// a remote instance of Sim Pool service.
type SimPoolClientProvider interface {
	GetClient() (pb.SimServiceClient, error)
}

type simPoolClientProvider struct {
	simPoolService pb.SimServiceClient
	simPoolHost    string
	timeout        time.Duration
}

func NewSimPoolClientProvider(simPoolHost string, timeout time.Duration) SimPoolClientProvider {
	return &simPoolClientProvider{simPoolHost: simPoolHost, timeout: timeout}
}

func (sp *simPoolClientProvider) GetClient() (pb.SimServiceClient, error) {
	if sp.simPoolService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Sim Pool service ", sp.simPoolHost)

		conn, err := grpc.NewClient(sp.simPoolHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Sim Pool service %s. Error: %v", sp.simPoolHost, err)

			return nil, fmt.Errorf("failed to connect to remote Sim Pool service: %w", err)
		}

		sp.simPoolService = pb.NewSimServiceClient(conn)
	}

	return sp.simPoolService, nil
}
