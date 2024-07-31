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
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
)

type RateClientProvider interface {
	GetClient() (pb.RateServiceClient, error)
}

type rateClientProvider struct {
	rateService pb.RateServiceClient
	rateHost    string
	timeout     time.Duration
}

func NewRateClientProvider(rateHost string, timeout time.Duration) RateClientProvider {
	return &rateClientProvider{rateHost: rateHost, timeout: timeout}
}

func (rt *rateClientProvider) GetClient() (pb.RateServiceClient, error) {
	if rt.rateService == nil {
		var conn *grpc.ClientConn

		log.Infoln("Connecting to Rate service ", rt.rateHost)

		conn, err := grpc.NewClient(rt.rateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Rate service %s. Error: %v", rt.rateHost, err)

			return nil, fmt.Errorf("failed to connect to remote Rate service: %w", err)
		}

		rt.rateService = pb.NewRateServiceClient(conn)
	}

	return rt.rateService, nil
}
