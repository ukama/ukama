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
	pb "github.com/ukama/ukama/testing/services/dummy/dsimfactory/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DsimfactoryProvider interface {
	GetClient() (pb.DsimfactoryServiceClient, error)
}

type dsimfactoryProvider struct {
	siteService pb.DsimfactoryServiceClient
	host        string
}

func NewDsimfactoryProvider(host string) DsimfactoryProvider {
	return &dsimfactoryProvider{host: host}
}

func (o *dsimfactoryProvider) GetClient() (pb.DsimfactoryServiceClient, error) {
	if o.siteService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Dsimfactory service ", o.host)

		conn, err := grpc.NewClient(o.host,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Dsimfactory service %s. Error: %v", o.host, err)

			return nil, fmt.Errorf("failed to connect to remote site service: %w", err)
		}

		o.siteService = pb.NewDsimfactoryServiceClient(conn)
	}

	return o.siteService, nil
}
