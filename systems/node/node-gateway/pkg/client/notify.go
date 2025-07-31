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
	"encoding/json"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/notify/pb/gen"
)

type Notify struct {
	conn    *grpc.ClientConn
	client  pb.NotifyServiceClient
	timeout time.Duration
	host    string
}

func NewNotify(notifyHost string, timeout time.Duration) *Notify {
	conn, err := grpc.NewClient(notifyHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewNotifyServiceClient(conn)

	return &Notify{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    notifyHost,
	}
}

func NewNotifyFromClient(mClient pb.NotifyServiceClient) *Notify {
	return &Notify{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (n *Notify) Close() {
	err := n.conn.Close()
	if err != nil {
		log.Warnf("Failed to gracefully close Notify Service connection: %v", err)
	}
}

func (n *Notify) Add(nodeId, severity, ntype, serviceName string, details json.RawMessage,
	status, time uint32) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	detailBytes, err := details.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return n.client.Add(ctx,
		&pb.AddRequest{
			NodeId:      nodeId,
			Severity:    severity,
			Type:        ntype,
			ServiceName: serviceName,
			Status:      status,
			Time:        time,
			Details:     detailBytes,
		})
}

func (n *Notify) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Get(ctx, &pb.GetRequest{
		NotificationId: id,
	})
}

func (n *Notify) List(nodeId, serviceName, nType string, count uint32, sort bool) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.List(ctx, &pb.ListRequest{
		NodeId:      nodeId,
		Type:        nType,
		ServiceName: serviceName,
		Count:       count,
		Sort:        sort,
	})
}

func (n *Notify) Delete(id string) (*pb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Delete(ctx, &pb.GetRequest{
		NotificationId: id,
	})
}

func (n *Notify) Purge(nodeId, serviceName, nType string) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Purge(ctx, &pb.PurgeRequest{
		NodeId:      nodeId,
		Type:        nType,
		ServiceName: serviceName,
	})
}
