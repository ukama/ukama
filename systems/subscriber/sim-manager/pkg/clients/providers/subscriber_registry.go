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
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
)

// SubscriberRegistryClientProvider creates a local client to interact with
// a remote instance of Subscriber Registry service.
type SubscriberRegistryClientProvider interface {
	GetClient() (pb.RegistryServiceClient, error)
}

type subscriberRegistryClientProvider struct {
	subscriberRegistryService pb.RegistryServiceClient
	subscriberRegistryHost    string
	timeout                   time.Duration
}

func NewSubscriberRegistryClientProvider(subscriberRegistryHost string, timeout time.Duration) SubscriberRegistryClientProvider {
	return &subscriberRegistryClientProvider{subscriberRegistryHost: subscriberRegistryHost, timeout: timeout}
}

func (p *subscriberRegistryClientProvider) GetClient() (pb.RegistryServiceClient, error) {
	if p.subscriberRegistryService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Subscriber Registry service ", p.subscriberRegistryHost)

		conn, err := grpc.NewClient(p.subscriberRegistryHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Subscriber Registry service %s. Error: %v", p.subscriberRegistryHost, err)

			return nil, fmt.Errorf("failed to connect to remote Subscriber Registry service: %w", err)
		}

		p.subscriberRegistryService = pb.NewRegistryServiceClient(conn)
	}

	return p.subscriberRegistryService, nil
}
