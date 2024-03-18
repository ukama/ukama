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
	pb "github.com/ukama/ukama/systems/services/gitClient/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GitClientProvider interface {
	GetClient() (pb.GitClientServiceClient, error)
}

type gitClientProvider struct {
	gitClientService pb.GitClientServiceClient
	gitHost          string
}

func NewGitClientProvider(gitHost string) GitClientProvider {
	return &gitClientProvider{gitHost: gitHost}
}

func (u *gitClientProvider) GetClient() (pb.GitClientServiceClient, error) {
	if u.gitClientService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to git client service ", u.gitHost)

		conn, err := grpc.DialContext(ctx, u.gitHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Org service %s. Error: %v", u.gitHost, err)
		}

		u.gitClientService = pb.NewGitClientServiceClient(conn)
	}

	return u.gitClientService, nil
}
