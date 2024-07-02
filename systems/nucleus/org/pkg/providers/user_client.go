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
	pb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClientProvider creates a local client to interact with
// a remote instance of  Users service.
type UserClientProvider interface {
	GetClient() (pb.UserServiceClient, error)
}

type userClientProvider struct {
	userService pb.UserServiceClient
	userHost    string
}

func NewUserClientProvider(userHost string) UserClientProvider {
	return &userClientProvider{userHost: userHost}
}

func (u *userClientProvider) GetClient() (pb.UserServiceClient, error) {
	if u.userService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to users service ", u.userHost)

		conn, err := grpc.NewClient(u.userHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to users service %s. Error: %v", u.userHost, err)
		}

		u.userService = pb.NewUserServiceClient(conn)
	}

	return u.userService, nil
}
