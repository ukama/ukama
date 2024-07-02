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
	pb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrgClientProvider creates a local client to interact with
// a remote instance of  Org service.
type OrgClientProvider interface {
	GetClient() (pb.OrgServiceClient, error)
}

type orgClientProvider struct {
	orgService pb.OrgServiceClient
	orgHost    string
}

func NewOrgClientProvider(orgHost string) OrgClientProvider {
	return &orgClientProvider{orgHost: orgHost}
}

func (u *orgClientProvider) GetClient() (pb.OrgServiceClient, error) {
	if u.orgService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Org service ", u.orgHost)

		conn, err := grpc.NewClient(u.orgHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Org service %s. Error: %v", u.orgHost, err)
		}

		u.orgService = pb.NewOrgServiceClient(conn)
	}

	return u.orgService, nil
}
