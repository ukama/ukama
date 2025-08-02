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
	Get(id string) (*pb.GetResponse, error)
	GetAll(orgId string, networkId string, subscriberId string, userId string) (*pb.GetAllResponse, error)
	UpdateStatus(id string, isRead bool) (*pb.UpdateStatusResponse, error)
}

type eventNotification struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.EventToNotifyServiceClient
	host    string
}

func NewEventNotification(host string, timeout time.Duration) (*eventNotification, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to EventNotification Service: %v", err)

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

func (e *eventNotification) Close() {
	if e.conn != nil {
		err := e.conn.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close EventNotification Service connection: %v", err)
		}
	}
}

func (e *eventNotification) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.Get(ctx, &pb.GetRequest{Id: id})
}

func (e *eventNotification) UpdateStatus(id string, isRead bool) (*pb.UpdateStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.UpdateStatus(ctx, &pb.UpdateStatusRequest{Id: id, IsRead: isRead})
}

func (e *eventNotification) GetAll(orgId string, networkId string, subscriberId string, userId string) (*pb.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.GetAll(ctx, &pb.GetAllRequest{OrgId: orgId,
		NetworkId:    networkId,
		SubscriberId: subscriberId,
		UserId:       userId,
	})
}
