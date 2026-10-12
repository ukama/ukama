/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/sql"

	log "github.com/sirupsen/logrus"
)

// Returns grpc error with code based on Sql error.
// Handles cases such as "not found" or "duplicate key"
func SqlErrorToGrpc(err error, entity string) error {
	if err == nil {
		return nil
	}

	log.Error(err)

	if sql.IsNotFoundError(err) {
		return status.Errorf(codes.NotFound, "%s record not found", entity)
	}

	if sql.IsDuplicateKeyError(err) {
		return status.Errorf(codes.AlreadyExists, "%s already exists", entity)
	}

	return status.Error(codes.Internal, err.Error())
}

func CreateGrpcConn(conf config.GrpcService) *grpc.ClientConn {
	var conn *grpc.ClientConn

	log.Infof("Connecting to service: %s", conf.Host)

	conn, err := grpc.NewClient(conf.Host, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				MinConnectTimeout: conf.Timeout,
			}))
	if err != nil {
		log.Fatalf("Failed to connect to service %s. Error: %v", conf.Host, err)
	}
	return conn
}
