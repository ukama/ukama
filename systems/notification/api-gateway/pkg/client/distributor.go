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
	GetNotificationStream(orgId string, networkId string, subscriberId string, userId string, scopes []string) (pb.DistributorService_GetNotificationStreamClient, error)
}

type distributor struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.DistributorServiceClient
	host    string
}

func NewDistributor(host string, timeout time.Duration) (*distributor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)

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

func (m *distributor) Close() {
	m.conn.Close()
}

func (n *distributor) GetNotificationStream(orgId string, networkId string, subscriberId string, userId string, scopes []string) (pb.DistributorService_GetNotificationStreamClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetNotificationStream(ctx,
		&pb.NotificationStreamRequest{
			OrgId:        orgId,
			NetworkId:    networkId,
			SubscriberId: subscriberId,
			UserId:       userId,
			Scopes:       scopes,
		})
}
