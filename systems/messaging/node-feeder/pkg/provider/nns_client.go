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
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NnsClientProvider interface {
	GetClient() (pb.NnsClient, error)
}

type nnsClientProvider struct {
	nnsClient pb.NnsClient
	nnsHost    string
}

func NewNnsClientProvider(nnsHost string) NnsClientProvider {
	return &nnsClientProvider{nnsHost: nnsHost}
}

func (u *nnsClientProvider) GetClient() (pb.NnsClient, error) {
	if u.nnsClient == nil {
		var conn *grpc.ClientConn

		conn, err := grpc.NewClient(u.nnsHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to NNS service %s. Error: %v", u.nnsHost, err)
		}

		u.nnsClient = pb.NewNnsClient(conn)
	}

	return u.nnsClient, nil
}
