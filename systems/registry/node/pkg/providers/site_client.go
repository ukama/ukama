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
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SiteClientProvider creates a local client to interact with

type SiteClientProvider interface {
	GetClient() (pb.SiteServiceClient, error)
}

type siteClientProvider struct {
	siteService pb.SiteServiceClient
	siteHost    string
}

func NewSiteClientProvider(siteHost string) SiteClientProvider {
	return &siteClientProvider{siteHost: siteHost}
}

func (o *siteClientProvider) GetClient() (pb.SiteServiceClient, error) {
	if o.siteService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Site service ", o.siteHost)

		conn, err := grpc.NewClient(o.siteHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Site service %s. Error: %v", o.siteHost, err)

			return nil, fmt.Errorf("failed to connect to remote site service: %w", err)
		}

		o.siteService = pb.NewSiteServiceClient(conn)
	}

	return o.siteService, nil
}
