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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
)

type LookupClientProvider interface {
	GetClient() (pb.LookupServiceClient, error)
}

type lookupClientProvider struct {
	lookupService pb.LookupServiceClient
	lookupHost    string
	timeout       time.Duration
}

func NewLookupClientProvider(lookupHost string, timeout time.Duration) LookupClientProvider {
	return &lookupClientProvider{lookupHost: lookupHost, timeout: timeout}
}

func (rt *lookupClientProvider) GetClient() (pb.LookupServiceClient, error) {
	if rt.lookupService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Lookup service ", rt.lookupHost)

		conn, err := grpc.NewClient(rt.lookupHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Lookup service %s. Error: %v", rt.lookupHost, err)

			return nil, fmt.Errorf("failed to connect to remote Lookup service: %w", err)
		}

		rt.lookupService = pb.NewLookupServiceClient(conn)
	}

	return rt.lookupService, nil
}
