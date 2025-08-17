/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
)

type Distributor interface {
	GetNotificationStream(ctx context.Context, orgId, networkId, subscriberId, userId string,
		scopes []string) (pb.DistributorService_GetNotificationStreamClient, error)
}

type distributor struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.DistributorServiceClient
	host    string
}

func NewDistributor(host string, timeout time.Duration) (*distributor, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Distributor Service: %v", err)

		return nil, err
	}

	client := pb.NewDistributorServiceClient(conn)

	return &distributor{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}, nil
}

func NewDistributorFromClient(client pb.DistributorServiceClient) *distributor {
	return &distributor{
		host:    "localhost",
		timeout: 10 * time.Second,
		conn:    nil,
		client:  client,
	}
}

func (d *distributor) Close() {
	if d.conn != nil {
		if err := d.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Distributor Service connection: %v", err)
		}
	}
}

func (d *distributor) GetNotificationStream(ctx context.Context, orgId, networkId, subscriberId, userId string,
	scopes []string) (pb.DistributorService_GetNotificationStreamClient, error) {

	return d.client.GetNotificationStream(ctx,
		&pb.NotificationStreamRequest{
			OrgId:        orgId,
			NetworkId:    networkId,
			SubscriberId: subscriberId,
			UserId:       userId,
			Scopes:       scopes,
		})
}
