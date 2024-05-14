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
	pb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

type EventNotification interface {
	Get(id string) (*pb.GetAllResponse, error)
	GetAll(orgId string, networkId string, subscriberId string, userId string, role string) (*pb.GetAllResponse, error)
	UpdateStatus(id string, isRead bool) (*pb.UpdateStatusResponse, error)
}

type eventNotification struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.EventToNotifyServiceClient
	host    string
}

func NewEventNotification(host string, timeout time.Duration) (*eventNotification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)

		return nil, err
	}

	client := pb.NewEventToNotifyServiceClient(conn)

	return &eventNotification{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}, nil
}

func NewEventToNotifyFromClient(client pb.EventToNotifyServiceClient) *eventNotification {
	return &eventNotification{
		host:    "localhost",
		timeout: 10 * time.Second,
		conn:    nil,
		client:  client,
	}
}

func (m *eventNotification) Close() {
	m.conn.Close()
}

func (n *eventNotification) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Get(ctx, &pb.GetRequest{Id: id})
}

func (n *eventNotification) UpdateStatus(id string, isRead bool) (*pb.UpdateStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.UpdateStatus(ctx, &pb.UpdateStatusRequest{Id: id, IsRead: isRead})
}

func (n *eventNotification) GetAll(orgId string, networkId string, subscriberId string, userId string, role string) (*pb.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetAll(ctx, &pb.GetAllRequest{OrgId: orgId,
		NetworkId:    networkId,
		SubscriberId: subscriberId,
		UserId:       userId,
		Role:         pb.RoleType(pb.RoleType_value[role])})
}
