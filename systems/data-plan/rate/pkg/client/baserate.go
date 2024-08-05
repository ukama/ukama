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

	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaserateClientProvider interface {
	GetClient() (bpb.BaseRatesServiceClient, error)
}

type baserateClientProvider struct {
	baserateService bpb.BaseRatesServiceClient
	baserateHost    string
	timeout         time.Duration
}

func NewBaseRateClientProvider(baserateHost string, timeout time.Duration) BaserateClientProvider {
	return &baserateClientProvider{baserateHost: baserateHost, timeout: timeout}
}

func (bs *baserateClientProvider) GetClient() (bpb.BaseRatesServiceClient, error) {
	if bs.baserateService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Rate service ", bs.baserateHost)

		conn, err := grpc.NewClient(bs.baserateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Rate service %s. Error: %v", bs.baserateHost, err)

			return nil, fmt.Errorf("failed to connect to remote Rate service: %w", err)
		}

		bs.baserateService = bpb.NewBaseRatesServiceClient(conn)
	}

	return bs.baserateService, nil
}
