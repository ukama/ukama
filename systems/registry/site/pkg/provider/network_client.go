/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NetworkClientProvider interface {
	GetClient() (pb.NetworkServiceClient, error)
}

type networkClientProvider struct {
	networkService pb.NetworkServiceClient
	networkHost    string
}

func NewNetworkClientProvider(networkHost string) NetworkClientProvider {
	return &networkClientProvider{networkHost: networkHost}
}

func (u *networkClientProvider) GetClient() (pb.NetworkServiceClient, error) {
	if u.networkService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Network service ", u.networkHost)

		conn, err := grpc.NewClient(u.networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Network service %s. Error: %v", u.networkHost, err)
		}

		u.networkService = pb.NewNetworkServiceClient(conn)
	}

	return u.networkService, nil
}
